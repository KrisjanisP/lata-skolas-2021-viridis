package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/KrisjanisP/viridis/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var l *log.Logger

func main() {
	l = log.New(os.Stdout, "[API] ", log.Ldate|log.Ltime)
	os.Mkdir("./data/images", 0755)

	var err error
	db, err = sql.Open("sqlite3", "./data/db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	utils.InitGeoJSON(db)

	go utils.StartWorker(db)

	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.GET("/", serveIndexHTML)
	router.Static("/index", "./front-page")
	router.Static("/images", "./assets/images")
	router.GET("/map.html", serveMapHTML)
	router.GET("/profile.html", serveProfileHTML)
	router.Static("/assets", "./assets") // serve assets like images
	router.Static("/dist", "./dist")     // serve javascript, css
	router.GET("/tiles", getTiles)
	router.POST("/tiles", postTiles)

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	router.Use(cors.New(config))
	router.Run(":8080")
}

func serveIndexHTML(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"index.html",
		gin.H{},
	)
}

func serveMapHTML(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"map.html",
		gin.H{},
	)
}

type FileRow struct {
	TileName string
	FileName string
	DownLink string
	FileType string
	ReqDate  string
}

func serveProfileHTML(c *gin.Context) {
	payload := []FileRow{
		{TileName: "2434-14_4",
			FileName: "4423-43_2_rgb.jpeg",
			FileType: "Krāsainā ortofotokarte",
			ReqDate:  "27.10.2021",
			DownLink: "asdfasdfa"},
		{TileName: "2434-14_4",
			FileName: "4423-43_2_rgb.jpeg",
			FileType: "Krāsainā ortofotokarte",
			ReqDate:  "27.10.2021",
			DownLink: "asdfasdfa"},
	}
	c.HTML(
		http.StatusOK,
		"profile.html",
		gin.H{
			"payload": payload},
	)
}

func getTiles(c *gin.Context) {
	c.File("./data/tiles.geojson")
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

	l.Println("Received tiles: " + strings.Join(tile_names, " "))

	for _, tile_name := range tile_names {
		if tile_name == "" {
			c.AbortWithStatus(400)
			return
		}
		tile, err := utils.GetTileUrlsRecord(db, tile_name)
		if err != nil {
			log.Fatal(err)
			c.AbortWithStatus(500)
			return
		}
		if tile.Name == "" {
			c.AbortWithStatus(400)
			return
		}
		time := time.Now().Format("2006-01-02 15:04:05")
		stmt, err := db.Prepare("INSERT INTO tasks_queue(tile_name, req_date, user_id) values(?,?,?)")
		if err != nil {
			log.Fatal(err)
			c.AbortWithStatus(500)
			return
		}

		_, err = stmt.Exec(tile_name, time, 1)
		if err != nil {
			log.Fatal(err)
			c.AbortWithStatus(500)
			return
		}
	}

	c.IndentedJSON(http.StatusCreated, tile_names)
}
