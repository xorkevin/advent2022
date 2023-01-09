package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

const (
	puzzleInput = "input.txt"
)

func main() {
	lineRegex := regexp.MustCompile(`Valve (\w+) has flow rate=(\d+); tunnels? leads? to valves? (.*)`)

	file, err := os.Open(puzzleInput)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	var valves []Valve
	var nodes []string
	var edges []Edge[string]

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		matches := lineRegex.FindStringSubmatch(scanner.Text())
		if len(matches) == 0 {
			log.Fatalln("Invalid line")
		}
		name := matches[1]
		rate, err := strconv.Atoi(matches[2])
		if err != nil {
			log.Fatalln(err)
		}
		nodes = append(nodes, name)
		if rate > 0 {
			valves = append(valves, Valve{
				Name: name,
				Rate: rate,
			})
		}
		for _, i := range strings.Split(matches[3], ", ") {
			edges = append(edges, Edge[string]{
				A: name,
				B: i,
				C: 1,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	sort.Slice(valves, func(i, j int) bool {
		return valves[i].Rate > valves[j].Rate
	})
	for n, i := range valves {
		k := i
		k.ID = 1 << n
		valves[n] = k
	}

	dist := FloydWarshall(nodes, edges)

	fmt.Println("Part 1:", searchMax(0, 30, "AA", map[string]struct{}{}, valves, dist, 0))
	allPaths := map[int]int{}
	searchPaths(0, 0, 26, "AA", map[string]struct{}{}, allPaths, valves, dist)
	maxFlow := 0
	for k, v := range allPaths {
		for k2, v2 := range allPaths {
			if k&k2 != 0 {
				continue
			}
			flow := v + v2
			if flow > maxFlow {
				maxFlow = flow
			}
		}
	}
	fmt.Println("Part 2:", maxFlow)
}

type (
	Valve struct {
		Name string
		ID   int
		Rate int
	}
)

func searchMax(acc int, remaining int, pos string, toggled map[string]struct{}, valves []Valve, dist PairwiseDistances[string], candidate int) int {
	if remaining <= 0 {
		return acc
	}

	{
		bound := 0
		t := remaining
		for _, i := range valves {
			if t <= 2 {
				break
			}
			if _, ok := toggled[i.Name]; ok {
				continue
			}
			t -= 2
			bound += t * i.Rate
		}
		if acc+bound < candidate {
			return 0
		}
	}

	maxFlow := acc
	if maxFlow > candidate {
		candidate = maxFlow
	}

	for _, i := range valves {
		if _, ok := toggled[i.Name]; ok {
			continue
		}
		cost, ok := dist.EdgeCost(pos, i.Name)
		if !ok {
			continue
		}
		nextRemaining := remaining - cost - 1
		if nextRemaining <= 0 {
			continue
		}
		toggled[i.Name] = struct{}{}
		flow := searchMax(acc+i.Rate*nextRemaining, nextRemaining, i.Name, toggled, valves, dist, candidate)
		delete(toggled, i.Name)
		if flow > maxFlow {
			maxFlow = flow
			if maxFlow > candidate {
				candidate = maxFlow
			}
		}
	}

	return maxFlow
}

func searchPaths(acc int, curPath int, remaining int, pos string, toggled map[string]struct{}, allPaths map[int]int, valves []Valve, dist PairwiseDistances[string]) {
	if remaining <= 0 {
		return
	}

	if v, ok := allPaths[curPath]; !ok || acc > v {
		allPaths[curPath] = acc
	}

	for _, i := range valves {
		if _, ok := toggled[i.Name]; ok {
			continue
		}
		cost, ok := dist.EdgeCost(pos, i.Name)
		if !ok {
			continue
		}
		nextRemaining := remaining - cost - 1
		if nextRemaining <= 0 {
			continue
		}
		toggled[i.Name] = struct{}{}
		searchPaths(acc+i.Rate*nextRemaining, curPath|i.ID, nextRemaining, i.Name, toggled, allPaths, valves, dist)
		delete(toggled, i.Name)
	}
}
