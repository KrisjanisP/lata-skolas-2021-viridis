package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/KrisjanisP/viridis/database"
	"github.com/KrisjanisP/viridis/utils"
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

func (dbapi DBAPI) DownloadFile(ctx *gin.Context) {
	session := sessions.Default(ctx)
	profile := session.Get("profile")
	m, ok := profile.(map[string]interface{})
	if !ok {
		ctx.AbortWithStatus(400)
		return
	}
	sub := m["sub"]
	userid := sub.(string)
	tileidStr := ctx.Param("tileid")
	id, err := strconv.Atoi(tileidStr)
	tileId := int64(id)
	if err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	tileName, err := dbapi.GetTileName(tileId)
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatus(500)
		return
	}
	fileType := ctx.Param("type")
	if fileType != "rgb" && fileType != "cir" && fileType != "ndv" && fileType != "ove" {
		ctx.AbortWithStatus(400)
		return
	}
	tilePossesion, err := dbapi.SelectPossesionRecordByUserId(tileId, userid)
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatus(400)
		return
	}
	if tilePossesion.UserId != userid || tilePossesion.TileId != tileId {
		ctx.AbortWithStatus(400)
		return
	}

	finishedTile, err := dbapi.GetFinishedTilesRecord(tileId)
	if err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	fileLoc := utils.GetTileRGBLocation(tileName)
	switch fileType {
	case "rgb":
		if finishedTile.Rgb == 0 {
			ctx.AbortWithStatus(400)
			return
		} else {
			fileLoc = utils.GetTileRGBLocation(tileName)
		}
	case "cir":
		if finishedTile.Cir == 0 {
			ctx.AbortWithStatus(400)
			return
		} else {
			fileLoc = utils.GetTileCIRLocation(tileName)
		}
	case "ndv":
		if finishedTile.Ndv == 0 {
			ctx.AbortWithStatus(400)
			return
		} else {
			fileLoc = utils.GetTileNDVILocation(tileName)
		}
	case "ove":
		if finishedTile.Ove == 0 {
			ctx.AbortWithStatus(400)
			return
		} else {
			fileLoc = utils.GetTileOverlayLocation(tileName)
		}
	}
	l.Println("Sending " + fileLoc)
	ctx.File(fileLoc)
	fmt.Println("request is correct")
}
