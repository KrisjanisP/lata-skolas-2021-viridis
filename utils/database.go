package utils

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/KrisjanisP/viridis/models"
)

func GetTileUrlsRecords(db *sql.DB) ([]models.Tile, error) {
	var result []models.Tile
	rows, err := db.Query("SELECT * FROM tile_urls")
	if err != nil {
		log.Fatal(err)
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var tile models.Tile
		err = rows.Scan(&tile.Id, &tile.Name, &tile.TfwURL, &tile.RgbURL, &tile.CirURL)
		if err != nil {
			log.Fatal(err)
			return result, err
		}
		result = append(result, tile)
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return result, err
	}

	return result, nil
}

func GetTileUrlsRecord(db *sql.DB, tile_name string) (models.Tile, error) {
	var result models.Tile
	rows, err := db.Query("SELECT * FROM tile_urls WHERE name = ?", tile_name)
	if err != nil {
		log.Fatal(err)
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		id := result.Id
		fmt.Println(id)
		err = rows.Scan(&result.Id, &result.Name, &result.TfwURL, &result.RgbURL, &result.CirURL)
		if err != nil {
			log.Fatal(err)
			return result, err
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
		return result, err
	}

	return result, nil
}

func GetTaskQueueRecords(db *sql.DB) ([]models.QueueTask, error) {
	var tasks []models.QueueTask
	rows, err := db.Query("SELECT * FROM tasks_queue")

	if err != nil {
		log.Fatal(err)
		return tasks, nil
	}

	for rows.Next() {
		var task models.QueueTask
		err = rows.Scan(&task.Id,
			&task.TileName, &task.ReqDate, &task.UserId)
		if err != nil {
			log.Fatal(err)
			return tasks, nil
		}
		tasks = append(tasks, task)
	}
	rows.Close() //good habit to close
	return tasks, nil
}
