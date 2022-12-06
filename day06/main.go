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
		c := file[i]

		if startpos1 < 0 {
			if _, ok := seen1[c]; !ok {
				seen1[c] = 0
			}
			seen1[c]++
			if i > 3 {
				if seen1[file[i-4]] < 2 {
					delete(seen1, file[i-4])
				} else {
					seen1[file[i-4]]--
				}
			}
			if i >= 3 && isUniq(seen1) {
				startpos1 = i + 1
			}
		}

		if startpos2 < 0 {
			if _, ok := seen2[c]; !ok {
				seen2[c] = 0
			}
			seen2[c]++
			if i > 13 {
				if seen2[file[i-14]] < 2 {
					delete(seen2, file[i-14])
				} else {
					seen2[file[i-14]]--
				}
			}
			if i >= 13 && isUniq(seen2) {
				startpos2 = i + 1
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
