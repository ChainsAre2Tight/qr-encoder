package errorcorrection_test

import (
	"fmt"
	"reflect"
	"testing"
	"writer/errorcorrection"
	"writer/qr"
)

func TestCodewordErrorCorrection(t *testing.T) {
	tt := []struct {
		in   []byte
		code *qr.QR
		out  []byte
	}{
		{
			in: []byte{32, 91, 11, 120, 209, 114, 220, 77, 67, 64, 236, 17, 236, 17, 236, 17},
			code: &qr.QR{
				ErrorCorrection: []byte{0, 251, 67, 46, 61, 118, 70, 64, 94, 32, 45},
			},
			out: []byte{196, 35, 39, 119, 235, 215, 231, 226, 93, 23},
		},
	}
	for _, td := range tt {
		t.Run(
			fmt.Sprintf("%v -> %v | %v", td.in, td.out, td.code.ErrorCorrection),
			func(t *testing.T) {
				result := errorcorrection.GenErrorCorrection(td.in, td.code)
				if !reflect.DeepEqual(result, td.out) {
					t.Fatalf("\nExpd: %v,\nGot:  %v,\nCode: %v", td.out, result, td.code.ErrorCorrection)
				}
			},
		)
	}
}
