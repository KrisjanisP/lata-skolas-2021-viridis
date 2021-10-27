package utils

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"os"
)

const SIZE = 10000

func GenerateOverlay(jpeg_ndvi_r io.Reader, jpeg_rgb_r io.Reader, jpeg_over_w io.Writer) {
	imgBACK, err := jpeg.Decode(jpeg_rgb_r)

	if err != nil {
		fmt.Println("Cant decode rgb")
		os.Exit(1)
	}

	imgNDVI, err := jpeg.Decode(jpeg_ndvi_r)

	if err != nil {
		fmt.Println("Cant decode ndvi")
		os.Exit(1)
	}

	rec := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{10000, 10000}})

	JPG := generateOverlayJPG(rec, getOverlayPixels(imgBACK, imgNDVI))

	jpeg.Encode(jpeg_over_w, JPG, &jpeg.Options{Quality: 75})
}

func getOverlayPixels(imgBACK image.Image, imgNDVI image.Image) [][][]uint8 {
	var pixels [][][]uint8
	for y := 0; y < 10000; y++ {
		var row [][]uint8
		for x := 0; x < 10000; x++ {
			row = append(row, overlayCalc(imgBACK, imgNDVI, x, y))
		}
		pixels = append(pixels, row)
	}
	return pixels
}

func overlay(lower float32, higher float32) uint8 {
	return uint8((lower / 255) * (lower + (((2.5 * higher) / 255) * (255 - lower))))
}

func overlayCalc(imgBACK image.Image, imgNDVI image.Image, x int, y int) []uint8 {
	Br, Bg, Bb, _ := imgBACK.At(x, y).RGBA()
	Nr, Ng, Nb, _ := imgNDVI.At(x, y).RGBA()
	var pixel []uint8
	pixel = append(pixel, overlay(float32(Br/257), float32(Nr/257)))
	pixel = append(pixel, overlay(float32(Bg/257), float32(Ng/257)))
	pixel = append(pixel, overlay(float32(Bb/257), float32(Nb/257)))
	return pixel
}

func generateOverlayJPG(JPG *image.RGBA, NewLayer [][][]uint8) *image.RGBA {
	for x := 0; x < SIZE; x++ {
		for y := 0; y < SIZE; y++ {
			JPG.SetRGBA(x, y, color.RGBA{NewLayer[y][x][0], NewLayer[y][x][1], NewLayer[y][x][2], 255})
		}
	}
	return JPG
}
