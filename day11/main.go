package main

import (
	"fmt"
)

//go:generate go run ./gen -o monkeys_gen.go

func main() {
	fmt.Println("hello")
}

type (
	Monkey struct {
		Items []int
		Op    func(old int) int
		Test  int
		JT    int
		JF    int
		count int
	}
)
