package main 

import "fmt"
import "math/big"
import "numthy"

func main() {
	fmt.Println("Hello, from numthy/main.")
	x, _ := new(big.Int).SetString("123", 10)
	y, _ := new(big.Int).SetString("456", 10)
	N, _ := new(big.Int).SetString("7890", 10)
	p := new(big.Int)

	numthy.MultMod(x, y, p, N)
	fmt.Printf("Product is now %v\n", p)

	fmt.Printf("Test PowerMod_RLB")
	g, _ := new(big.Int).SetString("2", 10)
	n, _ := new(big.Int).SetString("20", 10)
	numthy.PowerMod_RLB(g, n, p, N)
	fmt.Printf("2^20 (7890) = %v\n", p)

	fmt.Printf("Bye.\n")
}