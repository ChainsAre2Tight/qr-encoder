package format

import (
	"log"
	"strconv"
	"strings"
)

func DecimalToBinaryString(number int, length int) string {
	raw := strconv.FormatInt(int64(number), 2)
	padded := strings.Repeat("0", length-len(raw)) + raw
	log.Println("--->", padded)
	return padded
}
