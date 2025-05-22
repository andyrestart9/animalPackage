package main

import (
	"fmt"

	"github.com/andyrestart9/puppy"
)

// go get github.com/andyrestart9/puppy
// go mod tidy
func main() {
	s1 := puppy.Bark()
	s2 := puppy.Barks()

	fmt.Println(s1)
	fmt.Println(s2)

	// also like this
	fmt.Println(puppy.Bark())
	fmt.Println(puppy.Barks())
}
