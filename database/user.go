package database

import (
	"database/sql"
	"fmt"
	"log"
	"tweetsay/model"

	"github.com/dghubble/go-twitter/twitter"
)

// check if user already in database
func GetUser(username string) (*model.User, error) {
	fmt.Printf("Searching for user @%s in database\n", username)

	db, db_err := sql.Open("sqlite3", getPath())
	if db_err != nil {
		log.Fatal("Error opening database", db_err)
	}
	defer db.Close()

	var ID int64
	var CreatedAt string
	var ScreenName string
	var Active bool
	
	stmt, stmt_err := db.Prepare("SELECT * FROM Users WHERE ScreenName = ?")
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
		fmt.Printf("user @%s not found in database.\n", username)
		return nil, err
	}

	user := model.User {
		User: &twitter.User {
			ID: ID,
			ScreenName: ScreenName,
		},
		Active: Active,
	}

	return &user, err
}

// add user object to database
func AddUser(user *twitter.User) (model.User, error){
	fmt.Printf("Saving user @%s to database\n", user.ScreenName)
	
	db, db_err := sql.Open("sqlite3", getPath())
	if db_err != nil {
		log.Fatal("Error opening database", db_err)
	}

	stmt, stmt_err := db.Prepare("INSERT INTO Users (ID, Name, ScreenName) VALUES (?, ?, ?)")
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

	return model.User{
		User: &twitter.User {
			ID: user.ID,
			Name: user.Name,
			ScreenName: user.ScreenName,
		},
		Active: false,
	}, nil
}

// flips user active field
func ToggleUser(user *model.User) {
	db, db_err := sql.Open("sqlite3", getPath())
	if db_err != nil {
		log.Fatal("Error opening database", db_err)
	}
	defer db.Close()

	stmt, stmt_err := db.Prepare("UPDATE Users SET Active = ? Where ScreenName = ?")
	if stmt_err != nil {
		log.Fatalf((stmt_err.Error()))
	}
	defer stmt.Close()

	_, exec_err := stmt.Exec(!user.Active, user.User.ScreenName)
	if exec_err != nil {
		log.Fatal((exec_err.Error()))
	}
	print_string := "User @" + user.User.ScreenName + " is now"
	if !user.Active {
		print_string = print_string + " active."
	} else {
		print_string = print_string + " inactive."
	}
	fmt.Println(print_string)
}