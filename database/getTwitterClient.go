package database

import (
	"database/sql"
	"log"
	"tweetsay/helper"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func GetTwitterClient() (*twitter.Client){

	db, db_err := sql.Open("sqlite3", getPath())
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
		keyValue := helper.Prompt_key_from_user("TWITTER_API_KEY")
		addKey("TWITTER_API_KEY", keyValue)
	}

	if TWITTER_API_KEY_SECRET == "" {
		keyValue := helper.Prompt_key_from_user("TWITTER_API_KEY_SECRET")
		addKey("TWITTER_API_KEY_SECRET", keyValue)
	}

	if TWITTER_ACCESS_TOKEN == "" {
		keyValue := helper.Prompt_key_from_user("TWITTER_ACCESS_TOKEN")
		addKey("TWITTER_ACCESS_TOKEN", keyValue)
	}

	if TWITTER_ACCESS_TOKEN_SECRET == "" {
		keyValue := helper.Prompt_key_from_user("TWITTER_ACCESS_TOKEN_SECRET")
		addKey("TWITTER_ACCESS_TOKEN_SECRET", keyValue)
	}

	config := oauth1.NewConfig(TWITTER_API_KEY, TWITTER_API_KEY_SECRET)
	token := oauth1.NewToken(TWITTER_ACCESS_TOKEN, TWITTER_ACCESS_TOKEN_SECRET)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	return client
}