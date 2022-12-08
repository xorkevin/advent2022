package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
)

const (
	puzzleInput = "input.txt"
)

func main() {
	file, err := os.Open(puzzleInput)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	var cells [][]int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		row := make([]int, 0, len(line))
		for _, i := range line {
			if i < '0' || i > '9' {
				log.Fatalln("Invalid grid cell")
			}
			row = append(row, int(i)-'0')
		}
		cells = append(cells, row)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	grid, err := NewGrid(cells)
	if err != nil {
		log.Fatalln(err)
	}
	grid.computeVisibleSet()

	fmt.Println("Part 1:", len(grid.visible))
	fmt.Println("Part 2:", grid.maxPower())
}

type (
	Tuple2 struct {
		y, x int
	}

	Tuple4 struct {
		t, r, b, l int
	}

	Grid struct {
		w, h    int
		grid    [][]int
		visible map[Tuple2]struct{}
		power   map[Tuple2]*Tuple4
	}
)

func NewGrid(grid [][]int) (*Grid, error) {
	h := len(grid)
	if h == 0 {
		return nil, errors.New("Empty grid")
	}
	w := len(grid[0])
	for _, i := range grid {
		if len(i) != w {
			return nil, errors.New("Not rectangular grid")
		}
	}
	return &Grid{
		w:    w,
		h:    h,
		grid: grid,
	}, nil
}

func (t Tuple4) power() int {
	return t.t * t.r * t.b * t.l
}

func (g *Grid) maxPower() int {
	max := 0
	for _, v := range g.power {
		k := v.power()
		if k > max {
			max = k
		}
	}
	return max
}

func (g *Grid) computeVisibleSet() {
	g.visible = map[Tuple2]struct{}{}
	g.power = map[Tuple2]*Tuple4{}

	for y := 0; y < g.h; y++ {
		// l2r
		tallest := -1
		for x := 0; x < g.w; x++ {
			pos := Tuple2{x: x, y: y}
			k := g.grid[y][x]

			if k > tallest {
				tallest = k
				g.visible[pos] = struct{}{}
			}
		}

		// r2l
		tallest = -1
		for x := g.w - 1; x >= 0; x-- {
			pos := Tuple2{x: x, y: y}
			k := g.grid[y][x]

			if k > tallest {
				tallest = k
				g.visible[pos] = struct{}{}
			}
		}
	}

	for x := 0; x < g.w; x++ {
		// t2b
		tallest := -1
		for y := 0; y < g.h; y++ {
			pos := Tuple2{x: x, y: y}
			k := g.grid[y][x]

			if k > tallest {
				tallest = k
				g.visible[pos] = struct{}{}
			}
		}

		// b2t
		tallest = -1
		for y := g.h - 1; y >= 0; y-- {
			pos := Tuple2{x: x, y: y}
			k := g.grid[y][x]

			if k > tallest {
				tallest = k
				g.visible[pos] = struct{}{}
			}
		}
	}
}
