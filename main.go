package main

import (
	"fmt"
	"log"
	"os"
	"writer/engraving"
	"writer/errorcorrection"
	"writer/masking"
	"writer/output"
	"writer/qr"
	"writer/types"

	"golang.org/x/text/encoding/charmap"
)

func main() {
	if l := len(os.Args); l != 2 {
		log.Fatal("Unexpected number of arguments. Expected: 1, got: ", l-1)
	}
	input := os.Args[1]
	log.Println("input is:", input)

	matrix, err := Encode(input)
	if err != nil {
		log.Fatal(fmt.Errorf("main: %s", err))
	}

	output.MatrixToImage(matrix, false)
}

func Encode(input string) (types.Matrix, error) {
	fail := func(err error) (types.Matrix, error) {
		return nil, fmt.Errorf("Encode: %s", err)
	}
	// code := &QR{
	// 	Size:            29,
	// 	Capacity:        28,
	// 	ErrorCorrection: []uint8{0, 120, 104, 107, 109, 102, 161, 76, 3, 91, 191, 147, 169, 182, 194, 225, 120},
	// }

	code := &qr.QR{
		Size:                  21,
		Capacity:              16,
		ErrorCorrection:       []uint8{0, 251, 67, 46, 61, 118, 70, 64, 94, 32, 45},
		ErrorCorrectionMarker: "00",
		// FormatCorrectionCode:  []uint8{26, 16, 4},
	}

	log.Println("Code:")
	log.Printf("    Size: %d,\n", code.Size)
	log.Printf("    Capacity: %d,\n", code.Capacity)
	log.Println("    ErrorCorrection:", code.ErrorCorrection)

	// encode data into binary
	encoder := charmap.ISO8859_1.NewEncoder()
	encodedBytes, err := encoder.Bytes([]byte(input))
	if err != nil {
		return fail(err)
	}
	log.Println("Encoded data is")
	log.Println(encodedBytes)

	// add content length indicator
	encodedData := make([]byte, len(encodedBytes)+1)
	encodedData[0] = byte(len(encodedBytes))
	for i, val := range encodedBytes {
		encodedData[i+1] = val
	}

	binaryData := make([]bool, 8*(len(encodedData)+1))
	for i, val := range encodedData {
		start := i*8 + 4
		for j := range 8 {
			if val&(1<<j) > 0 {
				binaryData[start+7-j] = true
			}
		}
	}
	// log.Println(binaryData)

	// add mode indicator and separator
	// binaryData[0] = true
	binaryData[1] = true
	// binaryData[2] = true
	// binaryData[3] = true
	// log.Println(binaryData)

	data := make([]byte, code.Capacity)
	pos := 0

	for i := range binaryData {
		pos = i / 8
		if binaryData[i] {
			data[pos] += 1 << (7 - i%8)
		}
	}

	// asd := []byte{
	// 	16, 32, 12, 86, 97, 128,
	// }
	// for pos = range asd {
	// 	data[pos] = asd[pos]
	// }

	if pos < code.Capacity-1 {
		pos++
		flag := true
		for ; pos < code.Capacity; pos++ {
			if flag {
				data[pos] = 236
			} else {
				data[pos] = 17
			}
			flag = !flag
		}
	}

	log.Println("Data with padding:")
	log.Println(data)

	// generate error correction codes
	fec := errorcorrection.GenErrorCorrection(data, code)

	log.Println("Error correction codes:")
	log.Println(fec)

	// testFEC := []byte{
	// 	165, 36, 212, 193, 237, 54, 199, 135, 44, 85,
	// }
	// log.Println(testFEC)

	// convert data and fec to bit stream
	bitStream := make([]bool, 8*(len(data)+len(fec)))
	start := -8
	for _, val := range data {
		start += 8
		for j := range 8 {
			if val&(1<<j) > 0 {
				bitStream[start+7-j] = true
			}
		}
	}
	for _, val := range fec {
		start += 8
		for j := range 8 {
			if val&(1<<j) > 0 {
				bitStream[start+7-j] = true
			}
		}
	}

	data = make([]byte, len(bitStream)/8)
	for i := range bitStream {
		pos = i / 8
		if bitStream[i] {
			data[pos] += 1 << (7 - i%8)
		}
	}

	log.Println("bitstream\n", data)

	// place data onto matrix
	matrix := code.WriteDataOntoMatrix(
		bitStream,
		func(x int) bool { return x == 6 },
		func(x, y int) bool {
			return x <= 8 && y <= 8 || x <= 8 && y >= code.Size-8 || x >= code.Size-8 && y <= 8 || y == 6
		},
	)

	// evaluate masking patterns
	mask := "101"
	result := masking.ApplyMask(matrix, masking.Masks[mask])
	// result := matrix

	// place format data and its error corrections
	engraving.PlaceFinderPatterns(result, code)

	formatData := errorcorrection.ComputeFormatErrorCorrection(code.ErrorCorrectionMarker, mask)
	log.Println("Format data:")
	log.Println(formatData)

	engraving.PlaceFormatData(result, code, formatData)
	engraving.PlaceTimingPatterns(result, code)

	return result, nil
}
