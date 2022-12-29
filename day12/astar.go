package main

type (
	astarNode struct {
		g, h int
	}

	AStarEdge[K comparable] struct {
		Key K
		DG  int
	}
)

func astarNodeLess(a, b astarNode) bool {
	af := a.g + a.h
	bf := b.g + b.h
	if af == bf {
		return a.g < b.g
	}
	return af < bf
}

func AStarSearch[K comparable](start []K, goal K, neighbors func(k K) []AStarEdge[K], heuristic func(a, b K) int) ([]K, int) {
	open := NewHeap[K](astarNodeLess)
	gscore := map[K]int{}
	adjacent := map[K]K{}
	for _, i := range start {
		open.Upsert(i, astarNode{
			g: 0,
			h: heuristic(i, goal),
		})
		gscore[i] = 0
	}
	for k, v, ok := open.Remove(); ok; k, v, ok = open.Remove() {
		if k == goal {
			revpath := []K{k}
			for i, ok := adjacent[k]; ok; i, ok = adjacent[i] {
				revpath = append(revpath, i)
			}
			n := len(revpath)
			for i := 0; i < n/2; i++ {
				revpath[i], revpath[n-i-1] = revpath[n-i-1], revpath[i]
			}
			return revpath, v.g
		}

		for _, i := range neighbors(k) {
			ng := v.g + i.DG
			if g, ok := gscore[i.Key]; ok && ng >= g {
				continue
			}
			adjacent[i.Key] = k
			gscore[i.Key] = ng
			open.Upsert(i.Key, astarNode{
				g: ng,
				h: heuristic(i.Key, goal),
			})
		}
	}
	return nil, -1
}
