package utils

import (
	"fmt"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/KrisjanisP/viridis/models"
	"golang.org/x/image/tiff"
)

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func getTileRGBLocation(tile models.Tile) string {
	return fmt.Sprintf("./data/images/%s_rgb.jpeg", tile.Name)
}

func getTileCIRLocation(tile models.Tile) string {
	return fmt.Sprintf("./data/images/%s_cir.jpeg", tile.Name)
}

func getTileNDVILocation(tile models.Tile) string {
	return fmt.Sprintf("./data/images/%s_ndvi.jpeg", tile.Name)
}

func getTileOverlayLocation(tile models.Tile) string {
	return fmt.Sprintf("./data/images/%s_overlay.jpeg", tile.Name)
}

func downloadAndConvertTileRGB(tile models.Tile) {
	// Get the data
	resp, err := http.Get(tile.RgbURL)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(getTileRGBLocation(tile))

	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	ConvertTiffToJPEG(resp.Body, out)
}

func downloadAndConvertTileCIR(tile models.Tile) {
	// Get the data
	resp, err := http.Get(tile.CirURL)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(getTileCIRLocation(tile))
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	ConvertTiffToJPEG(resp.Body, out)
}

func generateTileNDVI(tile models.Tile) {
	rgbLoc := getTileRGBLocation(tile)
	rgbFile, err := os.Open(rgbLoc)
	if err != nil {
		log.Fatal(err)
	}
	defer rgbFile.Close()
	cirLoc := getTileCIRLocation(tile)
	cirFile, err := os.Open(cirLoc)
	if err != nil {
		log.Fatal(err)
	}
	defer cirFile.Close()
	ndviLoc := getTileNDVILocation(tile)
	ndviFile, err := os.Create(ndviLoc)
	if err != nil {
		log.Fatal(err)
	}
	defer ndviFile.Close()
	GenerateNDVI(rgbFile, cirFile, ndviFile)
}

func generateTileOverlay(tile models.Tile) {
	ndviLoc := getTileNDVILocation(tile)
	ndviFile, err := os.Open(ndviLoc)
	if err != nil {
		log.Fatal(err)
	}
	defer ndviFile.Close()
	rgbLoc := getTileRGBLocation(tile)
	rgbFile, err := os.Open(rgbLoc)
	if err != nil {
		log.Fatal(err)
	}
	defer rgbFile.Close()
	overlayLoc := getTileOverlayLocation(tile)
	overlayFile, err := os.Create(overlayLoc)
	if err != nil {
		log.Fatal(err)
	}
	defer ndviFile.Close()
	GenerateOverlay(ndviFile, rgbFile, overlayFile)
}

func ProcessTile(tile models.Tile) {
	fmt.Println("Processing " + tile.Name)

	rgbLoc := getTileRGBLocation(tile)
	if !fileExists(rgbLoc) {
		fmt.Println("Downloading RGB for " + tile.Name)
		downloadAndConvertTileRGB(tile)
		fmt.Println("Finished RGB for " + tile.Name)
	}

	cirLoc := getTileCIRLocation(tile)
	if !fileExists(cirLoc) {
		fmt.Println("Downloading CIR for " + tile.Name)
		downloadAndConvertTileCIR(tile)
		fmt.Println("Finished CIR for " + tile.Name)
	}

	ndviLoc := getTileNDVILocation(tile)
	if !fileExists(ndviLoc) {
		fmt.Println("Processing NDVI for " + tile.Name)
		generateTileNDVI(tile)
		fmt.Println("Finished NDVI for " + tile.Name)
	}

	overlayLoc := getTileOverlayLocation(tile)
	if !fileExists(overlayLoc) {
		fmt.Println("Processing overlay for " + tile.Name)
		generateTileOverlay(tile)
		fmt.Println("Finished overlay for " + tile.Name)
	}
	fmt.Println("Finished processing " + tile.Name)
}

func ConvertTiffToJPEG(tiff_r io.Reader, jpeg_w io.Writer) {
	//opening files
	img, err := tiff.Decode(tiff_r)

	if err != nil {
		fmt.Println("Cant decode file")
	}

	jpeg.Encode(jpeg_w, img, &jpeg.Options{Quality: 75})
}
