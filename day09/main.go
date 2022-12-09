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

	rope1 := NewRope(1)
	rope2 := NewRope(9)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		dir, countstr, ok := strings.Cut(scanner.Text(), " ")
		if !ok {
			log.Fatalln("Invalid line")
		}
		count, err := strconv.Atoi(countstr)
		if err != nil {
			log.Fatalln(err)
		}
		for i := 0; i < count; i++ {
			if err := rope1.Move(dir); err != nil {
				log.Fatalln(err)
			}
			if err := rope2.Move(dir); err != nil {
				log.Fatalln(err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Part 1:", len(rope1.history))
	fmt.Println("Part 2:", len(rope2.history))
}

type (
	Tuple2 struct {
		y, x int
	}

	Rope struct {
		h       Tuple2
		t       []Tuple2
		history map[Tuple2]struct{}
	}
)

func (t Tuple2) Move(dir string) (Tuple2, error) {
	switch dir {
	case "U":
		t.y -= 1
	case "R":
		t.x += 1
	case "D":
		t.y += 1
	case "L":
		t.x -= 1
	default:
		return Tuple2{}, errors.New("Invalid direction")
	}
	return t, nil
}

func (t Tuple2) Delta(p Tuple2) Tuple2 {
	t.x += p.x
	t.y += p.y
	return t
}

func (t Tuple2) Dist(p Tuple2) Tuple2 {
	return Tuple2{
		x: p.x - t.x,
		y: p.y - t.y,
	}
}

func (t Tuple2) Dir() Tuple2 {
	if t.x != 0 {
		if t.x > 0 {
			t.x = 1
		} else {
			t.x = -1
		}
	}
	if t.y != 0 {
		if t.y > 0 {
			t.y = 1
		} else {
			t.y = -1
		}
	}
	return t
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func (t Tuple2) MaxMag() int {
	a := abs(t.x)
	b := abs(t.y)
	if a > b {
		return a
	}
	return b
}

func NewRope(size int) *Rope {
	history := map[Tuple2]struct{}{
		{}: {},
	}
	tail := make([]Tuple2, size)
	return &Rope{
		h:       Tuple2{},
		t:       tail,
		history: history,
	}
}

func (r *Rope) Move(dir string) error {
	next, err := r.h.Move(dir)
	if err != nil {
		return err
	}
	r.h = next
	last := len(r.t) - 1
	for i, t := range r.t {
		if k := t.Dist(next); k.MaxMag() > 1 {
			next = t.Delta(k.Dir())
			if i == last {
				r.history[next] = struct{}{}
			}
		} else {
			next = t
		}
		r.t[i] = next
	}
	return nil
}
