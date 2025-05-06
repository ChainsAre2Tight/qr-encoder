package main

import (
	"fmt"
	"log"
	"os"
	"writer/errorcorrection"
	"writer/format"
	"writer/interfaces"
	"writer/output"
	"writer/qr"
)

var codes = map[string]interfaces.Code{
	"1-M": &qr.QR{
		Size:                  21,
		Capacity:              16,
		ErrorCorrection:       []uint8{0, 251, 67, 46, 61, 118, 70, 64, 94, 32, 45},
		ErrorCorrectionMarker: "00",
	},
}

var formats = map[string]interfaces.Format{
	"byte": &format.ByteFormat{},
}

func main() {
	fail := func(err error) {
		log.Fatal(fmt.Errorf("main: %s", err))
	}

	if l := len(os.Args); l != 2 {
		log.Fatal("Unexpected number of arguments. Expected: 1, got: ", l-1)
	}
	input := os.Args[1]
	log.Println("input is:", input)

	format := "byte"
	f := formats[format]
	log.Println("Selected format:", format)

	data, err := f.Encode(input)
	if err != nil {
		fail(err)
	}

	code := "1-M"
	c := codes[code]
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
