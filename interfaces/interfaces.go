package interfaces

import "writer/types"

type Code interface {
	X() int
	Y() int
	GetCapacity() int
	GetErrorCorrectionPolynomial() []byte
	WriteBitStream(bitStream []bool) (types.Matrix, error)
}

type Format interface {
	Encode(data string) ([]byte, error)
}

type DataEngraver interface {
	Write(bitStream []bool) (types.Matrix, string, error)
}

type MetadataEngraver interface {
	Write(matrix types.Matrix, mask string)
}
