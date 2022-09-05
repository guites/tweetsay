package helper

import (
	"log"
	"os"
)

func getDbPath() (string){
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	file := dirname + "/.tweetsay.db"
	return file
}