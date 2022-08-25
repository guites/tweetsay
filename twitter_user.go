package main

import (
	"fmt"
	"log"

	"github.com/dghubble/go-twitter/twitter"
)

func get_twitter_user_by_username(username string, client *twitter.Client) (*twitter.User){

	db_user, db_err := get_twitter_user_from_db(username)
	if db_err == nil {
		return db_user
	}

	log.Printf("Fetching user @%s from Twitter API\n", username)
	user, _, err := client.Users.Show(&twitter.UserShowParams{
		ScreenName: username, // "Antho_Repartie",
	})

	if err != nil {
		log.Fatalf("err: %s\n", err)
	}

	add_twitter_user_to_db(user)

	fmt.Printf("Account: @%s (%s) (%d)\n", user.ScreenName, user.Name, user.ID)
	fmt.Printf("Last Tweet ID: %d\n", user.Status.ID) // TODO: maybe this errors out when user has never tweeted
	return user
}