package api

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/KrisjanisP/viridis/database"
	"github.com/gin-gonic/gin"
)

type DBAPI struct {
	*database.DBAPI
}

var l *log.Logger

func init() {
	l = log.New(os.Stdout, "[API] ", log.Ldate|log.Ltime)
}

func (dbapi DBAPI) GetTiles(c *gin.Context) {
	l.Println("Sending tiles")
	c.File("./data/tiles.geojson")
}

func (dbapi DBAPI) PostTiles(c *gin.Context) {
	var tileNames []string

	if err := c.BindJSON(&tileNames); err != nil {
		return
	}

	if len(tileNames) >= 10 {
		c.AbortWithStatus(400)
		return
	}

	l.Println("Received tiles: " + strings.Join(tileNames, " "))
	for _, tileName := range tileNames {
		_, err := dbapi.GetTileId(tileName)
		if err == sql.ErrNoRows {
			c.AbortWithStatus(400)
		} else if err != nil {
			c.AbortWithStatus(500)
			return
		}
		//time := time.Now().Format("2006-01-02 15:04:05")

	}
	c.IndentedJSON(http.StatusCreated, tileNames)
}
