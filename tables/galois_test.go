package tables_test

import (
	"fmt"
	"testing"
	"writer/tables"
)

func Test_Logarithm(t *testing.T) {
	tt := []struct {
		in  int
		out uint8
	}{
		{89, 210}, {61, 228}, {138, 222}, {243, 233}, {149, 184}, {135, 13}, {183, 158}, {59, 120}, {184, 132}, {51, 125}, {41, 147}, {179, 171}, {70, 48}, {84, 143}, {107, 84},
	}
	for _, td := range tt {
		t.Run(
			fmt.Sprintf("%d -> %d", td.in, td.out),
			func(t *testing.T) {
				res := tables.GaloisFieldLogarithm[td.in]
				if res != td.out {
					t.Fatalf("Expected: %d,\nGot: %d\n", td.out, res)
				}
			},
		)
	}
}
