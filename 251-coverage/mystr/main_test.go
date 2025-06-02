package mystr

import (
	"fmt"
	"testing"
)

func TestGreet(t *testing.T) {
	s := Greet("James")
	if s != "Hello my dear, James" {
		t.Error("got", s, "want", "Hello my dear, James")
	}
}

func ExampleGreet() {
	fmt.Println(Greet("James"))
	// Output:
	// Hello my dear, James
}
