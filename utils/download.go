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
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	// Create the file
	loc := getTileRGBLocation(tileName)
	out, err := os.Create(loc)

	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	ConvertTiffToJPEG(resp.Body, out)
}

func downloadAndConvertTileCIR(tileName string, tileURLs database.TileURLs) {
	// Get the data
	resp, err := http.Get(tileURLs.CirURL)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	// Create the file
	loc := getTileCIRLocation(tileName)
	out, err := os.Create(loc)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	ConvertTiffToJPEG(resp.Body, out)
}
