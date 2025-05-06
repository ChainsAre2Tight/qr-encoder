package output

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"writer/types"
)

func MatrixToImage(matrix types.Matrix) {
	upperLeft, lowerRight := image.Point{0, 0}, image.Point{len(matrix), len(matrix[0])}
	img := image.NewGray(image.Rectangle{upperLeft, lowerRight})

	for x := range matrix {
		for y := range matrix[0] {
			if matrix[x][y] {
				img.Set(x, y, color.Black)
			} else {
				img.Set(x, y, color.White)
			}
		}
	}

	file, err := os.Create("output.png")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	png.Encode(file, img)
}
