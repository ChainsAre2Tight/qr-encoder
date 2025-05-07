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
	fail := func(err error) ([]byte, error) {
		return nil, fmt.Errorf("ByteFormat.Encode: %s", err)
	}

	encoder := charmap.ISO8859_1.NewEncoder()
	encodedBytes, err := encoder.Bytes([]byte(data))
	if err != nil {
		return fail(err)
	}
	length := len(encodedBytes)
	log.Println("Encoded data is")
	log.Println(encodedBytes)

	binaryString := ""
	for _, b := range encodedBytes {
		add, err := DecimalToBinaryString(int(b), 8)
		if err != nil {
			return fail(err)
		}
		binaryString = binaryString + add
	}

	// add content length indicator
	cci, err := DecimalToBinaryString(length, format.CCI)
	if err != nil {
		return fail(err)
	}

	// add mode indicator and separator
	log.Printf("Resulting string: %s %s %s %s", format.Indicator, cci, binaryString, format.Separator)
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
