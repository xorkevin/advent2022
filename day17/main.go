package main

import (
	"fmt"
	"log"
	"os"
	"unicode"
)

const (
	puzzleInput = "input.txt"
	puzzlePart1 = 2022
	puzzlePart2 = 1_000_000_000_000
)

func main() {
	puzzleBytes, err := os.ReadFile(puzzleInput)
	if err != nil {
		log.Fatalln(err)
	}
	if unicode.IsSpace(rune(puzzleBytes[len(puzzleBytes)-1])) {
		puzzleBytes = puzzleBytes[:len(puzzleBytes)-1]
	}

	sim := NewSim()
	shapeCount := 0
	nextShape := 0
	nextByte := 0

	seenStates := map[int]SimState{}

	foundCycle := false
	skipped := false

	var dShape int
	var dHeight int

	skippedShapes := 0
	skippedHeight := 0

	for {
		if !sim.hasShape() {
			if shapeCount == puzzlePart1 {
				fmt.Println("Part 1:", sim.top)
			}

			if !foundCycle {
				state := nextByte*len(shapes) + nextShape
				if v, ok := seenStates[state]; ok {
					if v.count == 2 {
						dShape = shapeCount - v.shapeCount
						dHeight = sim.top - v.height
						foundCycle = true
					}
					seenStates[state] = SimState{
						count:      v.count + 1,
						shapeCount: shapeCount,
						height:     sim.top,
					}
				} else {
					seenStates[state] = SimState{
						count:      1,
						shapeCount: shapeCount,
						height:     sim.top,
					}
				}
			}

			if !skipped && foundCycle && shapeCount >= puzzlePart1 {
				cycles := (puzzlePart2 - shapeCount) / dShape
				skippedShapes = cycles * dShape
				skippedHeight = cycles * dHeight
				skipped = true
			}

			if shapeCount+skippedShapes == puzzlePart2 {
				fmt.Println("Part 2:", sim.top+skippedHeight)
				break
			}

			sim.addShape(nextShape)
			shapeCount++
			nextShape = (nextShape + 1) % len(shapes)
		}
		b := puzzleBytes[nextByte]
		nextByte = (nextByte + 1) % len(puzzleBytes)
		var pushRight bool
		switch b {
		case '<':
			pushRight = false
		case '>':
			pushRight = true
		default:
			log.Fatalln("Invalid dir")
		}
		sim.pushDir(pushRight)
		if !sim.fall() {
			sim.commitShape()
		}
	}
}

type (
	SimState struct {
		count      int
		shapeCount int
		height     int
	}

	Pos struct {
		y, x int
	}

	Sim struct {
		pos  Pos
		kind int
		grid [][]byte
		top  int
	}
)

func NewSim() *Sim {
	s := &Sim{
		pos:  Pos{},
		kind: -1,
		grid: [][]byte{},
		top:  0,
	}
	s.addRows()
	return s
}

func (s *Sim) addRows() {
	for len(s.grid) < s.top+7 {
		s.grid = append(s.grid, []byte("......."))
	}
}

func (s *Sim) hasShape() bool {
	return s.kind >= 0
}

func (s *Sim) addShape(kind int) {
	s.pos = Pos{
		y: s.top + 3,
		x: 2,
	}
	s.kind = kind
}

func (s *Sim) commitShape() {
	shape := shapes[s.kind]
	h := len(shape)
	for yp, i := range shape {
		for x, j := range i {
			if j == '#' {
				y := h - yp - 1
				ny := s.pos.y + y
				s.grid[ny][s.pos.x+x] = '#'
				if t := ny + 1; t > s.top {
					s.top = t
				}
			}
		}
	}
	s.kind = -1
	s.addRows()
}

func (s *Sim) pushDir(right bool) {
	dir := -1
	if right {
		dir = 1
	}
	if s.checkShapeCollision(0, dir) {
		return
	}
	s.pos = Pos{
		y: s.pos.y,
		x: s.pos.x + dir,
	}
}

func (s *Sim) fall() bool {
	if s.checkShapeCollision(-1, 0) {
		return false
	}
	s.pos = Pos{
		y: s.pos.y - 1,
		x: s.pos.x,
	}
	return true
}

func (s *Sim) checkShapeCollision(dy, dx int) bool {
	shape := shapes[s.kind]
	h := len(shape)
	for yp, i := range shape {
		for x, j := range i {
			if j == '#' {
				y := h - yp - 1
				if s.isGridBlock(s.pos.y+y+dy, s.pos.x+x+dx) {
					return true
				}
			}
		}
	}
	return false
}

// isShapeBlock computes if coordinate is occupied by the shape in shape
// relative coordinates. Shape data is stored in the -y direction.
func (s *Sim) isShapeBlock(kind int, y, x int) bool {
	shape := shapes[kind]
	if y < 0 || y >= len(shape) {
		return false
	}
	row := shape[len(shape)-y-1]
	if x < 0 || x >= len(row) {
		return false
	}
	return row[x] == '#'
}

// isGridBlock computes if coordinate is occupied by a barrier in grid
// coordinates. Upwards is the +y direction.
func (s *Sim) isGridBlock(y, x int) bool {
	if y < 0 || y >= len(s.grid) {
		return true
	}
	row := s.grid[y]
	if x < 0 || x >= len(row) {
		return true
	}
	return row[x] == '#'
}

func (s *Sim) render() {
	for i := len(s.grid) - 1; i >= 0; i-- {
		if s.kind >= 0 && i >= s.pos.y && i < s.pos.y+4 {
			for j := 0; j < 7; j++ {
				if s.isShapeBlock(s.kind, i-s.pos.y, j-s.pos.x) {
					fmt.Print("@")
				} else if s.isGridBlock(i, j) {
					fmt.Print("#")
				} else {
					fmt.Print(".")
				}
			}
			fmt.Println()
		} else {
			fmt.Println(string(s.grid[i]))
		}
	}
}

var (
	shapes = [][][]byte{
		{
			[]byte("####"),
		},
		{
			[]byte(".#."),
			[]byte("###"),
			[]byte(".#."),
		},
		{
			[]byte("..#"),
			[]byte("..#"),
			[]byte("###"),
		},
		{
			[]byte("#"),
			[]byte("#"),
			[]byte("#"),
			[]byte("#"),
		},
		{
			[]byte("##"),
			[]byte("##"),
		},
	}
)
