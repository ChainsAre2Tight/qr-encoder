package errorcorrection

import (
	"writer/galois"
	"writer/interfaces"
	"writer/tables"
)

func GenErrorCorrection(b []byte, code interfaces.Code) []byte {
	// lenWords := len(b)
	// for range len(code.ErrorCorrection) {
	// 	b = append(b, 0)
	// }

	divisor := make([]uint8, len(code.GetErrorCorrectionPolynomial()))
	for i, power := range code.GetErrorCorrectionPolynomial() {
		divisor[i] = tables.GaloisFieldExponent[power]
	}

	result := galois.ByteDivRemainder(b, divisor)

	return result
}

var formatBCHPolynomial = []bool{true, false, true, false, false, true, true, false, true, true, true}
var formatMask = []bool{
	true, false, true, false, true,
	false, false, false, false, false,
	true, false, false, true, false,
}

func ComputeFormatErrorCorrection(level, mask string) []bool {
	combined := level + mask

	result := make([]bool, 15)
	for i := range 5 {
		result[i] = combined[i] == '1'
	}

	correction := galois.BinaryDivRemainder(result, formatBCHPolynomial)
	for i := 5; i < 15; i++ {
		result[i] = correction[i]
	}

	// apply masking
	for i := range 15 {
		result[i] = result[i] != formatMask[i]
	}

	return result
}
