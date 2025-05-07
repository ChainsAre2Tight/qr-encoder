package qr

import (
	"fmt"
	"qr-encoder/internal/engraving"
	"qr-encoder/internal/errorcorrection"
	"qr-encoder/internal/masking"
	"qr-encoder/internal/types"
)

type QRDataEngraver struct {
	Q *QR
}

func (e *QRDataEngraver) Write(bitStream []bool) (types.Matrix, string, error) {
	fail := func(err error) (types.Matrix, string, error) {
		return nil, "", fmt.Errorf("qrDataEngraver.Write: %s", err)
	}

	matrix := e.Q.InitMatrix()

	// place data onto matrix
	engraving.WriteDataOntoMatrix(
		matrix,
		e.Q.Size,
		e.Q.Size,
		bitStream,
		func(x int) bool { return x == 6 },
		func(x, y int) bool {
			return x <= 8 && y <= 8 || x <= 8 && y >= e.Q.Size-8 || x >= e.Q.Size-8 && y <= 8 || y == 6
		},
	)

	// TODO: evaluate masking patterns
	mask := "101"
	m, ok := masking.Masks[mask]
	if !ok {
		fail(fmt.Errorf("mask %s is not found", mask))
	}

	result := masking.ApplyMask(matrix, m)

	return result, mask, nil
}

type QRMetadataEngraver struct {
	q *QR
}

func (e *QRMetadataEngraver) Write(matrix types.Matrix, mask string) error {

	if len(mask) != 3 {
		return fmt.Errorf("invalid mask length: %d (%s)", len(mask), mask)
	}

	// finder patterns
	engraving.WriteSubmatrix(matrix, engraving.FinderPatternBackground, 0, 0)
	engraving.WriteSubmatrix(matrix, engraving.FinderPattern, 0, 0)

	engraving.WriteSubmatrix(matrix, engraving.FinderPatternBackground, e.q.Size-8, 0)
	engraving.WriteSubmatrix(matrix, engraving.FinderPattern, e.q.Size-7, 0)

	engraving.WriteSubmatrix(matrix, engraving.FinderPatternBackground, 0, e.q.Size-8)
	engraving.WriteSubmatrix(matrix, engraving.FinderPattern, 0, e.q.Size-7)

	formatData := errorcorrection.ComputeFormatErrorCorrection(e.q.ErrorCorrectionMarker, mask)
	qrPlaceFormatData(matrix, e.q, formatData)

	// timing pattern
	for i := 8; i < e.q.Size-8; i += 2 {
		matrix[6][i] = true
		matrix[6][i+1] = false
		matrix[i][6] = true
		matrix[i+1][6] = false
	}

	return nil
}

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
