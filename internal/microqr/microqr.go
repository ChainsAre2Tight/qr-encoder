package microqr

import (
	"log"
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
	Formats               map[string]types.FormatData
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
	// fail := func(err error) (types.Matrix, error) {
	// 	return nil, fmt.Errorf("microqr.WriteBitStream: %s", err)
	// }

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
	maskedMatrixes := make(map[string]types.Matrix)
	for mask, mk := range masking.MicroQRMasks {

		maskedMatrixes[mask] = masking.ApplyMask(matrix, mk)

		// place finder pattern
		engraving.WriteSubmatrix(maskedMatrixes[mask], engraving.FinderPatternBackground, 0, 0)
		engraving.WriteSubmatrix(maskedMatrixes[mask], engraving.FinderPattern, 0, 0)

		// place timing pattern
		for i := 8; i < m.Size; i += 2 {
			maskedMatrixes[mask][0][i] = true
			maskedMatrixes[mask][0][i-1] = false
			maskedMatrixes[mask][i][0] = true
			maskedMatrixes[mask][i-1][0] = false
		}

		// generate format data
		formatData := errorcorrection.ComputeFormatErrorCorrection(
			m.ErrorCorrectionMarker,
			mask,
			errorcorrection.FormatBCHPolynomial,
			errorcorrection.MicroQRMask,
		)

		log.Println(formatData)

		// place format data
		for i, pos := range microQRformatPositions {
			maskedMatrixes[mask][pos[0]][pos[1]] = formatData[i]
		}
	}

	// evaluate masking patterns
	log.Println("Evaluating masking patterns...")
	maskScores := make(map[string]int)
	for mask, matrix := range maskedMatrixes {
		var sum1, sum2 int
		l := len(matrix)
		for i := 1; i < l; i++ {
			if matrix[i][l-1] {
				sum2++
			}
			if matrix[l-1][i] {
				sum1++
			}
		}
		var score int
		if sum1 > sum2 {
			score = sum2*16 + sum1
		} else {
			score = sum1*16 + sum1
		}
		maskScores[mask] = score
		log.Printf("Mask: %s, score: %d", mask, score)
	}

	var max int = 0
	var resultMask string
	var result types.Matrix

	for mask, score := range maskScores {
		if score > max {
			max = score
			result = maskedMatrixes[mask]
			resultMask = mask
		}
	}

	log.Printf("Selected mask: %s", resultMask)

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
func (m *MicroQR) GetFormatData(format string) (bool, *types.FormatData) {
	data, ok := m.Formats[format]
	if !ok {
		return false, nil
	}
	return true, &data
}
