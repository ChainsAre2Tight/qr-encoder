package errorcorrection

import (
	"qr-encoder/internal/galois"
	"qr-encoder/internal/interfaces"
	"qr-encoder/internal/tables"
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

var FormatBCHPolynomial = []bool{true, false, true, false, false, true, true, false, true, true, true}
var FormatMask = []bool{
	true, false, true, false, true,
	false, false, false, false, false,
	true, false, false, true, false,
}
var MicroQRMask = []bool{
	true, false, false, false, true,
	false, false, false, true, false,
	false, false, true, false, true,
}

func ComputeFormatErrorCorrection(level, mask string, polynomial []bool, formatMask []bool) []bool {
	combined := level + mask

	result := make([]bool, 15)
	for i := range 5 {
		result[i] = combined[i] == '1'
	}

	correction := galois.BinaryDivRemainder(result, polynomial)
	for i := 5; i < 15; i++ {
		result[i] = correction[i]
	}

	// apply masking
	for i := range 15 {
		result[i] = result[i] != formatMask[i]
	}

	return result
}
