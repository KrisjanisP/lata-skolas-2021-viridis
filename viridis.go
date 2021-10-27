package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/KrisjanisP/viridis/models"
	process "github.com/KrisjanisP/viridis/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {

	os.Mkdir("./data/images", 0755)

	var err error
	db, err = sql.Open("sqlite3", "./data/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.GET("/", serveHTML)
	router.Static("/assets", "./assets") // serve assets like images
	router.Static("/dist", "./dist")     // serve javascript, css
	router.GET("/tiles", getTiles)
	router.POST("/tiles", postTiles)

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))
	router.Run("127.0.0.1:7080")
}

func serveHTML(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"index.html",
		gin.H{
			"title": "Home Page",
		},
	)
}

func getTiles(c *gin.Context) {
	c.File("./data/tiles.geojson")
}

func getTileUrlsRecord(tile_name string) (models.Tile, error) {
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

func postTiles(c *gin.Context) {
	var tile_names []string

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&tile_names); err != nil {
		return
	}

	if len(tile_names) >= 10 {
		c.AbortWithStatus(400)
		return
	}

	fmt.Println(tile_names)

	for _, tile_name := range tile_names {
		tile, err := getTileUrlsRecord(tile_name)
		if err != nil {
			c.AbortWithStatus(500)
			return
		}
		go process.ProcessTile(tile)
	}

	c.IndentedJSON(http.StatusCreated, tile_names)
}
