package interfaces

import "qr-encoder/internal/types"

type Code interface {
	X() int
	Y() int
	GetCapacity() int
	GetErrorCorrectionPolynomial() []byte
	WriteBitStream(bitStream []bool) (types.Matrix, error)
	GetFormatData(format string) (bool, *types.FormatData)
}

type Format interface {
	Encode(data string, format types.FormatData) ([]byte, error)
}
