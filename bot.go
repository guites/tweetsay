package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	_ "github.com/mattn/go-sqlite3"
)

func thething(client *twitter.Client) {
	user, _, err := client.Accounts.VerifyCredentials(&twitter.AccountVerifyParams{})
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}
	fmt.Printf("Account: @%s (%s)", user.ScreenName, user.Name)
	
	tweet, _, err := client.Statuses.Show(user.Status.ID, &twitter.StatusShowParams{
		TweetMode: "extended",
	})

	if err != nil {
		fmt.Printf("err: %s\n", err)
	}
	fmt.Printf("%s\n", tweet.FullText)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	available_commands := "Available options: [add_user_timeline|list_users]"

	db_err := createTables()
	if db_err != nil {
		log.Fatalf("Could not start sqlite3 database: %s", db_err.Error())
	}

	client := prepare_twitter_api()
	cmdArgs := os.Args[1:]
	if len(cmdArgs) < 1 {
		get_random_tweet("Antho_Repartie")
		return
	}
	chosenOption := cmdArgs[0]
	switch chosenOption {
	case "add_user_timeline":
		if len(cmdArgs) < 2 {
			log.Fatal("Usage: add_user_timeline @user_handle")
		}
		username := removeAt(cmdArgs[1])
		user := get_twitter_user_by_username(username, client)
		get_user_timeline(user, client)
	case "list_users":
		list_users()
	default:
		fmt.Println(available_commands)
	}
}