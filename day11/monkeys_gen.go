// Code generated; DO NOT EDIT.

package main

func getInputMonkeys() []Monkey {
	return []Monkey{
		{
			Items: []int{84, 66, 62, 69, 88, 91, 91},
			Op: func(old int) int {
				return old * 11
			},
			Test: 2,
			JT:   4,
			JF:   7,
		},
		{
			Items: []int{98, 50, 76, 99},
			Op: func(old int) int {
				return old * old
			},
			Test: 7,
			JT:   3,
			JF:   6,
		},
		{
			Items: []int{72, 56, 94},
			Op: func(old int) int {
				return old + 1
			},
			Test: 13,
			JT:   4,
			JF:   0,
		},
		{
			Items: []int{55, 88, 90, 77, 60, 67},
			Op: func(old int) int {
				return old + 2
			},
			Test: 3,
			JT:   6,
			JF:   5,
		},
		{
			Items: []int{69, 72, 63, 60, 72, 52, 63, 78},
			Op: func(old int) int {
				return old * 13
			},
			Test: 19,
			JT:   1,
			JF:   7,
		},
		{
			Items: []int{89, 73},
			Op: func(old int) int {
				return old + 5
			},
			Test: 17,
			JT:   2,
			JF:   0,
		},
		{
			Items: []int{78, 68, 98, 88, 66},
			Op: func(old int) int {
				return old + 6
			},
			Test: 11,
			JT:   2,
			JF:   5,
		},
		{
			Items: []int{70},
			Op: func(old int) int {
				return old + 7
			},
			Test: 5,
			JT:   1,
			JF:   3,
		},
	}
}
