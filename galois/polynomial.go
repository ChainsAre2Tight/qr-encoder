package galois

import "fmt"

// returns remainder of n / d
func ByteDivRemainder(n, d []uint8) []uint8 {
	r := make([]uint8, len(n)+len(d)-1)
	copy(r, n)

	for lead := range r {
		if lead+len(d) > len(r) {
			break
		}
		t := r[lead]
		for i := range d {
			sub := GaloisMul(d[i], t)
			r[i+lead] ^= sub
		}
		fmt.Println(r)
	}

	return r[len(n):]
}

func BinaryDivRemainder(n, d []bool) []bool {
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
