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

	term := NewTerm()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if err := term.ReadInput(scanner.Text()); err != nil {
			log.Fatalln(err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	totalSize, smallDirs := calcSmallDirSize(term.root)

	fmt.Println("Part 1:", smallDirs)

	freed := 0
	target := minUnused + totalSize - totalDisk
	if target > 0 {
		freed = findDirSizeTarget(term.root, target)
	}
	fmt.Println("Part 2:", freed)
}

const (
	totalDisk = 70000000
	minUnused = 30000000
)

func findDirSizeTarget(n *Node, target int) int {
	if n.size < target {
		return -1
	}
	atMost := n.size
	for _, v := range n.children {
		if !v.isDir {
			continue
		}
		k := findDirSizeTarget(v, target)
		if k >= 0 && k < atMost {
			atMost = k
		}
	}
	return atMost
}

const (
	smallDirLimit = 100000
)

func calcSmallDirSize(n *Node) (int, int) {
	total := 0
	cummulative := 0
	for _, v := range n.children {
		if v.isDir {
			t, c := calcSmallDirSize(v)
			total += t
			cummulative += c
		} else {
			total += v.size
		}
	}
	n.size = total
	if total <= smallDirLimit {
		cummulative += total
	}
	return total, cummulative
}

type (
	Term struct {
		pwd     []string
		root    *Node
		running string
	}

	Node struct {
		name     string
		children map[string]*Node
		isDir    bool
		size     int
	}
)

func NewTerm() *Term {
	return &Term{
		pwd:     nil,
		root:    NewDir("/"),
		running: "",
	}
}

func NewDir(name string) *Node {
	return &Node{
		name:     name,
		children: map[string]*Node{},
		isDir:    true,
		size:     0,
	}
}

func NewFile(name string, size int) *Node {
	return &Node{
		name:     name,
		children: map[string]*Node{},
		isDir:    false,
		size:     size,
	}
}

var (
	cmdRegex = regexp.MustCompile(`^\$ `)
)

func (t *Term) ReadInput(inp string) error {
	if cmdRegex.MatchString(inp) {
		t.running = ""
		return t.Exec(strings.Fields(inp[2:]))
	}
	if t.running != "" {
		switch t.running {
		case "ls":
			return t.ReadOutputLS(inp)
		}
	}
	return errors.New("Invalid input")
}

func (t *Term) Exec(cmd []string) error {
	if len(cmd) == 0 {
		return nil
	}
	switch cmd[0] {
	case "cd":
		{
			if len(cmd) != 2 {
				return errors.New("Invalid cd args")
			}
			return t.CD(cmd[1])
		}
	case "ls":
		{
			if len(cmd) != 1 {
				return errors.New("Invalid ls args")
			}
			t.running = "ls"
			return nil
		}
	default:
		return errors.New("Invalid cmd")
	}
}

func (t *Term) CD(d string) error {
	switch d {
	case "":
		{
			return errors.New("Invalid cd directory")
		}
	case "..":
		{
			if len(t.pwd) == 0 {
				return errors.New("No parent directory from root")
			}
			t.pwd = t.pwd[:len(t.pwd)-1]
			return nil
		}
	case "/":
		{
			t.pwd = t.pwd[:0]
			return nil
		}
	default:
		{
			t.pwd = append(t.pwd, d)
			return nil
		}
	}
}

func (t *Term) ReadOutputLS(inp string) error {
	kind, name, ok := strings.Cut(inp, " ")
	if !ok {
		return errors.New("Invalid ls output")
	}
	if name == "" {
		return errors.New("Invalid ls file name")
	}
	if kind == "dir" {
		return t.Mkdir(name)
	}
	size, err := strconv.Atoi(kind)
	if err != nil {
		return err
	}
	return t.Touch(name, size)
}

func (t *Term) MkdirPath() (*Node, error) {
	node := t.root
	for _, i := range t.pwd {
		if n, ok := node.children[i]; !ok {
			node.children[i] = NewDir(i)
		} else if !n.isDir {
			return nil, errors.New("Mkdir invalid path")
		}
		node = node.children[i]
	}
	return node, nil
}

func (t *Term) Mkdir(name string) error {
	node, err := t.MkdirPath()
	if err != nil {
		return err
	}
	if n, ok := node.children[name]; !ok {
		node.children[name] = NewDir(name)
	} else if !n.isDir {
		return errors.New("Mkdir on non-dir")
	}
	return nil
}

func (t *Term) Touch(name string, size int) error {
	node, err := t.MkdirPath()
	if err != nil {
		return err
	}
	if n, ok := node.children[name]; ok && n.isDir {
		return errors.New("Touch file on dir")
	}
	node.children[name] = NewFile(name, size)
	return nil
}
