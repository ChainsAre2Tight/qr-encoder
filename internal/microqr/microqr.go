package microqr

import "qr-encoder/internal/types"

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
	return nil, nil
}

func (m *MicroQR) InitMatrix() types.Matrix {
	result := make([][]bool, m.Size)
	for i := range result {
		result[i] = make([]bool, m.Size)
	}
	return result
}
