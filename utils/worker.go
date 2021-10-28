package utils

import (
	"database/sql"
	"log"
	"os"
	"time"
)

func StartWorker(db *sql.DB) {
	for {
		l := log.New(os.Stdout, "[Worker] ", log.Ldate|log.Ltime)
		l.Println("Checking for tasks.")
		tasks, err := GetTaskQueueRecords(db)
		if err != nil {
			log.Fatal(err)
		}
		for _, task := range tasks {
			tile, err := GetTileUrlsRecord(db, task.TileName)
			if err != nil {
				log.Fatal(err)
			}
			// process tile if necessary
			ProcessTile(tile)
		}
		time.Sleep(time.Second * 3)
	}
}
