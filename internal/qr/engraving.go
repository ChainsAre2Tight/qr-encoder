package qr

import (
	"qr-encoder/internal/types"
)

func qrPlaceFormatData(matrix types.Matrix, code *QR, formatData []bool) {
	formatPositionsUpperLeft := [15][2]int{
		{0, 8}, {1, 8}, {2, 8}, {3, 8}, {4, 8}, {5, 8},
		{7, 8}, {8, 8}, {8, 7},
		{8, 5}, {8, 4}, {8, 3}, {8, 2}, {8, 1}, {8, 0},
	}
	x := code.Size
	y := code.Size
	formatPositionsLowerRight := [15][2]int{
		{8, y - 1}, {8, y - 2}, {8, y - 3}, {8, y - 4}, {8, y - 5}, {8, y - 6}, {8, y - 7},
		{x - 8, 8}, {x - 7, 8}, {x - 6, 8}, {x - 5, 8}, {x - 4, 8}, {x - 3, 8}, {x - 2, 8}, {x - 1, 8},
	}

	for i := range formatData {
		posUL := formatPositionsUpperLeft[i]
		matrix[posUL[0]][posUL[1]] = formatData[i]

		posLR := formatPositionsLowerRight[i]
		matrix[posLR[0]][posLR[1]] = formatData[i]
	}

	// always black
	matrix[8][y-8] = true
}
