package main

import (
	"bufio"
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

	var start, end Pos
	hasStart, hasEnd := false, false

	var starts2 []Pos

	var grid [][]byte

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := []byte(scanner.Text())
		for x, i := range line {
			switch i {
			case 'S':
				start = Pos{
					y: len(grid),
					x: x,
				}
				hasStart = true
				line[x] = 'a'
				starts2 = append(starts2, Pos{
					y: len(grid),
					x: x,
				})
			case 'E':
				end = Pos{
					y: len(grid),
					x: x,
				}
				hasEnd = true
				line[x] = 'z'
			case 'a':
				starts2 = append(starts2, Pos{
					y: len(grid),
					x: x,
				})
			}
		}
		grid = append(grid, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	if !hasStart || !hasEnd {
		log.Fatalln("Missing start or end")
	}

	g := NewGrid(grid)
	{
		_, steps := AStarSearch([]Pos{start}, end, g.neighbors, manhattanDistance)
		fmt.Println("Part 1:", steps)
	}
	{
		_, steps := AStarSearch(starts2, end, g.neighbors, manhattanDistance)
		fmt.Println("Part 2:", steps)
	}
}

type (
	Pos struct {
		y, x int
	}

	Grid struct {
		h    int
		w    int
		grid [][]byte
	}
)

func NewGrid(grid [][]byte) *Grid {
	return &Grid{
		h:    len(grid),
		w:    len(grid[0]),
		grid: grid,
	}
}

func (g Grid) inBounds(k Pos) bool {
	return k.x >= 0 && k.x < g.w && k.y >= 0 && k.y < g.h
}

func (g Grid) get(k Pos) (byte, bool) {
	if !g.inBounds(k) {
		return 0, false
	}
	return g.grid[k.y][k.x], true
}

var (
	dirDeltas = []Pos{
		{-1, 0},
		{0, 1},
		{1, 0},
		{0, -1},
	}
)

func (g Grid) neighbors(k Pos) []AStarEdge[Pos] {
	limit := g.grid[k.y][k.x] + 1
	var e []AStarEdge[Pos]
	for _, i := range dirDeltas {
		k := Pos{
			y: k.y + i.y,
			x: k.x + i.x,
		}
		v, ok := g.get(k)
		if !ok || v > limit {
			continue
		}
		e = append(e, AStarEdge[Pos]{
			Key: k,
			DG:  1,
		})
	}
	return e
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func manhattanDistance(a, b Pos) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}
