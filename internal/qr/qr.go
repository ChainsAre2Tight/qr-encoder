package qr

import (
	"fmt"
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

func (q *QR) WriteBitStream(bitStream []bool) (types.Matrix, error) {
	fail := func(err error) (types.Matrix, error) {
		return nil, fmt.Errorf("qr.WriteBitStream: %s", err)
	}

	dataEngraver := &QRDataEngraver{Q: q}
	matrix, mask, err := dataEngraver.Write(bitStream)

	if err != nil {
		return fail(err)
	}

	formatEngraver := &QRMetadataEngraver{q: q}
	formatEngraver.Write(matrix, mask)

	return matrix, nil
}

func (q *QR) GetFormatData(format string) (bool, *types.FormatData) {
	data, ok := q.Formats[format]
	if !ok {
		return false, nil
	}
	return true, &data
}
