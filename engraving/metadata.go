package engraving

import (
	"qr-encoder/types"
)

var FinderPattern = types.Matrix{
	{true, true, true, true, true, true, true},
	{true, false, false, false, false, false, true},
	{true, false, true, true, true, false, true},
	{true, false, true, true, true, false, true},
	{true, false, true, true, true, false, true},
	{true, false, false, false, false, false, true},
	{true, true, true, true, true, true, true},
}

var FinderPatternBackground = types.Matrix{
	{false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false},
	{false, false, false, false, false, false, false, false},
}

func WriteSubmatrix(target, data types.Matrix, X, Y int) {
	for x := range data {
		for y := range data[x] {
			target[x+X][y+Y] = data[x][y]
		}
	}
}
