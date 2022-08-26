package main

import (
	"fmt"
	"log"

	"github.com/dghubble/go-twitter/twitter"
)

func get_twitter_user_by_username(username string, client *twitter.Client) (*User){

	db_user, db_err := get_twitter_user_from_db(username)
	if db_err == nil {
		return db_user
	}

	log.Printf("Fetching user @%s from Twitter API\n", username)
	twitter_user, _, err := client.Users.Show(&twitter.UserShowParams{
		ScreenName: username, // "Antho_Repartie",
	})

	if err != nil {
		log.Fatalf("err: %s\n", err)
	}

	user, user_err := add_twitter_user_to_db(twitter_user)
	if user_err != nil {
		log.Fatal("Error while adding user to database", user_err)
	}
	fmt.Printf("Account: @%s (%s) (%d)\n", user.User.ScreenName, user.User.Name, user.User.ID)
	return &user
}