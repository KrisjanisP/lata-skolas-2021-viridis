package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/KrisjanisP/viridis/utils"
	_ "github.com/mattn/go-sqlite3"
)

const (
	databaseFile  = "./data/db.sqlite"
	schemaSQLFile = "./database/schema.sql"
)

func main() {
	if utils.FileExists(databaseFile) {
		fmt.Println("Database file already exists.")
		fmt.Println("Be careful as this script creates a fresh database!")
		fmt.Println("If you wish to create a new database, delete the current one.")
		fmt.Println("If the script encounters an error, the new database will be deleted.")
		fmt.Println("The script is exiting.")
		return
	}
	_, err := os.Create(databaseFile)
	check(err)

	sqlDB, err := sql.Open("sqlite3", databaseFile)
	check(err)
	defer sqlDB.Close()

	schemaSQL, err := os.ReadFile(schemaSQLFile)
	check(err)

	ctx := context.Background()

	tx, err := sqlDB.BeginTx(ctx, nil)
	check(err)
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, string(schemaSQL))
	check(err)

	err = tx.Commit()
	check(err)

	fmt.Println("Database successfuly create.")
}

func check(e error) {
	if e != nil {
		// delete db file if error was encountered
		os.Remove(databaseFile)
		log.Panic(e)
	}
}
