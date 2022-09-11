package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"tweetsay/helper"
	"unicode/utf8"
)

// adds word to database, returns that word id
func AddWord (word string) int{
	db, db_err := sql.Open("sqlite3", getPath())
	if db_err != nil {
		log.Fatal("Error opening database", db_err)
	}
	defer db.Close()
	
	stmt, stmt_err := db.Prepare("INSERT INTO Words (word) VALUES (?)")
	if stmt_err != nil {
		log.Fatalf((stmt_err.Error()))
	}
	defer stmt.Close()

	// add all words in uppercase to database
	stripped_word := helper.RemoveSymbols(word)
	uppercase_word := strings.ToUpper(stripped_word)
	_, err := stmt.Exec(uppercase_word)
	if err != nil {
		if ! strings.Contains(err.Error(), "UNIQUE") {
			log.Fatalf("Error while inserting word %s into database, details: %s", word, err.Error())
		}
	}

	// whether the word was already in the database or was just inserted,
	// we need another query to retrieve its ID

	var wordID int
	row := db.QueryRow("SELECT ID FROM Words WHERE word = ?", uppercase_word)
	err_row := row.Scan(&wordID)
	if err_row != nil {
		log.Fatal((err_row.Error()))
	}

	return wordID
}

func RelateWordToTweet(wordID int, tweetID int64) {
	db, db_err := sql.Open("sqlite3", getPath())
	if db_err != nil {
		log.Fatal("Error opening database", db_err)
	}
	defer db.Close()

	stmt, stmt_err := db.Prepare("INSERT INTO WordAppearances (WordID, TweetID) VALUES (?, ?)")
	if stmt_err != nil {
		log.Fatalf((stmt_err.Error()))
	}
	defer stmt.Close()

	_, err := stmt.Exec(wordID, tweetID)
	if err != nil {
		if ! strings.Contains(err.Error(), "UNIQUE") {
			log.Fatalf("Error while creating word - tweet relationship (%d, %d) in  database, details: %s", wordID, tweetID, err.Error())
		}
	}
}

// TODO: maybe define this in models package?
type TweetExcerpt struct {
	Username string
	TweetID int64
	TweetExcerpt string
}

// return a list of tweets containing the given word
func FindTweetsContainingWord(word string) []TweetExcerpt {
	db, db_err := sql.Open("sqlite3", getPath())
	if db_err != nil {
		log.Fatal("Error opening database", db_err)
	}
	defer db.Close()

	var tweetExcerpts []TweetExcerpt
	var tweetID int64
	var tweetFullText string
	var tweetUsername string
	

	// words are all uppercased in database
	uppercase_word := strings.ToUpper(word)
	rows, query_err := db.Query("SELECT t.id, t.FullText, t.Username FROM Tweets t INNER JOIN WordAppearances wa ON wa.TweetID = t.ID INNER JOIN Words w ON w.ID = wa.WordID WHERE w.word = (?) AND SoftDeleted = 0", uppercase_word)
	if query_err != nil {
		log.Fatalf("Error while searching for tweets containing the word <<%s>>, detail: %s\n", word, query_err.Error())
	}

	for rows.Next() {
		scan_err := rows.Scan(&tweetID, &tweetFullText, &tweetUsername)
		if scan_err != nil {
			log.Fatalln("Error while iterating over results from db:", scan_err)
		}
		uppercase_fulltext := strings.ToUpper(tweetFullText)
		word_len := utf8.RuneCountInString(word)
		word_offset := strings.Index(uppercase_fulltext, uppercase_word)
		start_index := word_offset - 20
		if start_index < 0 {
			start_index = 0
		}
		end_index := word_offset + word_len + 20
		if end_index > len(tweetFullText) {
			end_index = len(tweetFullText)
		}
		tweetExcerpt := TweetExcerpt{
			Username: tweetUsername,
			TweetID: tweetID,
			TweetExcerpt: "..." + tweetFullText[start_index : end_index] + "...",
		}
		tweetExcerpts = append(tweetExcerpts, tweetExcerpt)
	}

	if len(tweetExcerpts) == 0 {
		fmt.Printf("You haven't seen any tweet containing the word <<%s>> yet!\n", word)
		os.Exit(0)
	}

	return tweetExcerpts
}