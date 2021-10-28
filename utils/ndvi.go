package utils

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"os"
)

func GenerateNDVI(jpeg_rgb_r io.Reader, jpeg_cir_r io.Reader, jpeg_ndvi_w io.Writer) {

	//decoding files

	imgNIR, err := jpeg.Decode(jpeg_cir_r)

	if err != nil {
		fmt.Println("Cant decode RGB")
		os.Exit(1)
	}

	imgRGB, err := jpeg.Decode(jpeg_rgb_r)

	if err != nil {
		fmt.Println("Cant decode RGB")
		os.Exit(1)
	}

	rec := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{SIZE, SIZE}})

	JPG := generateJPG(rec, getPixels(imgNIR, imgRGB))

	jpeg.Encode(jpeg_ndvi_w, JPG, &jpeg.Options{Quality: 75})

	imgNIR = nil
	imgRGB = nil
}

func getPixels(imgNIR image.Image, imgRGB image.Image) [][]float32 {
	var pixels [][]float32
	for y := 0; y < SIZE; y++ {
		var row []float32
		for x := 0; x < SIZE; x++ {
			row = append(row, NVDI_calc(imgNIR, imgRGB, x, y))
		}
		pixels = append(pixels, row)
	}
	return pixels
}

func NVDI_calc(imgNIR image.Image, imgRGB image.Image, x int, y int) float32 {
	NIR, _, _, _ := imgNIR.At(x, y).RGBA()
	RED, _, _, _ := imgRGB.At(x, y).RGBA()
	return float32(NIR-RED) / float32(NIR+RED)
}

func generateJPG(JPG *image.RGBA, NDVI [][]float32) *image.RGBA {
	for x := 0; x < SIZE; x++ {
		for y := 0; y < SIZE; y++ {
			v := NDVI[y][x]
			if v < 0 {
				JPG.SetRGBA(x, y, color.RGBA{0, 140 - uint8(200*(v-0.6)), 34 - uint8(85*(v-0.6)), 204})
			} else if v >= 0 && v < 0.2 {
				JPG.SetRGBA(x, y, color.RGBA{200 + uint8(150*v), 60 + uint8(600*v), 0, 204})
			} else if v >= 0.2 && v < 0.3 {
				JPG.SetRGBA(x, y, color.RGBA{230 - uint8(1100*(v-0.2)), 180 + uint8(200*(v-0.2)), 0 + uint8(600*(v-0.2)), 204})
			} else if v >= 0.3 && v < 0.6 {
				JPG.SetRGBA(x, y, color.RGBA{120 - uint8(400*(v-0.3)), 200 - uint8(200*(v-0.3)), 60 - uint8(85*(v-0.3)), 204})
			} else if v > 0.6 {
				JPG.SetRGBA(x, y, color.RGBA{200, 60, 0, 204})
			}
		}
	}
	return JPG
}
