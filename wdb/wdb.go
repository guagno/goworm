package wdb

import (
	"errors"
	"io/fs"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const DB_NAME = "/home/marco/Test/worm.db"

func Init_DB() {
	if _, err := os.Stat(DB_NAME); errors.Is(err, fs.ErrNotExist) {
		log.Printf("DB %v doesn't exists, we are going to create it\n", DB_NAME)
		file, err := os.Create(DB_NAME) // Create SQLite file
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
		log.Printf("DB %v created\n", DB_NAME)
	} else {
		log.Printf("DB %v exists\n", DB_NAME)
	}
}
