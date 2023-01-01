package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

const (
	puzzleInput = "input.txt"
	puzzleRow   = 2000000
	puzzleBound = 4000000
)

func main() {
	lineRegex := regexp.MustCompile(`x=(-?\d+).*y=(-?\d+).*x=(-?\d+).*y=(-?\d+)`)
	file, err := os.Open(puzzleInput)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	var sensors []Sensor
	beacons := map[Pos]struct{}{}
	first := true
	var leftBound, rightBound int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		matches := lineRegex.FindStringSubmatch(scanner.Text())
		if len(matches) == 0 {
			log.Fatalln("Invalid line")
		}
		x1, err := strconv.Atoi(matches[1])
		if err != nil {
			log.Fatalln(err)
		}
		y1, err := strconv.Atoi(matches[2])
		if err != nil {
			log.Fatalln(err)
		}
		x2, err := strconv.Atoi(matches[3])
		if err != nil {
			log.Fatalln(err)
		}
		y2, err := strconv.Atoi(matches[4])
		if err != nil {
			log.Fatalln(err)
		}
		pos := Pos{
			y: y1,
			x: x1,
		}
		beacon := Pos{
			y: y2,
			x: x2,
		}
		beacons[beacon] = struct{}{}
		radius := manhattanDistance(pos, beacon)
		sensor := Sensor{
			pos:    pos,
			radius: radius,
		}
		sensors = append(sensors, sensor)
		if x1, x2, ok := sensor.BoundsX(puzzleRow); ok {
			if first {
				first = false
				leftBound = x1
				rightBound = x2
			} else {
				if x1 < leftBound {
					leftBound = x1
				}
				if x2 > rightBound {
					rightBound = x2
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	{
		count := 0
		if !first {
			for x := leftBound; x <= rightBound; x++ {
				pos := Pos{
					y: puzzleRow,
					x: x,
				}
				if _, ok := beacons[pos]; ok {
					continue
				}
				for _, i := range sensors {
					if !i.InRange(pos) {
						continue
					}
					_, x2, ok := i.BoundsX(puzzleRow)
					if !ok {
						log.Fatalln("Invariant violated")
					}
					if x2 == x {
						count++
					} else {
						count += x2 - x
						x = x2 - 1
					}
					break
				}
			}
		}
		fmt.Println("Part 1:", count)
	}
	{
	outer2:
		for y := 0; y <= puzzleBound; y++ {
		outer:
			for x := 0; x <= puzzleBound; x++ {
				pos := Pos{
					y: y,
					x: x,
				}
				for _, i := range sensors {
					if !i.InRange(pos) {
						continue
					}
					_, x2, ok := i.BoundsX(y)
					if !ok {
						log.Fatalln("Invariant violated")
					}
					x = x2
					continue outer
				}
				fmt.Println("Part 2:", x*puzzleBound+y)
				break outer2
			}
		}
	}
}

type (
	Pos struct {
		y, x int
	}

	Sensor struct {
		pos    Pos
		radius int
	}
)

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func manhattanDistance(a, b Pos) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}

func (s *Sensor) InRange(pos Pos) bool {
	return manhattanDistance(s.pos, pos) <= s.radius
}

func (s *Sensor) BoundsX(y int) (int, int, bool) {
	vdelta := abs(s.pos.y - y)
	if vdelta > s.radius {
		return 0, 0, false
	}
	delta := s.radius - vdelta
	return s.pos.x - delta, s.pos.x + delta, true
}
