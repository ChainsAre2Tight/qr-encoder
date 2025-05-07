package main

import (
	"fmt"
	"log"
	"os"
	"qr-encoder/internal/errorcorrection"
	"qr-encoder/internal/format"
	"qr-encoder/internal/interfaces"
	"qr-encoder/internal/microqr"
	"qr-encoder/internal/output"
	"qr-encoder/internal/qr"
)

var erorrcorrectionpolynomials = map[string][]uint8{
	"10": {0, 251, 67, 46, 61, 118, 70, 64, 94, 32, 45},
}

var codes = map[string]interfaces.Code{
	"1-M": &qr.QR{
		Size:                  21,
		Capacity:              16,
		ErrorCorrection:       erorrcorrectionpolynomials["10"],
		ErrorCorrectionMarker: "00",
	},
	"M4-M": &microqr.MicroQR{
		Size:                  17,
		Capacity:              14,
		ErrorCorrection:       erorrcorrectionpolynomials["10"],
		ErrorCorrectionMarker: "110",
	},
}

var formats = map[string]interfaces.Format{
	"byte":         &format.ByteFormat{},
	"alphanumeric": &format.Alphanumeric{},
}

func main() {
	fail := func(err error) {
		log.Fatal(fmt.Errorf("main: %s", err))
	}

	printUsage := func() {
		log.Fatalf("Usage: format [string] code [string] input [string]")
	}

	if l := len(os.Args); l != 4 {
		log.Println("Unexpected number of arguments. Expected: 4, got: ", l-1)
		printUsage()
	}
	input := os.Args[3]
	log.Println("input is:", input)

	format := os.Args[1]
	f, ok := formats[format]
	if !ok {
		log.Printf("Unsupported format: %s. List of supported formats: %v", format, formats)
		printUsage()
	}
	log.Println("Selected format:", format)

	data, err := f.Encode(input)
	if err != nil {
		fail(err)
	}

	code := os.Args[2]
	c, ok := codes[code]
	if !ok {
		log.Printf("Unsupported code: %s. List of supported codes: %v", code, codes)
		printUsage()
	}
	log.Println("Selected code type:", code)
	log.Println(c)

	bitStream, err := PrepForEngraving(data, c)
	if err != nil {
		fail(err)
	}

	matrix, err := c.WriteBitStream(bitStream)
	if err != nil {
		fail(err)
	}

	output.MatrixToImage(matrix, false)
}

func PrepForEngraving(data []byte, code interfaces.Code) ([]bool, error) {
	fail := func(err error) ([]bool, error) {
		return nil, fmt.Errorf("PrepForEngraving: %s", err)
	}
	flag := true

	if len(data) > code.GetCapacity() {
		return fail(fmt.Errorf("exceeded maximum code capacity (%d > %d)", len(data), code.GetCapacity()))
	}

	// pad data
	for pos := len(data); pos < code.GetCapacity(); pos++ {
		if flag {
			data = append(data, 236)
		} else {
			data = append(data, 17)
		}
		flag = !flag
	}
	log.Println("Data with padding:")
	log.Println(data)

	// generate error correction codes
	fec := errorcorrection.GenErrorCorrection(data, code)

	log.Println("Error correction codes:")
	log.Println(fec)

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

	return bitStream, nil
}
