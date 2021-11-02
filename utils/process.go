package utils

import (
	"database/sql"
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"os"

	"github.com/KrisjanisP/viridis/database"
	"golang.org/x/image/tiff"
)

func setFinishedState(rgb int, cir int, ndvi int, overlay int, finished *database.FinishedTile) {
	finished.Rgb = rgb
	finished.Cir = cir
	finished.Ndv = ndvi
	finished.Ove = overlay
}

func checkUpdateFinishedTileRecordRes(cnt int64, res error) {
	if res != nil {
		log.Panic(res)
	}
	if cnt == 0 {
		fmt.Println("Something weird happened")
	}
}
func ProcessTile(tile database.Tile, tileURLs database.TileURLs, dbapi *database.DBAPI) {
	l := log.New(os.Stdout, "[Worker] ", log.Ldate|log.Ltime)

	finished, err := dbapi.GetFinishedTilesRecord(tile.Id)

	if err == sql.ErrNoRows {
		finished.TileId = tile.Id
		setFinishedState(0, 0, 0, 0, &finished)
		dbapi.InsertFinishedTilesRecord(finished)
	} else if err != nil {
		log.Panic(err)
	}

	rgbLoc := GetTileRGBLocation(tile.Name)
	if finished.Rgb == 0 || !FileExists(rgbLoc) {
		setFinishedState(0, 0, 0, 0, &finished)
		cnt, err := dbapi.UpdateFinishedTileRecord(finished)
		checkUpdateFinishedTileRecordRes(cnt, err)
		l.Println("Downloading RGB for " + tile.Name)
		downloadAndConvertTileRGB(tile.Name, tileURLs)
		l.Println("Finished RGB for " + tile.Name)
		setFinishedState(1, 0, 0, 0, &finished)
		cnt, err = dbapi.UpdateFinishedTileRecord(finished)
		checkUpdateFinishedTileRecordRes(cnt, err)
	}

	cirLoc := GetTileCIRLocation(tile.Name)
	if finished.Cir == 0 || !FileExists(cirLoc) {
		setFinishedState(1, 0, 0, 0, &finished)
		cnt, err := dbapi.UpdateFinishedTileRecord(finished)
		checkUpdateFinishedTileRecordRes(cnt, err)
		l.Println("Downloading CIR for " + tile.Name)
		downloadAndConvertTileCIR(tile.Name, tileURLs)
		l.Println("Finished CIR for " + tile.Name)
		setFinishedState(1, 1, 0, 0, &finished)
		cnt, err = dbapi.UpdateFinishedTileRecord(finished)
		checkUpdateFinishedTileRecordRes(cnt, err)
	}

	ndviLoc := GetTileNDVILocation(tile.Name)
	if finished.Ndv == 0 || !FileExists(ndviLoc) {
		setFinishedState(1, 1, 0, 0, &finished)
		cnt, err := dbapi.UpdateFinishedTileRecord(finished)
		checkUpdateFinishedTileRecordRes(cnt, err)
		l.Println("Processing NDVI for " + tile.Name)
		generateTileNDVI(tile.Name)
		l.Println("Finished NDVI for " + tile.Name)
		setFinishedState(1, 1, 1, 0, &finished)
		cnt, err = dbapi.UpdateFinishedTileRecord(finished)
		checkUpdateFinishedTileRecordRes(cnt, err)
	}

	overlayLoc := GetTileOverlayLocation(tile.Name)
	if finished.Ove == 0 || !FileExists(overlayLoc) {
		setFinishedState(1, 1, 1, 0, &finished)
		cnt, err := dbapi.UpdateFinishedTileRecord(finished)
		checkUpdateFinishedTileRecordRes(cnt, err)
		l.Println("Processing overlay for " + tile.Name)
		generateTileOverlay(tile.Name)
		l.Println("Finished overlay for " + tile.Name)
		setFinishedState(1, 1, 1, 1, &finished)
		cnt, err = dbapi.UpdateFinishedTileRecord(finished)
		checkUpdateFinishedTileRecordRes(cnt, err)
	}
}

func ConvertTiffToJPEG(tiff_r io.Reader, jpeg_w io.Writer) {
	//opening files
	img, err := tiff.Decode(tiff_r)

	if err != nil {
		fmt.Println("Cant decode file")
	}

	jpeg.Encode(jpeg_w, img, &jpeg.Options{Quality: 75})
}
