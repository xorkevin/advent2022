package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
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

	sum1 := 0
	sum2 := 0

	var groupCommon []byte
	groupSize := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := []byte(scanner.Text())
		if len(line)%2 != 0 {
			log.Fatalln("Invalid line format")
		}
		halflen := len(line) / 2
		c, ok := findCommon(line[:halflen], line[halflen:])
		if !ok {
			log.Fatalln("None in common")
		}
		p, err := prio(c)
		if err != nil {
			log.Fatalln(err)
		}
		sum1 += p

		if groupSize < 2 {
			if groupSize == 0 {
				groupCommon = line
			} else {
				groupCommon = findCommonAll(groupCommon, line)
			}
			groupSize++
			continue
		}

		c, ok = findCommon(groupCommon, line)
		if !ok {
			log.Fatalln("None in common")
		}
		p, err = prio(c)
		if err != nil {
			log.Fatalln(err)
		}
		sum2 += p
		groupCommon = nil
		groupSize = 0
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Part 1:", sum1)
	fmt.Println("Part 2:", sum2)
}

func findCommonAll(a, b []byte) []byte {
	first := make(map[byte]struct{}, len(a))
	for _, i := range a {
		first[i] = struct{}{}
	}
	common := make([]byte, 0, len(a))
	for _, i := range b {
		if _, ok := first[i]; ok {
			common = append(common, i)
		}
	}
	return common
}

func findCommon(a, b []byte) (byte, bool) {
	first := make(map[byte]struct{}, len(a))
	for _, i := range a {
		first[i] = struct{}{}
	}
	for _, i := range b {
		if _, ok := first[i]; ok {
			return i, true
		}
	}
	return 0, false
}

func prio(c byte) (int, error) {
	if c >= 'a' && c <= 'z' {
		return int(c) - 'a' + 1, nil
	}
	if c >= 'A' && c <= 'Z' {
		return int(c) - 'A' + 27, nil
	}
	return 0, errors.New("Invalid prio")
}
