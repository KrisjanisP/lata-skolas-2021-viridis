package utils

import (
	"database/sql"
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"os"

	"github.com/KrisjanisP/viridis/database"
	"github.com/KrisjanisP/viridis/models"
	"golang.org/x/image/tiff"
)

func setFinishedState(rgb int, cir int, ndvi int, overlay int, finished *models.FinishedTile) {
	finished.Rgb = rgb
	finished.Cir = cir
	finished.Ndv = ndvi
	finished.Ove = overlay
}

func ProcessTile(tile models.Tile, dbapi *database.DBAPI) {
	l := log.New(os.Stdout, "[Worker] ", log.Ldate|log.Ltime)

	tileURLs, err := dbapi.GetTileURLsRecord(tile.Id)
	if err != nil {
		log.Panic(err)
	}

	finished, err := dbapi.GetFinishedTilesRecord(tile.Id)

	if err == sql.ErrNoRows {
		finished.TileId = tile.Id
		setFinishedState(0, 0, 0, 0, &finished)
		dbapi.InsertFinishedTilesRecord(finished)
	} else if err != nil {
		log.Panic(err)
	}

	rgbLoc := getTileRGBLocation(tile.Name)
	if finished.Rgb == 0 || !FileExists(rgbLoc) {
		setFinishedState(0, 0, 0, 0, &finished)
		dbapi.UpdateFinishedTileRecord(finished)
		l.Println("Downloading RGB for " + tile.Name)
		downloadAndConvertTileRGB(tile.Name, tileURLs)
		l.Println("Finished RGB for " + tile.Name)
		setFinishedState(1, 0, 0, 0, &finished)
		dbapi.UpdateFinishedTileRecord(finished)
	}

	cirLoc := getTileCIRLocation(tile.Name)
	if finished.Cir == 0 || !FileExists(cirLoc) {
		setFinishedState(1, 0, 0, 0, &finished)
		dbapi.UpdateFinishedTileRecord(finished)
		l.Println("Downloading CIR for " + tile.Name)
		downloadAndConvertTileCIR(tile.Name, tileURLs)
		l.Println("Finished CIR for " + tile.Name)
		setFinishedState(1, 1, 0, 0, &finished)
		dbapi.UpdateFinishedTileRecord(finished)
	}

	ndviLoc := getTileNDVILocation(tile.Name)
	if finished.Ndv == 0 || !FileExists(ndviLoc) {
		setFinishedState(1, 1, 0, 0, &finished)
		dbapi.UpdateFinishedTileRecord(finished)
		l.Println("Processing NDVI for " + tile.Name)
		generateTileNDVI(tile.Name)
		l.Println("Finished NDVI for " + tile.Name)
		setFinishedState(1, 1, 1, 0, &finished)
		dbapi.UpdateFinishedTileRecord(finished)
	}

	overlayLoc := getTileOverlayLocation(tile.Name)
	if finished.Ove == 0 || !FileExists(overlayLoc) {
		setFinishedState(1, 1, 1, 0, &finished)
		dbapi.UpdateFinishedTileRecord(finished)
		l.Println("Processing overlay for " + tile.Name)
		generateTileOverlay(tile.Name)
		l.Println("Finished overlay for " + tile.Name)
		setFinishedState(1, 1, 1, 1, &finished)
		dbapi.UpdateFinishedTileRecord(finished)
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
