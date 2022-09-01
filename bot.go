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

	available_commands := `syntax: tweetsay [delete_last|list_words|list_users|add_user_timeline|toggle_user|help]
options:
delete_last                  Deletes the latest shown tweet.
list_words                   List words present in last tweet alongside a link to their wiktionary page.
list_users                   Lists all currently indexed users.
add_timeline <@user_name>    Indexes a user timeline into the app.
toggle_user <@user_name>     Toggles whether @user_name tweets can be shown on new terminals.
help                         Shows this help text.
	`

	db_err := createTables()
	if db_err != nil {
		log.Fatalf("Could not start sqlite3 database: %s", db_err.Error())
	}
	client := prepare_twitter_api()
	cmdArgs := os.Args[1:]
	if len(cmdArgs) < 1 {
		tweetID := get_random_tweet()
		set_last_shown_tweet(tweetID)
		return
	}

	chosenOption := cmdArgs[0]
	switch chosenOption {
	case "add_timeline":
		if len(cmdArgs) < 2 {
			log.Fatal("Usage: add_timeline @user_handle")
		}
		username := removeAt(cmdArgs[1])
		user := get_twitter_user_by_username(username, client)
		get_user_timeline(user, client)
	case "list_users":
		list_users()
	case "toggle_user":
		if len(cmdArgs) < 2 {
			log.Fatal("Usage: toggle_user @user_handle")
		}
		username := removeAt(cmdArgs[1])
		toggle_user(username)
	case "list_words":
		list_words_from_last_tweet()
	case "delete_last":
		delete_last_shown_tweet()
	case "help":
		fmt.Println(available_commands)
	default:
		fmt.Println(available_commands)
	}
}
