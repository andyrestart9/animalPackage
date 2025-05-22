// Package mymath provides ACME inc math solutions.
package mymath

// Sum adds an unlimited number of values of types int.
func Sum(xi ...int) int {
	sum := 0
	for _, v := range xi {
		sum += v
	}
	return sum
}
