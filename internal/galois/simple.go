package galois

import "qr-encoder/internal/tables"

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

	res := tables.GaloisFieldExponent[tables.GaloisFieldLogarithm[n]-tables.GaloisFieldLogarithm[d]]
	return res
}
