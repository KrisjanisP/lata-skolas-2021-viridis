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

func GetTileRGBLocation(tileName string) string {
	tileFileName := GetTileRGBFileName(tileName)
	return fmt.Sprintf("./data/images/%s", tileFileName)
}

func GetTileCIRLocation(tileName string) string {
	tileFileName := GetTileCIRFileName(tileName)
	return fmt.Sprintf("./data/images/%s", tileFileName)
}

func GetTileNDVILocation(tileName string) string {
	tileFileName := GetTileNDVIFileName(tileName)
	return fmt.Sprintf("./data/images/%s", tileFileName)
}

func GetTileOverlayLocation(tileName string) string {
	tileFileName := GetTileOverlayFileName(tileName)
	return fmt.Sprintf("./data/images/%s_ove.jpeg", tileFileName)
}

func GetImageDirLocation() string {
	return fmt.Sprint("./data/images/")
}
