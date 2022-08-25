package main

import (
	"database/sql"
	"envholder"
	"fmt"
	"log"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	_ "github.com/mattn/go-sqlite3"
)

func prepare_twitter_api() (*twitter.Client){
	var envVars envholder.EnvHolder = envholder.LoadEnv()
	consumerKey := envVars.GetVar("TWITTER_API_KEY")
	consumerSecret := envVars.GetVar("TWITTER_API_KEY_SECRET")
	accessToken := envVars.GetVar("TWITTER_ACCESS_TOKEN")
	accessTokenSecret := envVars.GetVar("TWITTER_ACCESS_TOKEN_SECRET")

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	return client
}

// check if user already in database
func get_twitter_user_from_db(username string) (*twitter.User, error) {
	db, db_err := sql.Open("sqlite3", "./updates.db")
	if db_err != nil {
		log.Fatal("Error opening database", db_err)
	}

	db_user, err := getUserFromDatabase(username, db)
	return db_user, err
}

// add user object to database
func add_twitter_user_to_db (user *twitter.User) {
	db, db_err := sql.Open("sqlite3", "./updates.db")
	if db_err != nil {
		log.Fatal("Error opening database", db_err)
	}

	stmt, stmt_err := db.Prepare("INSERT INTO DatabaseUsers (ID, Name, ScreenName) VALUES (?, ?, ?)")

	if stmt_err != nil {
		log.Fatalf((stmt_err.Error()))
	}
	res, err := stmt.Exec(
		user.ID,
		user.Name,
		user.ScreenName,
	)

	if err != nil {
		log.Fatal((err.Error()))
	}
	defer stmt.Close()
	log.Print(res)
}

func get_twitter_user_by_username(username string, client *twitter.Client) (*twitter.User){

	db_user, db_err := get_twitter_user_from_db(username)
	if db_err == nil {
		return db_user
	}

	user, _, err := client.Users.Show(&twitter.UserShowParams{
		ScreenName: username, // "Antho_Repartie",
	})

	if err != nil {
		log.Fatalf("err: %s\n", err)
	}

	add_twitter_user_to_db(user)

	fmt.Printf("Account: @%s (%s) (%d)\n", user.ScreenName, user.Name, user.ID)
	fmt.Printf("Last Tweet ID: %d\n", user.Status.ID)
	return user
}

func get_user_timeline(user *twitter.User, client *twitter.Client) ([]twitter.Tweet) {
		tweets, _, err := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		UserID: user.ID,
    	Count: 3200,
		TweetMode: "extended",
	})
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}
	for _, tweet := range tweets {
		fmt.Printf("%s\n", tweet.FullText)
	}
	return tweets
}

func add_user_timeline(client *twitter.Client) {

}

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

// remove @ if user included in command line arg
func removeAt(username string) (string) {
	if username[0:1] == "@" {
		username = username[1:]
	}
	return username
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	available_commands := "Available options: [add_user_timeline]"

	db_err := createTables()
	if db_err != nil {
		log.Fatalf("Could not start sqlite3 database: %s", db_err.Error())
	}

	client := prepare_twitter_api()
	cmdArgs := os.Args[1:]
	if len(cmdArgs) < 1 {
		log.Fatal(available_commands)
	}
	chosenOption := cmdArgs[0]
	switch chosenOption {
	case "add_user_timeline":
		if len(cmdArgs) < 2 {
			log.Fatal("Usage: add_user_timeline @user_handle")
		}
		username := removeAt(cmdArgs[1])
		get_twitter_user_by_username(username, client)
		// get_user_timeline(user, client)
	default:
		fmt.Println(available_commands)
	}
}