package format

import (
	"fmt"
	"log"
	"qr-encoder/internal/types"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

type ByteFormat struct{}

func (f *ByteFormat) Encode(data string, format types.FormatData) ([]byte, error) {
	encoder := charmap.ISO8859_1.NewEncoder()
	encodedBytes, err := encoder.Bytes([]byte(data))
	if err != nil {
		return nil, fmt.Errorf("ByteFormat.Encode: %s", err)
	}
	length := len(encodedBytes)
	log.Println("Encoded data is")
	log.Println(encodedBytes)

	binaryString := ""
	for _, b := range encodedBytes {
		binaryString = binaryString + DecimalToBinaryString(int(b), 8)
	}

	// add content length indicator
	cci := DecimalToBinaryString(length, format.CCI)

	// add mode indicator and separator
	binaryString = format.Indicator + cci + binaryString + format.Separator
	binaryString = binaryString + strings.Repeat("0", 8-len(binaryString)%8)

	result := make([]byte, len(binaryString)/8)
	pos := 0

	for i := range binaryString {
		pos = i / 8
		if binaryString[i] == '1' {
			result[pos] += 1 << (7 - i%8)
		}
	}
	return result, nil
}
