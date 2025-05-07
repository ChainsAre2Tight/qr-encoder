package format

import (
	"fmt"
	"log"
	"qr-encoder/internal/types"
	"strings"
)

type Alphanumeric struct{}

var alphanumericTable = map[rune]int{
	'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9, 'A': 10, 'B': 11, 'C': 12, 'D': 13, 'E': 14, 'F': 15, 'G': 16, 'H': 17, 'I': 18, 'J': 19, 'K': 20, 'L': 21, 'M': 22, 'N': 23, 'O': 24, 'P': 25, 'Q': 26, 'R': 27, 'S': 28, 'T': 29, 'U': 30, 'V': 31, 'W': 32, 'X': 33, 'Y': 34, 'Z': 35, ' ': 36, '$': 37, '%': 38, '*': 39, '+': 40, '-': 41, '.': 42, '/': 43, ':': 44,
}

func (f *Alphanumeric) Encode(data string, format types.FormatData) ([]byte, error) {
	length := len(data)
	var binaryString = ""

	conv := func(r rune) (int, error) {
		n, ok := alphanumericTable[r]
		if !ok {
			return 0, fmt.Errorf("convert: Unsupported symbol %s", string(r))
		}
		return n, nil
	}

	log.Println("converting to alphanumeric...")

	for i := 0; i < length/2; i++ {
		r1 := rune(data[2*i])
		first, err := conv(r1)
		if err != nil {
			return nil, err
		}
		r2 := rune(data[2*i+1])
		second, err := conv(r2)
		if err != nil {
			return nil, err
		}
		combo := first*45 + second
		str := fmt.Sprintf("%0.11b", combo)
		log.Printf("%s | %s ---> %s", string(r1), string(r2), str)
		binaryString = binaryString + str
	}
	if length%2 == 1 {
		r := rune(data[length-1])
		first, err := conv(r)
		if err != nil {
			return nil, err
		}
		str := fmt.Sprintf("%0.5b", first)
		log.Printf("%s ---> %s", string(r), str)
		binaryString = binaryString + str
	}

	cci := DecimalToBinaryString(length, format.CCI)
	log.Println("cci", cci)

	log.Printf("Resulting string: %s %s %s %s", format.Indicator, cci, binaryString, format.Separator)
	binaryString = format.Indicator + cci + binaryString + format.Separator

	padding := strings.Repeat("0", 8-len(binaryString)%8)
	binaryString = binaryString + padding

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
