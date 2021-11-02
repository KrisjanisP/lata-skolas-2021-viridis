package utils

import (
	"log"
	"os"
)

func generateTileNDVI(tileName string) {
	rgbLoc := GetTileRGBLocation(tileName)
	rgbFile, err := os.Open(rgbLoc)
	if err != nil {
		log.Panic(err)
	}
	defer rgbFile.Close()
	cirLoc := GetTileCIRLocation(tileName)
	cirFile, err := os.Open(cirLoc)
	if err != nil {
		log.Panic(err)
	}
	defer cirFile.Close()
	ndviLoc := GetTileNDVILocation(tileName)
	ndviFile, err := os.Create(ndviLoc)
	if err != nil {
		log.Panic(err)
	}
	defer ndviFile.Close()
	GenerateNDVI(rgbFile, cirFile, ndviFile)
}

func generateTileOverlay(tileName string) {
	ndviLoc := GetTileNDVILocation(tileName)
	ndviFile, err := os.Open(ndviLoc)
	if err != nil {
		log.Panic(err)
	}
	defer ndviFile.Close()
	rgbLoc := GetTileRGBLocation(tileName)
	rgbFile, err := os.Open(rgbLoc)
	if err != nil {
		log.Panic(err)
	}
	defer rgbFile.Close()
	overlayLoc := GetTileOverlayLocation(tileName)
	overlayFile, err := os.Create(overlayLoc)
	if err != nil {
		log.Panic(err)
	}
	defer ndviFile.Close()
	GenerateOverlay(ndviFile, rgbFile, overlayFile)
}
