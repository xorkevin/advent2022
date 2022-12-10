package main

import (
	"bufio"
	"errors"
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

	vm := NewVM()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		instr := line[0]
		arg := 0
		if len(line) > 1 {
			var err error
			arg, err = strconv.Atoi(line[1])
			if err != nil {
				log.Fatalln(err)
			}
		}
		vm.Exec(instr, arg)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Part 1:", vm.strength)
	fmt.Println("Part 2:")
	for _, i := range vm.grid {
		fmt.Println(string(i))
	}
}

type (
	VM struct {
		cycle       int
		rx          int
		strength    int
		cycleTarget int
		cycleIncr   int
		grid        [][]byte
		scanline    int
	}
)

const (
	screenWidth  = 40
	screenHeight = 6
	screenSize   = screenWidth * screenHeight
)

func NewVM() *VM {
	grid := make([][]byte, screenHeight)
	for i := range grid {
		grid[i] = make([]byte, screenWidth)
	}
	return &VM{
		cycle:       0,
		rx:          1,
		strength:    0,
		cycleTarget: 20,
		cycleIncr:   40,
		grid:        grid,
		scanline:    0,
	}
}

func (m *VM) Exec(instr string, arg int) error {
	switch instr {
	case "noop":
		m.execNoop()
	case "addx":
		m.execAddX(arg)
	default:
		return errors.New("Invalid instruction")
	}
	return nil
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func (m *VM) intersect(x int) bool {
	return abs(x-m.rx) <= 1
}

func (m *VM) checkCycle(cycle int, set bool, next int) {
	m.cycle += cycle
	if m.cycle >= m.cycleTarget {
		m.strength += m.cycleTarget * m.rx
		m.cycleTarget += m.cycleIncr
	}
	y := m.scanline / screenWidth
	x := m.scanline % screenWidth
	if m.intersect(x) {
		m.grid[y][x] = '#'
	} else {
		m.grid[y][x] = '.'
	}
	m.scanline = (m.scanline + 1) % screenSize
	if set {
		m.rx = next
	}
}

func (m *VM) execNoop() {
	m.checkCycle(1, false, 0)
}

func (m *VM) execAddX(arg int) {
	m.checkCycle(1, false, 0)
	m.checkCycle(1, true, m.rx+arg)
}
