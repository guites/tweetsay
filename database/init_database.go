package database

import (
	"database/sql"
)

const create string = `
	CREATE TABLE IF NOT EXISTS Tweets (
		ID INTEGER NOT NULL,
		CreatedAt TEXT NOT NULL,
		FullText TEXT NOT NULL,
		Lang TEXT NOT NULL DEFAULT "",
		Username TEXT NOT NULL,
		SoftDeleted BOOLEAN NOT NULL DEFAULT 0 CHECK (SoftDeleted IN (0,1))
	);
	CREATE TABLE IF NOT EXISTS Words (
		ID INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		Word TEXT UNIQUE
	);
	CREATE TABLE IF NOT EXISTS WordAppearances (
		WordID INTEGER NOT NULL,
		TweetID INTEGER NOT NULL,
		UNIQUE (WordID, TweetID)
	);
	CREATE TABLE IF NOT EXISTS Users (
		ID INTEGER NOT NULL,
		Name TEXT NOT NULL,
		ScreenName TEXT NOT NULL,
		Active BOOLEAN NOT NULL DEFAULT 0 CHECK (Active IN (0, 1))
	);
	CREATE TABLE IF NOT EXISTS TwitterCredentials (
		ID INTEGER PRIMARY KEY CHECK (id = 1),
		TWITTER_API_KEY TEXT,
		TWITTER_API_KEY_SECRET TEXT,
		TWITTER_ACCESS_TOKEN TEXT,
		TWITTER_ACCESS_TOKEN_SECRET TEXT
	);
	INSERT OR IGNORE INTO TwitterCredentials (TWITTER_API_KEY, TWITTER_API_KEY_SECRET, TWITTER_ACCESS_TOKEN, TWITTER_ACCESS_TOKEN_SECRET) VALUES ("", "", "", "");
	CREATE TABLE IF NOT EXISTS ShownTweets (
		ID INTEGER NOT NULL PRIMARY KEY,
		TweetID INTEGER NOT NULL,
		ShownAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`

func CreateTables() (error) {
	db, err := sql.Open("sqlite3", getPath())
	if err != nil {
		return err
	}
	defer db.Close()
	if _, err := db.Exec(create); err != nil {
		return err
	}
	return nil
}