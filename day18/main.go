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

	cloud := map[Point]struct{}{}
	first := true
	var start Point

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		strcoords := strings.Split(scanner.Text(), ",")
		if len(strcoords) != 3 {
			log.Fatalln("Invalid line")
		}
		x, err := strconv.Atoi(strcoords[0])
		if err != nil {
			log.Fatalln(err)
		}
		y, err := strconv.Atoi(strcoords[1])
		if err != nil {
			log.Fatalln(err)
		}
		z, err := strconv.Atoi(strcoords[2])
		if err != nil {
			log.Fatalln(err)
		}
		cloud[Point{x: x, y: y, z: z}] = struct{}{}
		if first {
			first = false
			start = Point{x: x, y: y, z: z}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	border := map[Point]struct{}{}

	surfaceArea := 0
	for k := range cloud {
		surfaceArea += 6 - numNeighbors(k, cloud)
		for _, i := range interCardinalDirs {
			p := k.add(i)
			if _, ok := cloud[p]; !ok {
				border[p] = struct{}{}
				if p.x > start.x {
					start = p
				}
			}
		}
	}
	fmt.Println("Part 1:", surfaceArea)

	extSurfaceArea := 0
	openSet := []Point{start}
	closedSet := map[Point]struct{}{start: {}}
	for len(openSet) != 0 {
		l := len(openSet) - 1
		p := openSet[l]
		openSet = openSet[:l]

		for _, i := range cardinalDirs {
			k := p.add(i)
			if _, ok := cloud[k]; ok {
				extSurfaceArea++
			} else if _, ok := border[k]; ok {
				if _, ok := closedSet[k]; !ok {
					openSet = append(openSet, k)
					closedSet[k] = struct{}{}
				}
			}
		}
	}
	fmt.Println("Part 2:", extSurfaceArea)
}

type (
	Point struct {
		x, y, z int
	}
)

func (p Point) add(other Point) Point {
	return Point{
		x: p.x + other.x,
		y: p.y + other.y,
		z: p.z + other.z,
	}
}

var (
	cardinalDirs = []Point{
		{1, 0, 0},
		{-1, 0, 0},
		{0, 1, 0},
		{0, -1, 0},
		{0, 0, 1},
		{0, 0, -1},
	}
	interCardinalDirs = []Point{
		{-1, -1, -1},
		{-1, -1, 0},
		{-1, -1, 1},
		{-1, 0, -1},
		{-1, 0, 0},
		{-1, 0, 1},
		{-1, 1, -1},
		{-1, 1, 0},
		{-1, 1, 1},
		{0, -1, -1},
		{0, -1, 0},
		{0, -1, 1},
		{0, 0, -1},
		{0, 0, 1},
		{0, 1, -1},
		{0, 1, 0},
		{0, 1, 1},
		{1, -1, -1},
		{1, -1, 0},
		{1, -1, 1},
		{1, 0, -1},
		{1, 0, 0},
		{1, 0, 1},
		{1, 1, -1},
		{1, 1, 0},
		{1, 1, 1},
	}
)

func numNeighbors(p Point, cloud map[Point]struct{}) int {
	count := 0
	for _, i := range cardinalDirs {
		if _, ok := cloud[p.add(i)]; ok {
			count++
		}
	}
	return count
}
