package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
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

	count1 := 0
	count2 := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		a, b, err := parseLine(scanner.Text())
		if err != nil {
			log.Fatalln(err)
		}
		if isFullyContained(a, b) {
			count1++
		} else if isOverlap(a, b) {
			count2++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Part 1:", count1)
	fmt.Println("Part 2:", count1+count2)
}

var (
	lineRegex = regexp.MustCompile(`^([0-9]+)-([0-9]+),([0-9]+)-([0-9]+)$`)
)

func parseLine(line string) (Pair, Pair, error) {
	match := lineRegex.FindStringSubmatch(line)
	if len(match) == 0 {
		return Pair{}, Pair{}, errors.New("Invalid line")
	}
	a, err := NewPairStr(match[1], match[2])
	if err != nil {
		return Pair{}, Pair{}, err
	}
	b, err := NewPairStr(match[3], match[4])
	if err != nil {
		return Pair{}, Pair{}, err
	}
	return a, b, nil
}

type (
	Pair struct {
		x, y int
	}
)

func NewPair(x, y int) Pair {
	if x < y {
		return Pair{
			x: x,
			y: y,
		}
	}
	return Pair{
		x: y,
		y: x,
	}
}

func NewPairStr(x, y string) (Pair, error) {
	i, err := strconv.Atoi(x)
	if err != nil {
		return Pair{}, err
	}
	j, err := strconv.Atoi(y)
	if err != nil {
		return Pair{}, err
	}
	return NewPair(i, j), nil
}

func isFullyContained(a, b Pair) bool {
	return isInclusive(a, b) || isInclusive(b, a)
}

func isInclusive(a, b Pair) bool {
	return b.x >= a.x && b.y <= a.y
}

func isOverlap(a, b Pair) bool {
	intervals := []Pair{
		{x: a.x, y: 0},
		{x: a.y + 1, y: 1},
		{x: b.x, y: 0},
		{x: b.y + 1, y: 1},
	}
	sort.Slice(intervals, func(i, j int) bool {
		a := intervals[i]
		b := intervals[j]
		if a.x != b.x {
			return a.x < b.x
		}
		if a.y != b.y {
			return a.y == 1
		}
		return false
	})

	count := 0

	for _, i := range intervals {
		if i.y == 0 {
			count++
			if count > 1 {
				return true
			}
		} else {
			count--
		}
	}

	return false
}
