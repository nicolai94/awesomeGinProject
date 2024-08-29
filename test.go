package main

import "fmt"

func newInt() *int {
	x := 100
	return &x
}

func main() {
	p := newInt()
	fmt.Println(*p) // Выведет 100
}
