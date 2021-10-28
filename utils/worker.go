package utils

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func StartWorker(db *sql.DB) {
	for {
		fmt.Println("Checking for tasks.")
		tasks, err := GetTaskQueueRecords(db)
		if err != nil {
			log.Fatal(err)
			continue
		}
		fmt.Println()
		for _, task := range tasks {
			fmt.Println(task.TileName)
			tile, err := GetTileUrlsRecord(db, task.TileName)
			if err != nil {
				log.Fatal(err)
				continue
			}
			ProcessTile(tile)
		}
		fmt.Println()
		time.Sleep(time.Second * 2)
	}
}
