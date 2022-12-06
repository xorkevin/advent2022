package main

import (
	"fmt"
	"log"
	"os"
)

const (
	puzzleInput = "input.txt"
)

func main() {
	file, err := os.ReadFile(puzzleInput)
	if err != nil {
		log.Fatalln(err)
	}
	if len(file) <= 4 {
		log.Fatalln("Invalid stream")
	}

	startpos1 := -1
	startpos2 := -1

	seen1 := map[byte]int{}
	seen2 := map[byte]int{}

	for i, l := 0, len(file); (startpos1 < 0 || startpos2 < 0) && i < l; i++ {
		if startpos1 < 0 {
			if _, ok := seen1[file[i]]; !ok {
				seen1[file[i]] = 0
			}
			seen1[file[i]]++
			if i > 3 {
				seen1[file[i-4]]--
				if seen1[file[i-4]] == 0 {
					delete(seen1, file[i-4])
				}
			}
			if i >= 3 {
				if isUniq(seen1) {
					startpos1 = i + 1
				}
			}
		}

		if startpos2 < 0 {
			if _, ok := seen2[file[i]]; !ok {
				seen2[file[i]] = 0
			}
			seen2[file[i]]++
			if i > 13 {
				seen2[file[i-14]]--
				if seen2[file[i-14]] == 0 {
					delete(seen2, file[i-14])
				}
			}
			if i >= 13 {
				if isUniq(seen2) {
					startpos2 = i + 1
				}
			}
		}
	}

	fmt.Println("Part 1:", startpos1)
	fmt.Println("Part 2:", startpos2)
}

func isUniq(seen map[byte]int) bool {
	for _, v := range seen {
		if v > 1 {
			return false
		}
	}
	return true
}
