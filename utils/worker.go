package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/KrisjanisP/viridis/database"
)

func StartWorker(dbapi *database.DBAPI) {
	fmt.Println("Hello from your dear worker.")
	for {
		l := log.New(os.Stdout, "[Worker] ", log.Ldate|log.Ltime)
		l.Println("Checking for tasks.")
		tiles, tileURLsArr, err := dbapi.JoinPossesionRecordTileIdsNamesURLs()
		if err != nil {
			log.Panic(err)
		}
		for i := 0; i < len(tiles); i++ {
			//fmt.Println(tiles[i].Name)
			//fmt.Println(tileURLsArr[i].RgbURL)
			ProcessTile(tiles[i], tileURLsArr[i], dbapi)
		}
		fmt.Println()
		time.Sleep(time.Second * 3)
	}
}
