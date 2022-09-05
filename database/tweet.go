package database

import (
	"database/sql"
	"fmt"
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

// prints a random tweet from active users,
// returns that tweet ID
func GetRandomTweet() (int64){
	db, db_err := sql.Open("sqlite3", getPath())
	if db_err != nil {
		log.Fatal("Error opening database", db_err)
	}
	defer db.Close()

	var tweet twitter.Tweet
	var username string
	row := db.QueryRow("SELECT t.ID, t.FullText, t.CreatedAt, t.Username FROM Tweets t INNER JOIN Users u on u.ScreenName = t.Username WHERE u.Active = 1  AND t.SoftDeleted = false ORDER BY RANDOM() LIMIT 1;")

	err := row.Scan(
		&tweet.ID,
		&tweet.FullText,
		&tweet.CreatedAt,
		&username,
	)
	
	if err != nil {
		log.Fatal("error while querying the database for tweets - ", err)
	}

	fmt.Printf("@%s tweets: %s\n", username, tweet.FullText)
	return tweet.ID
}


// updates the last shown tweet table with given tweetID
func SetLastShownTweet(tweetID int64) {
	db, db_err := sql.Open("sqlite3", getPath())
	if db_err != nil {
		log.Fatal("Error opening database", db_err)
	}

	stmt, stmt_err := db.Prepare("INSERT INTO ShownTweets (TweetID) VALUES (?)")
	if stmt_err != nil {
		log.Fatalf((stmt_err.Error()))
	}

	_, err := stmt.Exec(tweetID)
	if err != nil {
		log.Fatal((err.Error()))
	}
	defer stmt.Close()
}
