package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

	score1 := 0
	score2 := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		games := strings.Fields(scanner.Text())
		if len(games) != 2 {
			log.Fatalln("Invalid line")
		}

		theirMove := toMove(games[0])
		ownMove := toMove(games[1])

		score1 += winner(theirMove, ownMove) + ownMove + 1
		score2 += ownMove*3 + pickMove(theirMove, ownMove) + 1
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Part 1:", score1)
	fmt.Println("Part 2:", score2)
}

func pickMove(a, b int) int {
	return (a + b + 2) % 3
}

func winner(a, b int) int {
	return ((b - a + 3 + 1) % 3) * 3
}

func toMove(k string) int {
	switch k {
	case "A", "X":
		return 0
	case "B", "Y":
		return 1
	case "C", "Z":
		return 2
	default:
		log.Fatalln("Invalid move")
		return 0
	}
}
