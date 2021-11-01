package utils

import (
	"fmt"
)

func GetTileRGBFileName(tileName string) string {
	return fmt.Sprintf("%s_rgb.jpeg", tileName)
}

func GetTileCIRFileName(tileName string) string {
	return fmt.Sprintf("%s_cir.jpeg", tileName)
}

func GetTileNDVIFileName(tileName string) string {
	return fmt.Sprintf("%s_ndvi.jpeg", tileName)
}

func GetTileOverlayFileName(tileName string) string {
	return fmt.Sprintf("%s_overlay.jpeg", tileName)
}

func getTileRGBLocation(tileName string) string {
	tileFileName := GetTileRGBFileName(tileName)
	return fmt.Sprintf("./data/images/%s", tileFileName)
}

func getTileCIRLocation(tileName string) string {
	tileFileName := GetTileCIRFileName(tileName)
	return fmt.Sprintf("./data/images/%s", tileFileName)
}

func getTileNDVILocation(tileName string) string {
	tileFileName := GetTileNDVIFileName(tileName)
	return fmt.Sprintf("./data/images/%s", tileFileName)
}

func getTileOverlayLocation(tileName string) string {
	tileFileName := GetTileOverlayFileName(tileName)
	return fmt.Sprintf("./data/images/%s_ove.jpeg", tileFileName)
}

func getImageDirLocation() string {
	return fmt.Sprint("./data/images/")
}
