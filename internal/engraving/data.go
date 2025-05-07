package engraving

import (
	"qr-encoder/internal/types"
)

// func PlaceFinderPatterns(matrix types.Matrix, code *qr.QR) {
// 	writeSubmatrix(matrix, FinderPatternBackground, 0, 0)
// 	writeSubmatrix(matrix, FinderPattern, 0, 0)

// 	writeSubmatrix(matrix, FinderPatternBackground, code.Size-8, 0)
// 	writeSubmatrix(matrix, FinderPattern, code.Size-7, 0)

// 	writeSubmatrix(matrix, FinderPatternBackground, 0, code.Size-8)
// 	writeSubmatrix(matrix, FinderPattern, 0, code.Size-7)
// }

// func PlaceTimingPatterns(matrix types.Matrix, code *qr.QR) {
// 	for i := 8; i < code.Size-8; i += 2 {
// 		matrix[6][i] = true
// 		matrix[6][i+1] = false
// 		matrix[i][6] = true
// 		matrix[i+1][6] = false
// 	}
// }

func WriteDataOntoMatrix(
	matrix types.Matrix,
	X, Y int,
	bitstream []bool,
	skipColumn func(x int) bool,
	skipCell func(x, y int) bool,
) {
	up := true
	counter := 0

	// iterate throuth 2-wide columns in reverse order
	for mainX := X - 1; mainX >= 0; mainX -= 2 {
		if skipColumn(mainX) {
			mainX--
		}

		// main counter for Y coordinate, ignores direction
		for mainY := range Y {
			var y int
			// if going up, Y is actually sizeY - mainY - 1 (reverse order)
			if up {
				y = Y - mainY - 1
			} else {
				y = mainY
			}

			for x := mainX; x >= mainX-1; x-- {
				if x < 0 || y < 0 || skipCell(x, y) {
					continue
				}
				if counter >= len(bitstream) {
					return
				}
				matrix[x][y] = bitstream[counter]
				// log.Println(x, y, bitstream[counter])
				counter++
			}
		}
		up = !up
	}
}
