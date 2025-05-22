package main

import "fmt"

func main() {
	t := newTest("foo")
	st := *t
	fmt.Printf("%#v\n", t)
	fmt.Printf("%#v\n", st)
	fmt.Printf("%#v\n", *t)
	fmt.Printf("%#v\n", (*t).s)
}

func newTest(s string) *testString {
	return &testString{s}
}

type testString struct {
	s string
}
