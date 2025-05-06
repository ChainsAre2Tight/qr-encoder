package format

import (
	"fmt"
	"log"

	"golang.org/x/text/encoding/charmap"
)

type ByteFormat struct{}

func (f *ByteFormat) Encode(data string) ([]byte, error) {
	encoder := charmap.ISO8859_1.NewEncoder()
	encodedBytes, err := encoder.Bytes([]byte(data))
	if err != nil {
		return nil, fmt.Errorf("ByteFormat.Encode: %s", err)
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

	// add mode indicator and separator
	binaryData[1] = true

	result := make([]byte, len(binaryData)/8+1)
	pos := 0

	for i := range binaryData {
		pos = i / 8
		if binaryData[i] {
			result[pos] += 1 << (7 - i%8)
		}
	}
	return result, nil
}
