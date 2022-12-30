package main

import (
	"bufio"
	"errors"
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

	dfa := NewDFA[TKind, byte](tKindZero)
	{
		dfaNum := NewDFA[TKind, byte](tKindNum)
		dfa.AddDFA([]byte("0123456789"), dfaNum)
		dfaNum.AddDFA([]byte("0123456789"), dfaNum)

		dfa.AddPath([]byte("["), tKindOpen, tKindZero)
		dfa.AddPath([]byte("]"), tKindClose, tKindZero)
		dfa.AddPath([]byte(","), tKindComma, tKindZero)
	}
	lexer := NewLexer(dfa, tKindZero, tKindEOS, map[TKind]struct{}{
		tKindComma: {},
	})

	div1 := &Signal{isList: true, list: []*Signal{{num: 2}}}
	div2 := &Signal{isList: true, list: []*Signal{{num: 6}}}

	signals := []*Signal{div1, div2}

	var left *Signal

	count := 0
	numPairs := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			left = nil
			continue
		}
		tokens, err := lexer.Tokenize(line)
		if err != nil {
			log.Fatalln(err)
		}
		sig, rest, err := parseSignal(tokens)
		if err != nil {
			log.Fatalln(err)
		}
		if len(rest) != 1 || rest[0].Kind != tKindEOS {
			log.Fatalln("Did not parse all tokens")
		}
		if left == nil {
			left = sig
		} else {
			numPairs++
			if v, ok := compareSigs(left, sig); ok && v {
				count += numPairs
			}
		}
		signals = append(signals, sig)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Part 1:", count)

	sort.Slice(signals, func(i, j int) bool {
		v, ok := compareSigs(signals[i], signals[j])
		return v && ok
	})
	div1Idx := sort.Search(len(signals), func(i int) bool {
		if eqSigs(div1, signals[i]) {
			return true
		}
		v, ok := compareSigs(div1, signals[i])
		return v && ok
	})
	if div1Idx >= len(signals) || !eqSigs(div1, signals[div1Idx]) {
		log.Fatalln("Divider missing")
	}
	div2Idx := sort.Search(len(signals), func(i int) bool {
		if eqSigs(div2, signals[i]) {
			return true
		}
		v, ok := compareSigs(div2, signals[i])
		return v && ok
	})
	if div2Idx >= len(signals) || !eqSigs(div2, signals[div2Idx]) {
		log.Fatalln("Divider missing")
	}
	fmt.Println("Part 2:", (div1Idx+1)*(div2Idx+1))
}

func eqSigs(left, right *Signal) bool {
	if left.isList != right.isList {
		return false
	}
	if !left.isList {
		return left.num == right.num
	}
	if len(left.list) != len(right.list) {
		return false
	}
	for n, i := range left.list {
		if !eqSigs(i, right.list[n]) {
			return false
		}
	}
	return true
}

func compareSigs(left, right *Signal) (bool, bool) {
	if !left.isList && !right.isList {
		if left.num == right.num {
			return false, false
		}
		return left.num < right.num, true
	}
	if left.isList && right.isList {
		for i := 0; i < len(left.list) && i < len(right.list); i++ {
			v, ok := compareSigs(left.list[i], right.list[i])
			if ok {
				return v, true
			}
		}
		if len(left.list) == len(right.list) {
			return false, false
		}
		return len(left.list) < len(right.list), true
	}
	if !left.isList {
		left = &Signal{
			isList: true,
			list:   []*Signal{left},
		}
	} else {
		right = &Signal{
			isList: true,
			list:   []*Signal{right},
		}
	}
	return compareSigs(left, right)
}

var (
	ErrorParse = errors.New("Failed to parse")
)

type (
	Signal struct {
		num    int
		isList bool
		list   []*Signal
	}
)

type (
	TKind int
)

const (
	tKindZero TKind = iota
	tKindEOS
	tKindNum
	tKindOpen
	tKindClose
	tKindComma
)

func parseSignal(tokens []Token[TKind, byte]) (*Signal, []Token[TKind, byte], error) {
	if len(tokens) == 0 {
		return nil, tokens, fmt.Errorf("No tokens: %w", ErrorParse)
	}
	top := tokens[0]
	switch top.Kind {
	case tKindNum:
		{
			num, err := strconv.Atoi(string(top.Val))
			if err != nil {
				return nil, tokens, fmt.Errorf("%w: %v", ErrorParse, err)
			}
			return &Signal{
				num: num,
			}, tokens[1:], nil
		}
	case tKindOpen:
		{
			tokens = tokens[1:]
			var signals []*Signal
			for {
				if len(tokens) == 0 {
					return nil, tokens, fmt.Errorf("No tokens: %w", ErrorParse)
				}
				if tokens[0].Kind == tKindClose {
					tokens = tokens[1:]
					break
				}
				sig, rest, err := parseSignal(tokens)
				if err != nil {
					return nil, rest, err
				}
				signals = append(signals, sig)
				tokens = rest
			}
			return &Signal{
				isList: true,
				list:   signals,
			}, tokens, nil
		}
	default:
		return nil, tokens, errors.New("Failed to parse")
	}
}
