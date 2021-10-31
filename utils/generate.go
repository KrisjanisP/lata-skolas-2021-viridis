package utils

import (
	"log"
	"os"
)

func generateTileNDVI(tileName string) {
	rgbLoc := getTileRGBLocation(tileName)
	rgbFile, err := os.Open(rgbLoc)
	if err != nil {
		log.Fatal(err)
	}
	defer rgbFile.Close()
	cirLoc := getTileCIRLocation(tileName)
	cirFile, err := os.Open(cirLoc)
	if err != nil {
		log.Fatal(err)
	}
	defer cirFile.Close()
	ndviLoc := getTileNDVILocation(tileName)
	ndviFile, err := os.Create(ndviLoc)
	if err != nil {
		log.Fatal(err)
	}
	defer ndviFile.Close()
	GenerateNDVI(rgbFile, cirFile, ndviFile)
}

func generateTileOverlay(tileName string) {
	ndviLoc := getTileNDVILocation(tileName)
	ndviFile, err := os.Open(ndviLoc)
	if err != nil {
		log.Fatal(err)
	}
	defer ndviFile.Close()
	rgbLoc := getTileRGBLocation(tileName)
	rgbFile, err := os.Open(rgbLoc)
	if err != nil {
		log.Fatal(err)
	}
	defer rgbFile.Close()
	overlayLoc := getTileOverlayLocation(tileName)
	overlayFile, err := os.Create(overlayLoc)
	if err != nil {
		log.Fatal(err)
	}
	defer ndviFile.Close()
	GenerateOverlay(ndviFile, rgbFile, overlayFile)
}
