package qr

import (
	"log"
	"qr-encoder/internal/engraving"
	"qr-encoder/internal/errorcorrection"
	"qr-encoder/internal/masking"
	"qr-encoder/internal/types"
)

type QR struct {
	Size                  int
	Capacity              int
	ErrorCorrection       []uint8
	ErrorCorrectionMarker string
	Formats               map[string]types.FormatData
}

func (c *QR) InitMatrix() types.Matrix {
	result := make([][]bool, c.Size)
	for i := range result {
		result[i] = make([]bool, c.Size)
	}
	return result
}

func (q *QR) X() int {
	return q.Size
}

func (q *QR) Y() int {
	return q.Size
}

func (q *QR) GetErrorCorrectionPolynomial() []byte {
	return q.ErrorCorrection
}

func (q *QR) GetCapacity() int {
	return q.Capacity
}

type mm struct {
	m     types.Matrix
	score int
}

func (q *QR) WriteBitStream(bitStream []bool) (types.Matrix, error) {

	matrix := q.InitMatrix()

	// place data onto matrix
	engraving.WriteDataOntoMatrix(
		matrix,
		q.Size,
		q.Size,
		bitStream,
		func(x int) bool { return x == 6 },
		func(x, y int) bool {
			return x <= 8 && y <= 8 || x <= 8 && y >= q.Size-8 || x >= q.Size-8 && y <= 8 || y == 6
		},
	)

	// evaluate masking patterns
	log.Println("Evaluating mask patterns...")
	maskedMatrixes := make(map[string]*mm)
	for mask, m := range masking.Masks {
		matrix := mm{
			m:     masking.ApplyMask(matrix, m),
			score: 0,
		}
		maskedMatrixes[mask] = &matrix

		// finder patterns
		engraving.WriteSubmatrix(matrix.m, engraving.FinderPatternBackground, 0, 0)
		engraving.WriteSubmatrix(matrix.m, engraving.FinderPattern, 0, 0)

		engraving.WriteSubmatrix(matrix.m, engraving.FinderPatternBackground, q.Size-8, 0)
		engraving.WriteSubmatrix(matrix.m, engraving.FinderPattern, q.Size-7, 0)

		engraving.WriteSubmatrix(matrix.m, engraving.FinderPatternBackground, 0, q.Size-8)
		engraving.WriteSubmatrix(matrix.m, engraving.FinderPattern, 0, q.Size-7)

		formatData := errorcorrection.ComputeFormatErrorCorrection(
			q.ErrorCorrectionMarker,
			mask,
			errorcorrection.FormatBCHPolynomial,
			errorcorrection.FormatMask,
		)
		qrPlaceFormatData(matrix.m, q, formatData)

		// timing pattern
		for i := 8; i < q.Size-8; i += 2 {
			matrix.m[6][i] = true
			matrix.m[6][i+1] = false
			matrix.m[i][6] = true
			matrix.m[i+1][6] = false
		}

		matrix.score = evaluateSymbol(matrix.m)
		log.Printf("Mask: %s, score: %d", mask, matrix.score)
	}

	var min = 9999999
	var result types.Matrix
	var resultMask string
	for mask, matrix := range maskedMatrixes {
		if matrix.score < min {
			min = matrix.score
			result = matrix.m
			resultMask = mask
		}
	}

	log.Printf("Selected mask: %s", resultMask)

	return result, nil
}

func (q *QR) GetFormatData(format string) (bool, *types.FormatData) {
	data, ok := q.Formats[format]
	if !ok {
		return false, nil
	}
	return true, &data
}
