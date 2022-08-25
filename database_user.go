package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/dghubble/go-twitter/twitter"
)

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

// prints a list of all users currently in database
func list_users() {
	fmt.Println("Listing all users registered to the database:")
	fmt.Println("---------------------------------------------")

	db, db_err := sql.Open("sqlite3", "./updates.db")
	if db_err != nil {
		log.Fatal("Error opening database", db_err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT u.ID, u.Name, u.ScreenName, count(u.ScreenName) FROM DatabaseUsers u JOIN DatabaseTweets t ON u.ScreenName = t.Username GROUP BY u.ID")
	if err != nil {
		log.Fatal("Error while querying the database for users", err)
	}
	defer rows.Close()

	// type UserWithDBTweetCount struct {
	// 	User twitter.User
	// 	DbTweetCount int
	// }

	// var db_users []UserWithDBTweetCount
	var ID int64
	var Name string
	var ScreenName string
	var DbTweetCount int

	for rows.Next() {
		err := rows.Scan(
			&ID,
			&Name,
			&ScreenName,
			&DbTweetCount,
		)
		if err != nil {
			log.Fatal("Error while iterating over timeline results from db", err)
		}
		fmt.Printf("@%s (%s) - %d tweets\n", ScreenName, Name, DbTweetCount)
		// db_users = append(db_users, UserWithDBTweetCount{User: twitter.User{ID: ID, Name: Name, ScreenName: ScreenName}, DbTweetCount: DbTweetCount})
	}

}