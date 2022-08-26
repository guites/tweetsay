package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/dghubble/go-twitter/twitter"
)

// check if user already in database
func get_twitter_user_from_db(username string) (*User, error) {
	log.Printf("Searching for user @%s in database\n", username)

	db, db_err := sql.Open("sqlite3", getDbPath())
	if db_err != nil {
		log.Fatal("Error opening database", db_err)
	}
	defer db.Close()

	var ID int64
	var CreatedAt string
	var ScreenName string
	var Active bool
	
	stmt, stmt_err := db.Prepare("SELECT * FROM DatabaseUsers WHERE ScreenName = ?")
	if stmt_err != nil {
		log.Fatal("Error preparing statement", stmt_err)
	}
	defer stmt.Close()

	err := stmt.QueryRow(username).Scan(
		&ID,
		&CreatedAt,
		&ScreenName,
		&Active,
	)
	
	if err != nil {
		log.Printf("user @%s not found in database\n", username)
		return nil, err
	}

	user := User {
		User: &twitter.User {
			ID: ID,
			ScreenName: ScreenName,
		},
		Active: Active,
	}

	return &user, err
}

// add user object to database
func add_twitter_user_to_db (user *twitter.User) (User, error){
	log.Printf("Saving user @%s to database\n", user.ScreenName)
	
	db, db_err := sql.Open("sqlite3", getDbPath())
	if db_err != nil {
		log.Fatal("Error opening database", db_err)
	}

	stmt, stmt_err := db.Prepare("INSERT INTO DatabaseUsers (ID, Name, ScreenName) VALUES (?, ?, ?)")
	if stmt_err != nil {
		log.Fatalf((stmt_err.Error()))
	}
	defer stmt.Close()

	_, err := stmt.Exec(
		user.ID,
		user.Name,
		user.ScreenName,
	)
	if err != nil {
		log.Fatal((err.Error()))
	}

	return User{
		User: &twitter.User {
			ID: user.ID,
			Name: user.Name,
			ScreenName: user.ScreenName,
		},
		Active: false,
	}, nil
}

// prints a list of all users currently in database
func list_users() {
	fmt.Println("Listing all users registered to the database:")
	fmt.Println("---------------------------------------------")

	db, db_err := sql.Open("sqlite3", getDbPath())
	if db_err != nil {
		log.Fatal("Error opening database", db_err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT u.ID, u.Name, u.ScreenName, count(u.ScreenName), u.Active FROM DatabaseUsers u JOIN DatabaseTweets t ON u.ScreenName = t.Username GROUP BY u.ID")
	if err != nil {
		log.Fatal("Error while querying the database for users", err)
	}
	defer rows.Close()

	var ID int64
	var Name string
	var ScreenName string
	var DbTweetCount int
	var Active bool
	var ActiveStr string

	for rows.Next() {
		err := rows.Scan(
			&ID,
			&Name,
			&ScreenName,
			&DbTweetCount,
			&Active,
		)
		if err != nil {
			log.Fatal("Error while iterating over timeline results from db", err)
		}
		if Active {
			ActiveStr = "[x]"
		} else {
			ActiveStr = "[ ]"
		}
		fmt.Printf("@%s (%s) %s - %d tweets\n", ScreenName, Name, ActiveStr, DbTweetCount)
	}
}

// flips user active column
func toggle_user(username string) {
	user, err := get_twitter_user_from_db(username)
	if err != nil {
		fmt.Printf("User @%s not registered in database. Please run <add_user_timeline @%s>\n", username, username)
	}
	
	db, db_err := sql.Open("sqlite3", getDbPath())
	if db_err != nil {
		log.Fatal("Error opening database", db_err)
	}
	defer db.Close()

	stmt, stmt_err := db.Prepare("UPDATE DatabaseUsers SET Active = ? Where ScreenName = ?")
	if stmt_err != nil {
		log.Fatalf((stmt_err.Error()))
	}
	defer stmt.Close()

	_, exec_err := stmt.Exec(!user.Active, user.User.ScreenName)
	if exec_err != nil {
		log.Fatal((err.Error()))
	}
	print_string := "User @" + user.User.ScreenName + " is now"
	if !user.Active {
		print_string = print_string + " active."
	} else {
		print_string = print_string + " inactive."
	}
	fmt.Println(print_string)

}