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

type DataEngraver interface {
	Write(bitStream []bool) (types.Matrix, string, error)
}

type MetadataEngraver interface {
	Write(matrix types.Matrix, mask string)
}
