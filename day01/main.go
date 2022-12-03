package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

	var weights []int
	current := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			weights = append(weights, current)
			current = 0
			continue
		}

		num, err := strconv.Atoi(line)
		if err != nil {
			log.Fatalln(err)
		}
		current += num
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	weights = append(weights, current)
	sort.Ints(weights)
	if len(weights) == 0 {
		return
	}
	last := len(weights) - 1
	fmt.Println("Part 1:", weights[last])
	if last < 2 {
		return
	}
	fmt.Println("Part 2:", weights[last]+weights[last-1]+weights[last-2])
}
