package qr

import "writer/types"

type QR struct {
	Size                  int
	Capacity              int
	ErrorCorrection       []uint8
	ErrorCorrectionMarker string
	// FormatCorrectionCode  []uint8
}

func (c *QR) initMatrix() types.Matrix {
	result := make([][]bool, c.Size)
	for i := range result {
		result[i] = make([]bool, c.Size)
	}
	return result
}

func (q *QR) WriteDataOntoMatrix(
	bitstream []bool,
	skipColumn func(x int) bool,
	skipCell func(x, y int) bool,
) types.Matrix {
	up := true

	matrix := q.initMatrix()
	// return matrix
	counter := 0

	// iterate throuth 2-wide columns in reverse order
	for mainX := q.Size - 1; mainX >= 0; mainX -= 2 {
		if skipColumn(mainX) {
			mainX--
		}

		// main counter for Y coordinate, ignores direction
		for mainY := range q.Size {
			var y int
			// if going up, Y is actually sizeY - mainY - 1 (reverse order)
			if up {
				y = q.Size - mainY - 1
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
