package main

import "database/sql"

const create string = `
	CREATE TABLE IF NOT EXISTS DatabaseTweets (
		ID INTEGER NOT NULL,
		CreatedAt TEXT NOT NULL,
		FullText TEXT NOT NULL,
		Username TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS DatabaseUsers (
		ID INTEGER NOT NULL,
		Name TEXT NOT NULL,
		ScreenName TEXT NOT NULL
	);
	`

const file string = "updates.db"

func createTables() (error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return err
	}
	defer db.Close()
	if _, err := db.Exec(create); err != nil {
		return err
	}
	return nil
}