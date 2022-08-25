package main

import (
	"database/sql"

	"github.com/dghubble/go-twitter/twitter"
	_ "github.com/mattn/go-sqlite3"
)

func getUserFromDatabase (username string, db *sql.DB) (*twitter.User, error){
	var user twitter.User
	err := db.
	QueryRow("SELECT * FROM DatabaseUsers;").
	Scan(
		&user.ID,
		&user.CreatedAt,
		&user.ScreenName,
	)
	if err != nil {
		return nil, err
	}
	return &user, err
}
