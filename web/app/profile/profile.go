package profile

import (
	"fmt"
	"log"
	"net/http"

	"github.com/KrisjanisP/viridis/database"
	"github.com/KrisjanisP/viridis/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type FileRow struct {
	TileName   string
	FileName   string
	DownLink   string
	FileType   string
	ReqDate    string
	IsFinished bool
}

type DBAPI struct {
	*database.DBAPI
}

// Handler for our logged-in user page.
func (dbapi DBAPI) Handler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	profile := session.Get("profile")
	m, ok := profile.(map[string]interface{})
	if !ok {
		ctx.AbortWithStatus(400)
	}
	sub := m["sub"]
	var payload []FileRow
	tiles, finishedTiles, err := dbapi.JoinPossesionRecordTileIdsNamesFinished(sub.(string))
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(len(tiles))
	for i := 0; i < len(tiles); i++ {
		var fileRow FileRow
		fileRow.TileName = tiles[i].Name
		fileRow.IsFinished = finishedTiles[i].Rgb == 1
		fileRow.FileName = utils.GetTileRGBFileName(tiles[i].Name)
		fileRow.DownLink = fmt.Sprintf("/download/%d?type=rgb", tiles[i].Id)
		fileRow.FileType = "krāsainā ortofotokarte"
		payload = append(payload, fileRow)
		fileRow.IsFinished = finishedTiles[i].Cir == 1
		fileRow.FileName = utils.GetTileCIRFileName(tiles[i].Name)
		fileRow.DownLink = fmt.Sprintf("/download/%d?type=cir", tiles[i].Id)
		fileRow.FileType = "infrasarkanā ortofotokarte"
		payload = append(payload, fileRow)
		fileRow.IsFinished = finishedTiles[i].Ndv == 1
		fileRow.FileName = utils.GetTileNDVIFileName(tiles[i].Name)
		fileRow.DownLink = fmt.Sprintf("/download/%d?type=ndv", tiles[i].Id)
		fileRow.FileType = "apstrādātā ortofotokarte"
		payload = append(payload, fileRow)
		fileRow.IsFinished = finishedTiles[i].Ove == 1
		fileRow.FileName = utils.GetTileOverlayFileName(tiles[i].Name)
		fileRow.DownLink = fmt.Sprintf("/download/%d?type=ove", tiles[i].Id)
		fileRow.FileType = "apstrādātā ar pārklājumu"
		payload = append(payload, fileRow)
	}
	fmt.Println(sub)
	ctx.HTML(
		http.StatusOK,
		"profile.html",
		gin.H{
			"payload": payload,
			"profile": profile},
	)
}
