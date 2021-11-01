package api

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/KrisjanisP/viridis/database"
	"github.com/gin-contrib/sessions"
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

func (dbapi DBAPI) PostTiles(ctx *gin.Context) {
	var tileNames []string
	if err := ctx.BindJSON(&tileNames); err != nil {
		return
	}
	if len(tileNames) >= 10 {
		ctx.AbortWithStatus(400)
		return
	}

	session := sessions.Default(ctx)
	profile := session.Get("profile")
	m, ok := profile.(map[string]interface{})
	if !ok {
		ctx.AbortWithStatus(400)
		return
	}
	sub := m["sub"]
	if sub == nil {
		ctx.AbortWithStatus(400)
		return
	}

	l.Println("Received tiles: " + strings.Join(tileNames, " "))
	var tileIds []int64
	for _, tileName := range tileNames {
		tileId, err := dbapi.GetTileId(tileName)
		if err == sql.ErrNoRows {
			log.Println(err)
			ctx.AbortWithStatus(400)
			return
		} else if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(500)
			return
		}
		tileIds = append(tileIds, tileId)
	}
	var tilePossesions []database.TilePossesion
	for _, tileId := range tileIds {
		tilePossesion := database.TilePossesion{TileId: tileId, UserId: sub.(string)}
		tilePossesions = append(tilePossesions, tilePossesion)
	}
	err := dbapi.InsertOrIgnoreTilePossesionRecords(tilePossesions)
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatus(400)
	}
	ctx.IndentedJSON(http.StatusCreated, tileNames)
}
