package main

import (
	"errors"
	"fmt"
	"log"
	"math"
)

var ErrWrongMath = errors.New("wrong math: square root of negative number")

func main() {
	fmt.Printf("%T\n", ErrWrongMath)
	_, err := sqrt(-10)
	if err != nil {
		log.Fatalln(err)
	}
}

func sqrt(f float64) (float64, error) {
	if f < 0 {
		return 0, ErrWrongMath
	}
	return math.Sqrt(f), nil
}
