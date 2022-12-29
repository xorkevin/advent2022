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

	astarAdjacent[T comparable] struct {
		g       int
		hasPrev bool
		prev    T
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
	adjacent := map[T]astarAdjacent[T]{}
	for _, i := range start {
		open.Add(astarNode[T]{
			v: i,
			g: 0,
			h: heuristic(i, goal),
		})
		adjacent[i] = astarAdjacent[T]{
			g: 0,
		}
	}
	for current, ok := open.Remove(); ok; current, ok = open.Remove() {
		if current.v == goal {
			revpath := []T{current.v}
			for i, ok := adjacent[current.v]; ok; {
				if !i.hasPrev {
					break
				}
				revpath = append(revpath, i.prev)
				i, ok = adjacent[i.prev]
			}
			n := len(revpath)
			for i := 0; i < n/2; i++ {
				revpath[i], revpath[n-i-1] = revpath[n-i-1], revpath[i]
			}
			return revpath, current.g
		}

		for _, i := range neighbors(current.v) {
			ng := current.g + i.DG
			if g, ok := adjacent[i.Value]; ok && ng >= g.g {
				continue
			}
			adjacent[i.Value] = astarAdjacent[T]{
				g:       ng,
				hasPrev: true,
				prev:    current.v,
			}
			open.Add(astarNode[T]{
				v: i.Value,
				g: ng,
				h: heuristic(i.Value, goal),
			})
		}
	}
	return nil, -1
}
