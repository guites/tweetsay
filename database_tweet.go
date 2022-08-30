package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/dghubble/go-twitter/twitter"
)

// add tweet object to database
func add_tweet_to_db (tweet *twitter.Tweet) {
	log.Printf("Saving tweet @%d to database: https://twitter.com/%s/status/%d\n", tweet.ID, tweet.User.ScreenName, tweet.ID)
	
	db, db_err := sql.Open("sqlite3", getDbPath())
	if db_err != nil {
		log.Fatal("Error opening database", db_err)
	}
	stmt, stmt_err := db.Prepare("INSERT INTO Tweets (ID, CreatedAt, FullText, UserName, Lang) VALUES (?, ?, ?, ?, ?)")

	if stmt_err != nil {
		log.Fatalf((stmt_err.Error()))
	}
	_, err := stmt.Exec(
		tweet.ID,
		tweet.CreatedAt,
		tweet.FullText,
		tweet.User.ScreenName,
		tweet.Lang,
	)

	if err != nil {
		log.Fatal((err.Error()))
	}
	defer stmt.Close()
}

// check if user already has tweets in database
func get_user_timeline_from_db(user *twitter.User) (bool) {
	log.Printf("Searching for tweets from @%s in database\n", user.ScreenName)
	
	db, db_err := sql.Open("sqlite3", getDbPath())
	if db_err != nil {
		log.Fatal("Error opening database", db_err)
	}
	defer db.Close()

	stmt, stmt_err := db.Prepare("SELECT * FROM Tweets WHERE UserName = ?")
	if stmt_err != nil {
		log.Fatal("Error preparing statement", stmt_err)
	}
	defer stmt.Close()

	rows, rows_err := stmt.Query(user.ScreenName)
	if rows_err != nil {
		log.Fatal("Error executing statement:", rows_err)
	}
	defer rows.Close()

	var timeline []twitter.Tweet
	var ID        int64
	var CreatedAt string
	var FullText  string
	var UserName  string

	for rows.Next() {
		err := rows.Scan(
			&ID,
			&CreatedAt,
			&FullText,
			&UserName,
		)
		if err != nil {
			log.Fatal("Error while iterating over timeline results from db", err)
		}
		timeline = append(timeline, twitter.Tweet{ID: ID, CreatedAt: CreatedAt, FullText: FullText})
	}
	if len(timeline) > 0 {
		log.Printf("found %d tweets from @%s in the database\n", len(timeline), user.ScreenName)
		return true
	}
	return false
}

// prints a random tweet from active users,
// returns that tweet ID
func get_random_tweet() (int64){
	db, db_err := sql.Open("sqlite3", getDbPath())
	if db_err != nil {
		log.Fatal("Error opening database", db_err)
	}
	defer db.Close()

	var tweet twitter.Tweet
	row := db.QueryRow("SELECT t.ID, t.FullText, t.CreatedAt FROM Tweets t INNER JOIN Users u on u.ScreenName = t.Username WHERE u.Active = 1  AND t.SoftDeleted = false ORDER BY RANDOM() LIMIT 1;")

	err := row.Scan(
		&tweet.ID,
		&tweet.FullText,
		&tweet.CreatedAt,
	)
	
	if err != nil {
		log.Fatal("error while querying the database for tweets - ", err)
	}

	fmt.Println(tweet.FullText)
	return tweet.ID
}

func delete_last_shown_tweet() {
	db, db_err := sql.Open("sqlite3", getDbPath())
	if db_err != nil {
		log.Fatal("Error opening database", db_err)
	}
	defer db.Close()

	var tweet twitter.Tweet
	var lastShownTweetID int
	row := db.QueryRow("SELECT t.ID, t.FullText, t.CreatedAt, t.Lang, s.ID FROM Tweets t INNER JOIN ShownTweets s ON t.ID = s.TweetID ORDER BY s.ID DESC LIMIT 1;")

	err := row.Scan(
		&tweet.ID,
		&tweet.FullText,
		&tweet.CreatedAt,
		&tweet.Lang,
		&lastShownTweetID,
	)

	if err != nil {
		log.Fatal("error while querying the database for tweets - ", err)
	}

	fmt.Println("Delete tweet with ID",tweet.ID,"?")
	if len(tweet.FullText) > 30 {
		fmt.Println(tweet.FullText[:30],"...")
	} else {
		fmt.Println(tweet.FullText,"...")
	}
	
	var confirmDelete string
    fmt.Print("Please type 'delete' to confirm removal:")
    fmt.Scanf("%s", &confirmDelete)
	if strings.ToLower(confirmDelete) == "delete" {

		// Sets tweet's SoftDeleted flag to TRUE
		stmt, stmt_err := db.Prepare("UPDATE Tweets SET SoftDeleted = ? Where ID = ?")
		if stmt_err != nil {
			log.Fatalf((stmt_err.Error()))
		}
		defer stmt.Close()

		_, exec_err := stmt.Exec(true, tweet.ID)
		if exec_err != nil {
			log.Fatal((err.Error()))
		}

		del_stmt, del_stmt_err := db.Prepare("DELETE FROM ShownTweets WHERE ID = ?")
		if del_stmt_err != nil {
			log.Fatalf((del_stmt_err.Error()))
		}
		defer del_stmt.Close()

		_, del_exec_err := del_stmt.Exec(lastShownTweetID)
		if del_exec_err != nil {
			log.Fatal((err.Error()))
		}

		fmt.Println("Tweet was deleted successfully.")
	} else {
		fmt.Println("Tweet was not deleted.")
	}
}

func list_words_from_last_tweet() {
	db, db_err := sql.Open("sqlite3", getDbPath())
	if db_err != nil {
		log.Fatal("Error opening database", db_err)
	}
	defer db.Close()

	var tweet twitter.Tweet
	row := db.QueryRow("SELECT t.ID, t.FullText, t.CreatedAt, t.Lang FROM Tweets t INNER JOIN ShownTweets s ON t.ID = s.TweetID ORDER BY s.ID DESC LIMIT 1;")

	err := row.Scan(
		&tweet.ID,
		&tweet.FullText,
		&tweet.CreatedAt,
		&tweet.Lang,
	)
	
	if err != nil {
		log.Fatal("error while querying the database for tweets - ", err)
	}

	// formats wiktionary url based on tweet Lang
	var wiktionary_url string
	if tweet.Lang == "und" ||
	tweet.Lang == "qme" || 
	tweet.Lang == "" {
		wiktionary_url = "https://wiktionary.org/wiki/"
	} else {
		wiktionary_url =  "https://" + tweet.Lang + ".wiktionary.org/wiki/"
	}

	// fmt.Println(tweet.FullText)
	words := strings.Fields(tweet.FullText)

	var max_length int
	for _, word := range words {
		curr_length := utf8.RuneCountInString(word)
		if curr_length > max_length {
			max_length = curr_length
		}
	}

	for index, word := range words {
		if word == "!" ||
		word == "?" ||
		word == "-" ||
		word == ":" {
			continue
		}
		wiktionary_link := wiktionary_url + word
		if strings.HasPrefix(word, "http") {
			wiktionary_link = ""
		}
		wrapped_word := word
		if utf8.RuneCountInString(word) < max_length {
			wrapped_word = word + strings.Repeat(" ", max_length - utf8.RuneCountInString(word))
		}
		padded_index := "0"+strconv.Itoa(index + 1)
		fmt.Println(padded_index[len(padded_index)-2:], wrapped_word, wiktionary_link)
	}
}