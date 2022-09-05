package database

import (
	"log"
	"os"
)

func getPath() (string){
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	file := dirname + "/.tweetsay.db"
	return file
}