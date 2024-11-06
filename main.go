package main

import (
	"fmt"
	"math/rand/v2"
)

func main() {
	numAs := 1 + rand.IntN(20)  // Random number of A's (1-50)
	numHs := 1 + rand.IntN(100) // Random number of H's (1-50)

	scream := "A"
	for i := 1; i < numAs; i++ {
		scream += "A"
	}
	for i := 0; i < numHs; i++ {
		scream += "H"
	}
	fmt.Println(scream)
}
