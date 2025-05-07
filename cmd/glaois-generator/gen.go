package main

import "fmt"

func main() {
	// value := 1
	// for range 256 {
	// 	value <<= 1
	// 	if value > 255 {
	// 		value = value ^ 285
	// 	}
	// 	fmt.Printf("%d, ", value)
	// }
	// res := [256]int{}
	// for i, val := range tables.GaloisFieldExponent {
	// 	res[val] = i
	// }
	// for i, val := range res {
	// 	if val == 0 {
	// 		fmt.Println(false, i, val)
	// 	}
	// 	fmt.Printf("%d, ", val)
	// }
	// fmt.Println(true)
	asd := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ $%*+-./:"
	for i, val := range asd {
		fmt.Printf("'%s': %d, ", string(val), i)
	}
}
