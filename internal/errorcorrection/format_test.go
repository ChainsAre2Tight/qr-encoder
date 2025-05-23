package errorcorrection_test

import (
	"fmt"
	"qr-encoder/internal/errorcorrection"
	"testing"
)

func TestFormatErrorCorrectionGeneration(t *testing.T) {
	tt := []struct {
		mode string
		mask string
		out  string
	}{
		{
			mode: "11",
			mask: "000",
			out:  "011010101011111",
		}, {
			mode: "00",
			mask: "101",
			out:  "100000011001110",
		}, {
			mode: "00",
			mask: "010",
			out:  "101111001111100",
		},
	}
	for _, td := range tt {
		t.Run(
			fmt.Sprintf("%s | %s -> %s", td.mode, td.mask, td.out),
			func(t *testing.T) {
				raw_result := errorcorrection.ComputeFormatErrorCorrection(td.mode, td.mask, errorcorrection.FormatBCHPolynomial, errorcorrection.FormatMask)
				res := ""
				for _, val := range raw_result {
					if val {
						res = res + "1"
					} else {
						res = res + "0"
					}
				}
				if res != td.out {
					t.Fatalf("\nWant: %s,\nGot:  %s.", td.out, res)
				}
			},
		)
	}
}

func TestMicroQrFormat(t *testing.T) {
	tt := []struct {
		mode string
		mask string
		out  string
	}{
		{
			mode: "110",
			mask: "00",
			out:  "010010100001000",
		},
	}
	for _, td := range tt {
		t.Run(
			fmt.Sprintf("%s | %s -> %s", td.mode, td.mask, td.out),
			func(t *testing.T) {
				raw_result := errorcorrection.ComputeFormatErrorCorrection(td.mode, td.mask, errorcorrection.FormatBCHPolynomial, errorcorrection.MicroQRMask)
				res := ""
				for _, val := range raw_result {
					if val {
						res = res + "1"
					} else {
						res = res + "0"
					}
				}
				if res != td.out {
					t.Fatalf("\nWant: %s,\nGot:  %s.", td.out, res)
				}
			},
		)
	}
}
