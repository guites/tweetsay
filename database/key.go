package database

import (
	"database/sql"
	"log"
)

func AddKey(keyName string, keyValue string)() {
	log.Printf("Saving key @%s to database\n", keyName)

	db, db_err := sql.Open("sqlite3", getPath())
	if db_err != nil {
		log.Fatal("Error opening database", db_err)
	}
	sql_string := "UPDATE TwitterCredentials SET " + keyName + " = ?"
	stmt, stmt_err := db.Prepare(sql_string)

	if stmt_err != nil {
		log.Fatalf((stmt_err.Error()))
	}
	_, err := stmt.Exec(keyValue)

	if err != nil {
		log.Fatal((err.Error()))
	}
	defer stmt.Close()
}