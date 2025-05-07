package types

type Matrix [][]bool

type Mask func(x, y int) bool

type FormatData struct {
	Indicator string
	CCI       int
	Separator string
}
