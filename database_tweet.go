package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/dghubble/go-twitter/twitter"
)

// add tweet object to database
func add_tweet_to_db (tweet *twitter.Tweet) {
	log.Printf("Saving tweet @%d to database: https://twitter.com/%s/status/%d\n", tweet.ID, tweet.User.ScreenName, tweet.ID)
	
	db, db_err := sql.Open("sqlite3", getDbPath())
	if db_err != nil {
		log.Fatal("Error opening database", db_err)
	}
	stmt, stmt_err := db.Prepare("INSERT INTO DatabaseTweets (ID, CreatedAt, FullText, UserName) VALUES (?, ?, ?, ?)")

	if stmt_err != nil {
		log.Fatalf((stmt_err.Error()))
	}
	_, err := stmt.Exec(
		tweet.ID,
		tweet.CreatedAt,
		tweet.FullText,
		tweet.User.ScreenName,
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

	stmt, stmt_err := db.Prepare("SELECT * FROM DatabaseTweets WHERE UserName = ?")
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
	var ID int64
	var CreatedAt string
	var FullText string
	var UserName string

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

// prints a random tweet from active users
func get_random_tweet() {
	db, db_err := sql.Open("sqlite3", getDbPath())
	if db_err != nil {
		log.Fatal("Error opening database", db_err)
	}
	defer db.Close()

	var tweet twitter.Tweet
	row := db.QueryRow("SELECT t.ID, t.FullText, t.CreatedAt FROM DatabaseTweets t INNER JOIN DatabaseUsers u on u.ScreenName = t.Username WHERE u.Active = 1 ORDER BY RANDOM() LIMIT 1;")

	err := row.Scan(
		&tweet.ID,
		&tweet.FullText,
		&tweet.CreatedAt,
	)
	
	if err != nil {
		log.Fatal("error while querying the database for tweets - ", err)
	}

	fmt.Println(tweet.FullText)
}