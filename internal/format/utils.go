package format

import (
	"fmt"
	"strconv"
	"strings"
)

func DecimalToBinaryString(number int, length int) (string, error) {
	raw := strconv.FormatInt(int64(number), 2)
	if len(raw) > length {
		return "", fmt.Errorf("decimalToBinaryString: Recieved number whose binary representation is longer than exepected output (%d > %d)", len(raw), length)
	}
	padded := strings.Repeat("0", length-len(raw)) + raw
	return padded, nil
}
