package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"writer/tables"

	"golang.org/x/text/encoding/charmap"
)

type Matrix [][]bool

type QR struct {
	Size                  int
	Capacity              int
	ErrorCorrection       []uint8
	ErrorCorrectionMarker string
	// FormatCorrectionCode  []uint8
}

func (code *QR) initMatrix() Matrix {
	result := make([][]bool, code.Size)
	for i := range result {
		result[i] = make([]bool, code.Size)
	}
	return result
}

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

	MatrixToImage(matrix)
}

func GaloisMul(a, b uint8) uint8 {
	// fmt.Println(a, b,
	// 	tables.GaloisFieldLogarithm[a], tables.GaloisFieldLogarithm[b],
	// 	(tables.GaloisFieldLogarithm[a] + tables.GaloisFieldLogarithm[b]),
	// 	tables.GaloisFieldExponent[tables.GaloisFieldLogarithm[a]+tables.GaloisFieldLogarithm[b]],
	// )
	if a == 0 || b == 0 {
		return 0
	}
	var res uint8

	// account for integer overflow
	dif := 255 - tables.GaloisFieldLogarithm[a]
	if tables.GaloisFieldLogarithm[b] > dif {
		res = tables.GaloisFieldExponent[tables.GaloisFieldLogarithm[a]+tables.GaloisFieldLogarithm[b]+1]
	} else {
		res = tables.GaloisFieldExponent[tables.GaloisFieldLogarithm[a]+tables.GaloisFieldLogarithm[b]]
	}

	return res
}

func GaloisDiv(n, d uint8) uint8 {
	if n == 0 {
		return 0
	}
	if d == 0 {
		panic("division by zero")
	}

	return tables.GaloisFieldExponent[tables.GaloisFieldLogarithm[n]-tables.GaloisFieldLogarithm[d]]
}

// returns remainder of n / d
func PolynomialGaloisDivRemainder(n, d []uint8) []uint8 {
	r := make([]uint8, len(n))
	copy(r, n)

	for lead := range r {
		if lead+len(d) > len(r) {
			break
		}
		t := GaloisDiv(r[lead], d[0])
		for i := range d {
			sub := GaloisMul(d[i], t)
			r[i+lead] -= sub
		}
	}

	return r[len(n)-len(d)+1:]
}

func GenErrorCorrection(b []byte, code *QR) []byte {
	divisor := make([]uint8, len(code.ErrorCorrection))
	for i, power := range code.ErrorCorrection {
		divisor[i] = tables.GaloisFieldExponent[power]
	}

	result := PolynomialGaloisDivRemainder(b, divisor)

	return result
}

func WriteDataOntoMatrix(code *QR, bitstream []bool, skipColumn func(x int) bool, skipCell func(x, y int) bool) Matrix {
	up := true

	matrix := code.initMatrix()
	counter := 0

	// iterate throuth 2-wide columns in reverse order
	for mainX := code.Size - 1; mainX >= 0; mainX -= 2 {
		if skipColumn(mainX) {
			mainX--
		}

		// main counter for Y coordinate, ignores direction
		for mainY := range code.Size {
			var y int
			// if going up, Y is actually sizeY - mainY - 1 (reverse order)
			if up {
				y = code.Size - mainY - 1
			} else {
				y = mainY
			}

			for x := mainX; x >= mainX-1; x-- {
				if x < 0 || y < 0 || skipCell(x, y) {
					continue
				}
				if counter >= len(bitstream) {
					return matrix
				}
				matrix[x][y] = bitstream[counter]
				counter++
			}
		}
		up = !up
	}
	return matrix
}

func Encode(input string) (Matrix, error) {
	fail := func(err error) (Matrix, error) {
		return nil, fmt.Errorf("Encode: %s", err)
	}
	// code := &QR{
	// 	Size:            29,
	// 	Capacity:        28,
	// 	ErrorCorrection: []uint8{0, 120, 104, 107, 109, 102, 161, 76, 3, 91, 191, 147, 169, 182, 194, 225, 120},
	// }

	code := &QR{
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

	// for i := range binaryData {
	// 	pos = i / 8
	// 	if binaryData[i] {
	// 		data[pos] += 1 << (7 - i%8)
	// 	}
	// }

	asd := []byte{
		16, 32, 12, 86, 97, 128,
	}
	for pos = range asd {
		data[pos] = asd[pos]
	}

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
	fec := GenErrorCorrection(data, code)

	log.Println("Error correction codes:")
	log.Println(fec)

	testFEC := []byte{
		165, 36, 212, 193, 237, 54, 199, 135, 44, 85,
	}
	log.Println(testFEC)

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

	// log.Println(bitStream)

	// place data onto matrix
	matrix := WriteDataOntoMatrix(
		code,
		bitStream,
		func(x int) bool { return x == 6 },
		func(x, y int) bool {
			return x <= 8 && y <= 8 || x <= 8 && y >= code.Size-8 || x >= code.Size-8 && y <= 8
		},
	)

	// evaluate masking patterns
	mask := "010"
	result := ApplyMask(matrix, Masks[mask])

	// place format data and its error corrections
	PlaceFinderPatterns(result, code)

	formatData := ComputeFormatErrorCorrection(code.ErrorCorrectionMarker, mask)
	log.Println("Format data:")
	log.Println(formatData)

	PlaceFormatData(result, code, formatData)
	PlaceTimingPatterns(result, code)

	return result, nil
}

func PlaceFormatData(matrix Matrix, code *QR, formatData []bool) {
	formatPositionsUpperLeft := [15][2]int{
		{0, 8}, {1, 8}, {2, 8}, {3, 8}, {4, 8}, {5, 8},
		{7, 8}, {8, 8}, {8, 7},
		{8, 5}, {8, 4}, {8, 3}, {8, 2}, {8, 1}, {8, 0},
	}
	x := code.Size
	y := code.Size
	formatPositionsLowerRight := [15][2]int{
		{8, y - 1}, {8, y - 2}, {8, y - 3}, {8, y - 4}, {8, y - 5}, {8, y - 6}, {8, y - 7},
		{x - 8, 8}, {x - 7, 8}, {x - 6, 8}, {x - 5, 8}, {x - 4, 8}, {x - 3, 8}, {x - 2, 8}, {x - 1, 8},
	}

	for i := range formatData {
		posUL := formatPositionsUpperLeft[i]
		matrix[posUL[0]][posUL[1]] = formatData[i]

		posLR := formatPositionsLowerRight[i]
		matrix[posLR[0]][posLR[1]] = formatData[i]
	}

	matrix[8][y-8] = true
}

func PlaceFinderPatterns(matrix Matrix, code *QR) {
	WriteSubmatrix(matrix, FinderPatternBackground, 0, 0)
	WriteSubmatrix(matrix, FinderPattern, 0, 0)

	WriteSubmatrix(matrix, FinderPatternBackground, code.Size-8, 0)
	WriteSubmatrix(matrix, FinderPattern, code.Size-7, 0)

	WriteSubmatrix(matrix, FinderPatternBackground, 0, code.Size-8)
	WriteSubmatrix(matrix, FinderPattern, 0, code.Size-7)
}

func PlaceTimingPatterns(matrix Matrix, code *QR) {
	for i := 8; i < code.Size-8; i += 2 {
		matrix[6][i] = true
		matrix[6][i+1] = false
		matrix[i][6] = true
		matrix[i+1][6] = false
	}
}

func BinaryPolynomyalGaloisDivRemainder(n, d []bool) []bool {
	r := make([]bool, len(n))
	copy(r, n)

	for lead := range r {
		if lead+len(d) > len(r) {
			break
		}
		t := r[lead] == d[0]
		if !t {
			continue
		}
		for i := range d {
			if d[i] {
				r[lead+i] = !r[lead+i]
			}
		}
	}
	return r
}

var FormatBCHPolynomial = []bool{true, false, true, false, false, true, true, false, true, true, true}
var FormatMask = []bool{
	true, false, true, false, true,
	false, false, false, false, false,
	true, false, false, true, false,
}

func ComputeFormatErrorCorrection(level, mask string) []bool {
	combined := level + mask

	result := make([]bool, 15)
	for i := range 5 {
		result[i] = combined[i] == '1'
	}

	correction := BinaryPolynomyalGaloisDivRemainder(result, FormatBCHPolynomial)
	for i := 5; i < 15; i++ {
		result[i] = correction[i]
	}

	// apply masking
	for i := range 15 {
		result[i] = result[i] != FormatMask[i]
	}

	return result
}

type Mask func(x, y int) bool

func ApplyMask(matrix Matrix, mask Mask) Matrix {
	result := make([][]bool, len(matrix))
	for x := range result {
		result[x] = make([]bool, len(matrix[0]))
		for y := range result[x] {
			result[x][y] = mask(x, y) != matrix[x][y]
		}
	}

	return result
}

var Masks = map[string]Mask{
	"000": func(x, y int) bool { return (x+y)%2 == 0 },
	"001": func(x, y int) bool { return y%2 == 0 },
	"010": func(x, y int) bool { return x%3 == 0 },
	"011": func(x, y int) bool { return (x+y)%3 == 0 },
	"100": func(x, y int) bool { return (y/2+x/3)%2 == 0 },
	"101": func(x, y int) bool { return (x*y)%2+(x*y)%3 == 0 },
	"110": func(x, y int) bool { return ((x*y)%2+(x*y)%3)%2 == 0 },
	"111": func(x, y int) bool { return ((x+y)%2+(x*y)%3)%2 == 0 },
}

func MatrixToImage(matrix Matrix) {
	upperLeft, lowerRight := image.Point{0, 0}, image.Point{len(matrix), len(matrix[0])}
	img := image.NewGray(image.Rectangle{upperLeft, lowerRight})

	for x := range matrix {
		for y := range matrix[0] {
			if matrix[x][y] {
				img.Set(x, y, color.Black)
			} else {
				img.Set(x, y, color.White)
			}
		}
	}

	file, err := os.Create("output.png")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	png.Encode(file, img)
}

var FinderPattern = Matrix{
	{true, true, true, true, true, true, true},
	{true, false, false, false, false, false, true},
	{true, false, true, true, true, false, true},
	{true, false, true, true, true, false, true},
	{true, false, true, true, true, false, true},
	{true, false, false, false, false, false, true},
	{true, true, true, true, true, true, true},
}

var FinderPatternBackground = Matrix{
	{false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false},
}

func WriteSubmatrix(target, data Matrix, X, Y int) {
	for x := range data {
		for y := range data[x] {
			target[x+X][y+Y] = data[x][y]
		}
	}
}
