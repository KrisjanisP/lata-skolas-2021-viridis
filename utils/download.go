package utils

import (
	"log"
	"net/http"
	"os"

	"github.com/KrisjanisP/viridis/database"
)

func downloadAndConvertTileRGB(tileName string, tileURLs database.TileURLs) {
	// Get the data
	resp, err := http.Get(tileURLs.RgbURL)
	if err != nil {
		log.Panic(err)
		return
	}
	defer resp.Body.Close()
	// Create Image dir
	dir := GetImageDirLocation()
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		log.Panic(err)
	}
	// Create the file
	loc := GetTileRGBLocation(tileName)
	out, err := os.Create(loc)
	if err != nil {
		log.Panic(err)
	}
	defer out.Close()

	ConvertTiffToJPEG(resp.Body, out)
}

func downloadAndConvertTileCIR(tileName string, tileURLs database.TileURLs) {
	// Get the data
	resp, err := http.Get(tileURLs.CirURL)
	if err != nil {
		log.Panic(err)
		return
	}
	defer resp.Body.Close()

	// Create the file
	loc := GetTileCIRLocation(tileName)
	out, err := os.Create(loc)
	if err != nil {
		log.Panic(err)
	}
	defer out.Close()

	ConvertTiffToJPEG(resp.Body, out)
}
