package utils

import (
	"fmt"

	"github.com/KrisjanisP/viridis/database"
)

func StartWorker(dbapi *database.DBAPI) {
	fmt.Println("Hello from your dear worker.")
	/*
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
				ProcessTile(tile, db)
			}
			time.Sleep(time.Second * 3)
		}*/
}
