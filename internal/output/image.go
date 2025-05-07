package output

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"qr-encoder/internal/types"

	"github.com/nfnt/resize"
)

func MatrixToImage(matrix types.Matrix, include_border bool) {
	var upperLeft, lowerRight image.Point
	if include_border {
		upperLeft, lowerRight = image.Point{-4, -4}, image.Point{len(matrix) + 4, len(matrix[0]) + 4}
	} else {
		upperLeft, lowerRight = image.Point{0, 0}, image.Point{len(matrix), len(matrix[0])}
	}

	img := image.NewGray(image.Rectangle{upperLeft, lowerRight})

	if include_border {
		for x := -4; x < len(matrix)+4; x++ {
			for y := -4; y < len(matrix)+4; y++ {
				img.Set(x, y, color.White)
			}
		}
	}

	for x := range matrix {
		for y := range matrix[0] {
			if matrix[x][y] {
				img.Set(x, y, color.Black)
			} else {
				img.Set(x, y, color.White)
			}
		}
	}

	newImage := resize.Resize(500, 500, img, resize.NearestNeighbor)

	file, err := os.Create("output.png")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	png.Encode(file, newImage)
}
