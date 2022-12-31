package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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

	grid1 := map[Pos]struct{}{}
	grid2 := map[Pos]struct{}{}
	lowest := map[int]int{}
	floor := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var points []Pos
		for _, i := range strings.Split(scanner.Text(), " -> ") {
			sx, sy, ok := strings.Cut(i, ",")
			if !ok {
				log.Fatalln("Invalid point")
			}
			x, err := strconv.Atoi(sx)
			if err != nil {
				log.Fatalln(err)
			}
			y, err := strconv.Atoi(sy)
			if err != nil {
				log.Fatalln(err)
			}
			points = append(points, Pos{y: y, x: x})
		}
		last := points[0]
		grid1[last] = struct{}{}
		grid2[last] = struct{}{}
		if v, ok := lowest[last.x]; !ok || last.y > v {
			lowest[last.x] = last.y
			if last.y > floor {
				floor = last.y
			}
		}
		for _, i := range points[1:] {
			delta := unitDelta(last, i)
			for last != i {
				last = addPos(last, delta)
				grid1[last] = struct{}{}
				grid2[last] = struct{}{}
				if v, ok := lowest[last.x]; !ok || last.y > v {
					lowest[last.x] = last.y
					if last.y > floor {
						floor = last.y
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	floor += 2

	{
		count := 0
		for dropParticle1(Pos{y: 0, x: 500}, grid1, lowest) {
			count++
		}
		fmt.Println("Part 1:", count)
	}
	{
		count := 0
		for dropParticle2(Pos{y: 0, x: 500}, grid2, floor) {
			count++
		}
		fmt.Println("Part 2:", count)
	}
}

var (
	dirs = []Pos{
		{y: 1, x: 0},
		{y: 1, x: -1},
		{y: 1, x: 1},
	}
)

func dropParticle1(a Pos, grid map[Pos]struct{}, lowest map[int]int) bool {
	if _, ok := grid[a]; ok {
		return false
	}
outer:
	for {
		if v, ok := lowest[a.x]; !ok || a.y > v {
			return false
		}
		for _, i := range dirs {
			next := addPos(a, i)
			if _, ok := grid[next]; !ok {
				a = next
				continue outer
			}
		}
		grid[a] = struct{}{}
		return true
	}
}

func dropParticle2(a Pos, grid map[Pos]struct{}, floor int) bool {
	if _, ok := grid[a]; ok {
		return false
	}
outer:
	for {
		if a.y+1 < floor {
			for _, i := range dirs {
				next := addPos(a, i)
				if _, ok := grid[next]; !ok {
					a = next
					continue outer
				}
			}
		}
		grid[a] = struct{}{}
		return true
	}
}

type (
	Pos struct {
		y, x int
	}
)

func unitDir(a, b int) int {
	if a < b {
		return 1
	} else if a > b {
		return -1
	}
	return 0
}

func unitDelta(a, b Pos) Pos {
	return Pos{
		y: unitDir(a.y, b.y),
		x: unitDir(a.x, b.x),
	}
}

func addPos(a, b Pos) Pos {
	return Pos{
		y: a.y + b.y,
		x: a.x + b.x,
	}
}
