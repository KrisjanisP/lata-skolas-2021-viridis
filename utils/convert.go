package utils

import (
	"fmt"
	"image/jpeg"
	"io"

	"golang.org/x/image/tiff"
)

func ConvertTiffToJPEG(tiff_r io.Reader, jpeg_w io.Writer) {
	//opening files
	img, err := tiff.Decode(tiff_r)

	if err != nil {
		fmt.Println("Cant decode file")
	}

	jpeg.Encode(jpeg_w, img, &jpeg.Options{Quality: 75})
}
