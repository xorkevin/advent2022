package main

import (
	"errors"
)

type (
	DFA[T, C comparable] struct {
		kind  T
		nodes map[C]*DFA[T, C]
	}
)

func NewDFA[T, C comparable](kind T) *DFA[T, C] {
	return &DFA[T, C]{
		kind:  kind,
		nodes: map[C]*DFA[T, C]{},
	}
}

func (d *DFA[T, C]) AddDFA(s []C, dfa *DFA[T, C]) {
	for _, i := range s {
		d.nodes[i] = dfa
	}
}

func (d *DFA[T, C]) AddPath(path []C, kind, zero T) *DFA[T, C] {
	k := d
	for _, i := range path {
		var n *DFA[T, C]
		if v, ok := k.nodes[i]; ok {
			n = v
		} else {
			n = NewDFA[T, C](zero)
			k.nodes[i] = n
		}
		k = n
	}
	k.kind = kind
	return k
}

func (d *DFA[T, C]) Match(c C) (*DFA[T, C], bool) {
	n, ok := d.nodes[c]
	if !ok {
		return nil, false
	}
	return n, true
}

func (d *DFA[T, C]) Kind() T {
	return d.kind
}

type (
	Token[T, C comparable] struct {
		Kind T
		Val  []C
	}

	Lexer[T, C comparable] struct {
		dfa       *DFA[T, C]
		zero, eos T
		ignored   map[T]struct{}
	}
)

func NewLexer[T, C comparable](dfa *DFA[T, C], zero, eos T, ignored map[T]struct{}) *Lexer[T, C] {
	return &Lexer[T, C]{
		dfa:     dfa,
		zero:    zero,
		eos:     eos,
		ignored: ignored,
	}
}

var (
	ErrorLex = errors.New("Invalid token")
)

func (l *Lexer[T, C]) Next(chars []C) (*Token[T, C], []C, error) {
	n := l.dfa
	end := 0
	for i, c := range chars {
		next, ok := n.Match(c)
		if !ok {
			break
		}
		n = next
		end = i + 1
	}
	if n.Kind() == l.zero {
		if end == 0 && len(chars) == 0 {
			return &Token[T, C]{Kind: l.eos, Val: nil}, chars, nil
		}
		return nil, chars, ErrorLex
	}
	return &Token[T, C]{Kind: n.Kind(), Val: chars[:end]}, chars[end:], nil
}

func (l *Lexer[T, C]) Tokenize(chars []C) ([]Token[T, C], error) {
	var tokens []Token[T, C]
	for {
		t, next, err := l.Next(chars)
		if err != nil {
			return nil, err
		}
		if _, ok := l.ignored[t.Kind]; !ok {
			tokens = append(tokens, *t)
		}
		chars = next
		if t.Kind == l.eos {
			return tokens, nil
		}
	}
}
