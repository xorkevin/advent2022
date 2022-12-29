package main

type (
	astarNode[T comparable] struct {
		v    T
		g, h int
	}

	AStarEdge[T comparable] struct {
		Value T
		DG    int
	}
)

func astarNodeLess[T comparable](a, b astarNode[T]) bool {
	af := a.g + a.h
	bf := b.g + b.h
	if af == bf {
		return a.g < b.g
	}
	return af < bf
}

func AStarSearch[T comparable](start []T, goal T, neighbors func(k T) []AStarEdge[T], heuristic func(a, b T) int) ([]T, int) {
	open := NewHeap(astarNodeLess[T])
	gscore := map[T]int{}
	adjacent := map[T]T{}
	for _, i := range start {
		open.Add(astarNode[T]{
			v: i,
			g: 0,
			h: heuristic(i, goal),
		})
		gscore[i] = 0
	}
	for current, ok := open.Remove(); ok; current, ok = open.Remove() {
		if current.v == goal {
			revpath := []T{current.v}
			for i, ok := adjacent[current.v]; ok; i, ok = adjacent[i] {
				revpath = append(revpath, i)
			}
			n := len(revpath)
			for i := 0; i < n/2; i++ {
				revpath[i], revpath[n-i-1] = revpath[n-i-1], revpath[i]
			}
			return revpath, current.g
		}

		for _, i := range neighbors(current.v) {
			ng := current.g + i.DG
			if g, ok := gscore[i.Value]; ok && ng >= g {
				continue
			}
			adjacent[i.Value] = current.v
			gscore[i.Value] = ng
			open.Add(astarNode[T]{
				v: i.Value,
				g: ng,
				h: heuristic(i.Value, goal),
			})
		}
	}
	return nil, -1
}
