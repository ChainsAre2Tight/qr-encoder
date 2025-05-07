package format

import (
	"fmt"
	"strings"
)

type Alphanumeric struct{}

var alphanumeric = map[rune]int{
	'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9, 'A': 10, 'B': 11, 'C': 12, 'D': 13, 'E': 14, 'F': 15, 'G': 16, 'H': 17, 'I': 18, 'J': 19, 'K': 20, 'L': 21, 'M': 22, 'N': 23, 'O': 24, 'P': 25, 'Q': 26, 'R': 27, 'S': 28, 'T': 29, 'U': 30, 'V': 31, 'W': 32, 'X': 33, 'Y': 34, 'Z': 35, ' ': 36, '$': 37, '%': 38, '*': 39, '+': 40, '-': 41, '.': 42, '/': 43, ':': 44,
}

func (f *Alphanumeric) Encode(data string) ([]byte, error) {
	length := len(data)
	var binaryString = ""

	conv := func(r rune) (int, error) {
		n, ok := alphanumeric[r]
		if !ok {
			return 0, fmt.Errorf("convert: Unsupported symbol %s", string(r))
		}
		return n, nil
	}

	for i := 0; i < length/2; i++ {
		first, err := conv(rune(data[2*i]))
		if err != nil {
			return nil, err
		}
		second, err := conv(rune(data[2*i+1]))
		if err != nil {
			return nil, err
		}
		combo := first*45 + second
		str := fmt.Sprintf("%0.11b", combo)
		fmt.Println(str)
		binaryString = binaryString + str
	}
	if length%2 == 1 {
		first, err := conv(rune(data[length-1]))
		if err != nil {
			return nil, err
		}
		str := fmt.Sprintf("%0.5b", first)
		fmt.Println(str)
		binaryString = binaryString + str
	}

	fmt.Println("binary string", binaryString)

	cci := fmt.Sprintf("%0.9b", length)

	fmt.Println("cci", cci)
	binaryString = "0010" + cci + binaryString
	binaryString = binaryString + strings.Repeat("0", 8+8-len(binaryString)%8)

	fmt.Println(binaryString)
	fmt.Println(len(binaryString))

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
