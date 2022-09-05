package database

import (
	"database/sql"
	"log"

	"github.com/dghubble/go-twitter/twitter"
)

// add tweet object to database
func AddTweet (tweet *twitter.Tweet) {
	// log.Printf("Saving tweet @%d to database: https://twitter.com/%s/status/%d\n", tweet.ID, tweet.User.ScreenName, tweet.ID)
	
	db, db_err := sql.Open("sqlite3", getPath())
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