package qr

import "qr-encoder/internal/types"

func evaluateSymbol(matrix types.Matrix) int {
	size := len(matrix)

	var n1 int

	var blackCounter, whiteCounter int
	for x := range size {
		flag := false
		counter := 0
		for y := range size {
			if matrix[x][y] {
				blackCounter++
				if flag {
					counter++
				} else {
					if counter >= 5 {
						n1 += counter - 2
						counter = 1
						flag = true
					}
				}
			} else {
				whiteCounter++
				if !flag {
					counter++
				} else if counter >= 5 {
					n1 += counter - 2
					counter = 1
					flag = false
				}
			}
		}
	}
	for y := range size {
		flag := false
		counter := 0
		for x := range size {
			if matrix[x][y] {
				if flag {
					counter++
				} else {
					if counter >= 5 {
						n1 += counter - 2
						counter = 1
						flag = true
					}
				}
			} else {
				if !flag {
					counter++
				} else if counter >= 5 {
					n1 += counter - 2
					counter = 1
					flag = false
				}
			}
		}
	}

	blackPercent := float64(blackCounter) / float64(blackCounter+whiteCounter) * 100
	var n4 int
	if blackPercent > 50.0 {
		blackPercent -= 50.0
		n4 = (int(blackPercent) % 5) * 10
	}

	return n1 + n4
}
