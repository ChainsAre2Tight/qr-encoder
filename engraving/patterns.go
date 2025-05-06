package engraving

import (
	"writer/qr"
	"writer/types"
)

var FinderPattern = types.Matrix{
	{true, true, true, true, true, true, true},
	{true, false, false, false, false, false, true},
	{true, false, true, true, true, false, true},
	{true, false, true, true, true, false, true},
	{true, false, true, true, true, false, true},
	{true, false, false, false, false, false, true},
	{true, true, true, true, true, true, true},
}

var FinderPatternBackground = types.Matrix{
	{false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false},
}

func writeSubmatrix(target, data types.Matrix, X, Y int) {
	for x := range data {
		for y := range data[x] {
			target[x+X][y+Y] = data[x][y]
		}
	}
}

func PlaceFinderPatterns(matrix types.Matrix, code *qr.QR) {
	writeSubmatrix(matrix, FinderPatternBackground, 0, 0)
	writeSubmatrix(matrix, FinderPattern, 0, 0)

	writeSubmatrix(matrix, FinderPatternBackground, code.Size-8, 0)
	writeSubmatrix(matrix, FinderPattern, code.Size-7, 0)

	writeSubmatrix(matrix, FinderPatternBackground, 0, code.Size-8)
	writeSubmatrix(matrix, FinderPattern, 0, code.Size-7)
}

func PlaceTimingPatterns(matrix types.Matrix, code *qr.QR) {
	for i := 8; i < code.Size-8; i += 2 {
		matrix[6][i] = true
		matrix[6][i+1] = false
		matrix[i][6] = true
		matrix[i+1][6] = false
	}
}
