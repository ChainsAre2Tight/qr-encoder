package masking

import "writer/types"

var Masks = map[string]types.Mask{
	"000": func(x, y int) bool { return (x+y)%2 == 0 },
	"001": func(x, y int) bool { return y%2 == 0 },
	"010": func(x, y int) bool { return x%3 == 0 },
	"011": func(x, y int) bool { return (x+y)%3 == 0 },
	"100": func(x, y int) bool { return (y/2+x/3)%2 == 0 },
	"101": func(x, y int) bool { return (x*y)%2+(x*y)%3 == 0 },
	"110": func(x, y int) bool { return ((x*y)%2+(x*y)%3)%2 == 0 },
	"111": func(x, y int) bool { return ((x+y)%2+(x*y)%3)%2 == 0 },
}

// Returns a new, masked matrix
func ApplyMask(matrix types.Matrix, mask types.Mask) types.Matrix {
	result := make([][]bool, len(matrix))
	for x := range result {
		result[x] = make([]bool, len(matrix[0]))
		for y := range result[x] {
			result[x][y] = mask(x, y) != matrix[x][y]
		}
	}

	return result
}
