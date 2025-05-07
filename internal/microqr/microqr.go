package microqr

import (
	"fmt"
	"qr-encoder/internal/engraving"
	"qr-encoder/internal/errorcorrection"
	"qr-encoder/internal/masking"
	"qr-encoder/internal/types"
)

type MicroQR struct {
	Size                  int
	Capacity              int
	ErrorCorrection       []uint8
	ErrorCorrectionMarker string
}

func (m *MicroQR) X() int {
	return m.Size
}

func (m *MicroQR) Y() int {
	return m.Size
}

func (m *MicroQR) GetCapacity() int {
	return m.Capacity
}

func (m *MicroQR) GetErrorCorrectionPolynomial() []byte {
	return m.ErrorCorrection
}

func (m *MicroQR) WriteBitStream(bitStream []bool) (types.Matrix, error) {
	fail := func(err error) (types.Matrix, error) {
		return nil, fmt.Errorf("microqr.WriteBitStream: %s", err)
	}

	matrix := m.InitMatrix()

	// place data onto matrix
	engraving.WriteDataOntoMatrix(
		matrix,
		m.Size,
		m.Size,
		bitStream,
		func(x int) bool { return x == 0 },
		func(x, y int) bool {
			return x <= 8 && y <= 8 || y == 0
		},
	)

	// TODO: evaluate masking patterns
	mask := "00"
	mk, ok := masking.MicroQRMasks[mask]
	if !ok {
		return fail(fmt.Errorf("mask %s is not found", mask))
	}

	result := masking.ApplyMask(matrix, mk)

	// place finder pattern
	engraving.WriteSubmatrix(result, engraving.FinderPatternBackground, 0, 0)
	engraving.WriteSubmatrix(result, engraving.FinderPattern, 0, 0)

	// place timing pattern
	for i := 8; i < m.Size; i += 2 {
		result[0][i] = true
		result[0][i-1] = false
		result[i][0] = true
		result[i-1][0] = false
	}

	// generate format data
	formatData := errorcorrection.ComputeFormatErrorCorrection(
		m.ErrorCorrectionMarker,
		mask,
		errorcorrection.FormatBCHPolynomial,
		errorcorrection.MicroQRMask,
	)

	// place format data
	for i, pos := range microQRformatPositions {
		result[pos[0]][pos[1]] = formatData[i]
	}

	return result, nil
}

var microQRformatPositions = [15][2]int{
	{1, 8}, {2, 8}, {3, 8}, {4, 8}, {5, 8}, {6, 8}, {7, 8}, {8, 8},
	{8, 7}, {8, 6}, {8, 5}, {8, 4}, {8, 3}, {8, 2}, {8, 1},
}

func (m *MicroQR) InitMatrix() types.Matrix {
	result := make([][]bool, m.Size)
	for i := range result {
		result[i] = make([]bool, m.Size)
	}
	return result
}
