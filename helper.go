package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// remove @ if user included in command line arg
func removeAt(username string) (string) {
	if username[0:1] == "@" {
		username = username[1:]
	}
	return username
}

func add_key_to_db(keyName string, keyValue string)() {
	log.Printf("Saving key @%s to database\n", keyName)

	db, db_err := sql.Open("sqlite3", getDbPath())
	if db_err != nil {
		log.Fatal("Error opening database", db_err)
	}
	sql_string := "UPDATE TwitterCredentials SET " + keyName + " = ?"
	stmt, stmt_err := db.Prepare(sql_string)

	if stmt_err != nil {
		log.Fatalf((stmt_err.Error()))
	}
	_, err := stmt.Exec(keyValue)

	if err != nil {
		log.Fatal((err.Error()))
	}
	defer stmt.Close()
}

func prompt_key_from_user(keyName string) (string) {
	var keyValue string
	fmt.Printf("Please enter you %s: ", keyName)
	fmt.Scanf("%s", &keyValue)
	return keyValue
}

func prepare_twitter_api() (*twitter.Client){

	db, db_err := sql.Open("sqlite3", getDbPath())
	if db_err != nil {
		log.Fatal("Error opening database", db_err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT TWITTER_API_KEY, TWITTER_API_KEY_SECRET, TWITTER_ACCESS_TOKEN, TWITTER_ACCESS_TOKEN_SECRET FROM TwitterCredentials")
	if err != nil {
		log.Fatal("Error while querying the database for users", err)
	}
	defer rows.Close()

	var TWITTER_API_KEY string
	var TWITTER_API_KEY_SECRET string
	var TWITTER_ACCESS_TOKEN string
	var TWITTER_ACCESS_TOKEN_SECRET string

	db.QueryRow("SELECT TWITTER_API_KEY, TWITTER_API_KEY_SECRET, TWITTER_ACCESS_TOKEN, TWITTER_ACCESS_TOKEN_SECRET FROM TwitterCredentials").Scan(
		&TWITTER_API_KEY,
		&TWITTER_API_KEY_SECRET,
		&TWITTER_ACCESS_TOKEN,
		&TWITTER_ACCESS_TOKEN_SECRET,
	)

	if TWITTER_API_KEY == "" {
		keyValue := prompt_key_from_user("TWITTER_API_KEY")
		add_key_to_db("TWITTER_API_KEY", keyValue)
	}

	if TWITTER_API_KEY_SECRET == "" {
		keyValue := prompt_key_from_user("TWITTER_API_KEY_SECRET")
		add_key_to_db("TWITTER_API_KEY_SECRET", keyValue)
	}

	if TWITTER_ACCESS_TOKEN == "" {
		keyValue := prompt_key_from_user("TWITTER_ACCESS_TOKEN")
		add_key_to_db("TWITTER_ACCESS_TOKEN", keyValue)
	}

	if TWITTER_ACCESS_TOKEN_SECRET == "" {
		keyValue := prompt_key_from_user("TWITTER_ACCESS_TOKEN_SECRET")
		add_key_to_db("TWITTER_ACCESS_TOKEN_SECRET", keyValue)
	}

	config := oauth1.NewConfig(TWITTER_API_KEY, TWITTER_API_KEY_SECRET)
	token := oauth1.NewToken(TWITTER_ACCESS_TOKEN, TWITTER_ACCESS_TOKEN_SECRET)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	return client
}

func getDbPath() (string){
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	file := dirname + "/.tweetsay.db"
	return file
}

func set_last_shown_tweet(tweetID int64) {
	db, db_err := sql.Open("sqlite3", getDbPath())
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
