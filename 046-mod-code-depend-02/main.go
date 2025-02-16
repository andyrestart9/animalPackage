package main

import (
	"fmt"

	"github.com/andyrestart9/puppy"
)

func main() {
	s1 := puppy.Bark()
	s2 := puppy.Barks()

	fmt.Println(s1)
	fmt.Println(s2)

	s3 := puppy.BidBark()
	s4 := puppy.BidBarks()

	fmt.Println(s3)
	fmt.Println(s4)

	// also like this
	fmt.Println(puppy.Bark())
	fmt.Println(puppy.Barks())
}
