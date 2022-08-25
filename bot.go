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
	log.Printf("Searching for user @%s in database\n", username)
	
	defer db.Close()

	var user twitter.User
	stmt, stmt_err := db.Prepare("SELECT * FROM DatabaseUsers WHERE ScreenName = ?")
	
	if stmt_err != nil {
		log.Fatal("Error preparing statement", stmt_err)
	}

	defer stmt.Close()

	err := stmt.QueryRow(username).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.ScreenName,
	)
	
	if err != nil {
		log.Printf("user @%s not found in database\n", username)
		return nil, err
	}

	log.Printf("user @%s already in database\n", username)
	return &user, err
}

// add user object to database
func add_twitter_user_to_db (user *twitter.User) {
	log.Printf("Saving user @%s to database\n", user.ScreenName)
	
	db, db_err := sql.Open("sqlite3", "./updates.db")
	if db_err != nil {
		log.Fatal("Error opening database", db_err)
	}
	stmt, stmt_err := db.Prepare("INSERT INTO DatabaseUsers (ID, Name, ScreenName) VALUES (?, ?, ?)")

	if stmt_err != nil {
		log.Fatalf((stmt_err.Error()))
	}
	_, err := stmt.Exec(
		user.ID,
		user.Name,
		user.ScreenName,
	)

	if err != nil {
		log.Fatal((err.Error()))
	}
	defer stmt.Close()
}

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

// add tweet object to database
func add_tweet_to_db (tweet *twitter.Tweet) {
	log.Printf("Saving tweet @%d to database: https://twitter.com/%s/status/%d\n", tweet.ID, tweet.User.ScreenName, tweet.ID)
	
	db, db_err := sql.Open("sqlite3", "./updates.db")
	if db_err != nil {
		log.Fatal("Error opening database", db_err)
	}
	stmt, stmt_err := db.Prepare("INSERT INTO DatabaseTweets (ID, CreatedAt, FullText, UserName) VALUES (?, ?, ?, ?)")

	if stmt_err != nil {
		log.Fatalf((stmt_err.Error()))
	}
	_, err := stmt.Exec(
		tweet.ID,
		tweet.CreatedAt,
		tweet.FullText,
		tweet.User.ScreenName,
	)

	if err != nil {
		log.Fatal((err.Error()))
	}
	defer stmt.Close()
}

// check if user timeline already in database
func get_user_timeline_from_db(user *twitter.User) (bool) {
	log.Printf("Searching for tweets from @%s in database\n", user.ScreenName)
	
	db, db_err := sql.Open("sqlite3", "./updates.db")
	if db_err != nil {
		log.Fatal("Error opening database", db_err)
	}
	defer db.Close()

	stmt, stmt_err := db.Prepare("SELECT * FROM DatabaseTweets WHERE UserName = ?")
	if stmt_err != nil {
		log.Fatal("Error preparing statement", stmt_err)
	}
	defer stmt.Close()

	rows, rows_err := stmt.Query(user.ScreenName)
	if rows_err != nil {
		log.Fatal("Error executing statement:", rows_err)
	}
	defer rows.Close()

	var timeline []twitter.Tweet
	var ID int64
	var CreatedAt string
	var FullText string
	var UserName string

	for rows.Next() {
		err := rows.Scan(
			&ID,
			&CreatedAt,
			&FullText,
			&UserName,
		)
		if err != nil {
			log.Fatal("Error while iterating over timeline results from db", err)
		}
		timeline = append(timeline, twitter.Tweet{ID: ID, CreatedAt: CreatedAt, FullText: FullText})
	}
	if len(timeline) > 0 {
		log.Printf("found %d tweets from @%s in the database\n", len(timeline), user.ScreenName)
		return true
	}
	return false
}


func get_user_timeline(user *twitter.User, client *twitter.Client) () {

	if get_user_timeline_from_db(user) {
		return
	}
	log.Printf("Fetching user @%s tweets from Twitter API\n", user.ScreenName)
	tweets, _, err := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		UserID: user.ID,
    	Count: 3200,
		TweetMode: "extended",
	})
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}
	for _, tweet := range tweets {
		// fmt.Printf("%s\n", tweet.FullText)
		add_tweet_to_db(&tweet)
	}
	return
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
		user := get_twitter_user_by_username(username, client)
		get_user_timeline(user, client)
	default:
		fmt.Println(available_commands)
	}
}