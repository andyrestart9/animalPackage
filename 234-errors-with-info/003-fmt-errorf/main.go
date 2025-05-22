package main

import (
	"fmt"
	"log"
	"math"
)

func main() {
	_, err := sqrt(-10)
	if err != nil {
		log.Fatalln(err)
	}
}

func sqrt(f float64) (float64, error) {
	if f < 0 {
		ErrWrongMath:=fmt.Errorf("wrong math: square root of negative number: %v",f)
		return 0, ErrWrongMath
	}
	return math.Sqrt(f), nil
}
