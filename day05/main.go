package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
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

	var grid1 *Grid
	var grid2 *Grid

	var rows [][]byte

	modeGrid := true

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if modeGrid {
			line := scanner.Bytes()
			if len(line) == 0 {
				g, err := newGridFromRows(rows)
				if err != nil {
					log.Fatalln(err)
				}
				grid1 = g
				grid2 = grid1.Clone()
				modeGrid = false
				continue
			}
			row, err := parseGridRow(line)
			if err != nil {
				log.Fatalln(err)
			}
			if len(row) != 0 {
				rows = append(rows, row)
			}
			continue
		}
		instr, err := parseInstrLine(scanner.Text())
		if err != nil {
			log.Fatalln(err)
		}
		grid1.processInstr1(instr)
		grid2.processInstr2(instr)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Part 1:", string(grid1.tops()))
	fmt.Println("Part 2:", string(grid2.tops()))
}

type (
	Grid struct {
		grid [][]byte
	}
)

func newGridFromRows(rows [][]byte) (*Grid, error) {
	h := len(rows)
	if h == 0 {
		return nil, errors.New("No rows")
	}
	w := len(rows[0])
	for _, i := range rows {
		if len(i) != w {
			return nil, errors.New("Mismatched rows")
		}
	}
	grid := make([][]byte, 0, w)
	for i := 0; i < w; i++ {
		col := make([]byte, 0, h)
		for j := 0; j < h; j++ {
			c := rows[h-j-1][i]
			if c == '.' {
				break
			}
			col = append(col, c)
		}
		grid = append(grid, col)
	}
	return &Grid{
		grid: grid,
	}, nil
}

func (g Grid) Clone() *Grid {
	grid := make([][]byte, 0, len(g.grid))
	for _, i := range g.grid {
		row := make([]byte, len(i))
		copy(row, i)
		grid = append(grid, row)
	}
	return &Grid{
		grid: grid,
	}
}

func (g *Grid) processInstr1(instr Instr) {
	for i := 0; i < instr.a; i++ {
		g.push(instr.c, g.pop(instr.b))
	}
}

func (g *Grid) processInstr2(instr Instr) {
	stack := make([]byte, 0, instr.a)
	for i := 0; i < instr.a; i++ {
		stack = append(stack, g.pop(instr.b))
	}
	for i, l := 0, instr.a; i < l; i++ {
		g.push(instr.c, stack[l-i-1])
	}
}

func (g *Grid) pop(col int) byte {
	l := len(g.grid[col]) - 1
	b := g.grid[col][l]
	g.grid[col] = g.grid[col][:l]
	return b
}

func (g *Grid) push(col int, b byte) {
	g.grid[col] = append(g.grid[col], b)
}

func (g *Grid) top(col int) byte {
	return g.grid[col][len(g.grid[col])-1]
}

func (g *Grid) tops() []byte {
	t := make([]byte, 0, len(g.grid))
	for i := range g.grid {
		t = append(t, g.top(i))
	}
	return t
}

func (g Grid) String() string {
	var b strings.Builder
	for n, i := range g.grid {
		b.WriteString(strconv.Itoa(n + 1))
		b.WriteString(": ")
		b.WriteString(string(i))
		b.WriteString("\n")
	}
	return b.String()
}

type (
	Instr struct {
		a, b, c int
	}
)

var (
	lineRegex = regexp.MustCompile(`^move (\d+) from (\d+) to (\d+)$`)
)

func parseInstrLine(line string) (Instr, error) {
	captures := lineRegex.FindStringSubmatch(line)
	if len(captures) == 0 {
		return Instr{}, errors.New("Invalid line")
	}
	a, err := strconv.Atoi(captures[1])
	if err != nil {
		return Instr{}, err
	}
	b, err := strconv.Atoi(captures[2])
	if err != nil {
		return Instr{}, err
	}
	c, err := strconv.Atoi(captures[3])
	if err != nil {
		return Instr{}, err
	}
	return Instr{
		a: a,
		b: b - 1,
		c: c - 1,
	}, nil
}

func parseGridRow(line []byte) ([]byte, error) {
	row := make([]byte, 0, (len(line)/4)+1)
	for i, l := 0, len(line)-1; i < l; {
		var c byte
		switch line[i] {
		case ' ':
			if line[i+1] != ' ' {
				return nil, nil
			}
			c = '.'
		case '[':
			c = line[i+1]
		default:
			return nil, errors.New("Invalid line")
		}
		row = append(row, c)
		i += 4
	}
	return row, nil
}
