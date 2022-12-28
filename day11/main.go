package main

import (
	"fmt"
	"sort"
)

//go:generate go run ./gen -i input.txt -o monkeys_gen.go

func main() {
	{
		monkeys := getInputMonkeys()
		numMonkeys := len(monkeys)
		counts := make([]int, numMonkeys)
		for i := 0; i < 20; i++ {
			for n, i := range monkeys {
				for _, j := range i.Items {
					k, t := i.Process1(j)
					monkeys[t].Add(k)
				}
				counts[n] += len(i.Items)
				i.Discard()
			}
		}
		sort.Ints(counts)
		fmt.Println("Part 1:", counts[numMonkeys-1]*counts[numMonkeys-2])
	}
	{
		monkeys := getInputMonkeys()
		modulus := 1
		for _, i := range monkeys {
			modulus *= i.Test
		}
		numMonkeys := len(monkeys)
		counts := make([]int, numMonkeys)
		for i := 0; i < 10000; i++ {
			for n, i := range monkeys {
				for _, j := range i.Items {
					k, t := i.Process2(j)
					monkeys[t].Add(k % modulus)
				}
				counts[n] += len(i.Items)
				i.Discard()
			}
		}
		sort.Ints(counts)
		fmt.Println("Part 2:", counts[numMonkeys-1]*counts[numMonkeys-2])
	}
}

type (
	Monkey struct {
		Items []int
		Op    func(old int) int
		Test  int
		JT    int
		JF    int
	}
)

func (m *Monkey) Discard() {
	if len(m.Items) != 0 {
		m.Items = m.Items[:0]
	}
}

func (m *Monkey) Add(val int) {
	m.Items = append(m.Items, val)
}

func (m *Monkey) Process1(val int) (int, int) {
	k := m.Op(val) / 3
	if k%m.Test == 0 {
		return k, m.JT
	} else {
		return k, m.JF
	}
}

func (m *Monkey) Process2(val int) (int, int) {
	k := m.Op(val)
	if k%m.Test == 0 {
		return k, m.JT
	} else {
		return k, m.JF
	}
}
