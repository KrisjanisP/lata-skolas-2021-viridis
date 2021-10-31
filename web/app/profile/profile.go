package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type FileRow struct {
	TileName string
	FileName string
	DownLink string
	FileType string
	ReqDate  string
}

// Handler for our logged-in user page.
func Handler(ctx *gin.Context) {
	/*
		finisheds, err := utils.GetFinishedTilesRecords(db)
		if err != nil {
			log.Fatal(err)
		}
	*/
	var payload []FileRow
	/*
		for _, finished := range finisheds {
			var file FileRow
			file.TileName = finished.Name
			file.FileName = utils.GetTileRGBFileName(finished.Name)
			file.DownLink = "1234"
			file.ReqDate = "27.10.2021"
			file.FileType = "Krāsainā ortofotokarte"
			payload = append(payload, file)
		}
	*/

	ctx.HTML(
		http.StatusOK,
		"profile.html",
		gin.H{
			"payload": payload},
	)
}
