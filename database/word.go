package database

import (
	"database/sql"
	"log"
	"strings"
)

// adds word to database, returns that word id
func AddWord (word string) int{
	db, db_err := sql.Open("sqlite3", getPath())
	if db_err != nil {
		log.Fatal("Error opening database", db_err)
	}
	
	stmt, stmt_err := db.Prepare("INSERT INTO Words (word) VALUES (?)")
	if stmt_err != nil {
		log.Fatalf((stmt_err.Error()))
	}
	defer stmt.Close()

	_, err := stmt.Exec(word)
	if err != nil {
		if ! strings.Contains(err.Error(), "UNIQUE") {
			log.Fatalf("Error while inserting word %s into database, details: %s", word, err.Error())
		}
	}

	// whether the word was already in the database or was just inserted,
	// we need another query to retrieve its ID

	var wordID int
	row := db.QueryRow("SELECT ID FROM Words WHERE word = ?", word)
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