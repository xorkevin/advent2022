package main

type (
	PairwiseDistances[T comparable] map[T]map[T]int

	Edge[T any] struct {
		A, B T
		C    int
	}
)

func (m PairwiseDistances[T]) EdgeCost(a, b T) (int, bool) {
	d, ok := m[a]
	if !ok {
		return 0, false
	}
	k, ok := d[b]
	if !ok {
		return 0, false
	}
	return k, true
}

func FloydWarshall[T comparable](nodes []T, edges []Edge[T]) map[T]map[T]int {
	dist := PairwiseDistances[T]{}
	for _, i := range edges {
		if _, ok := dist[i.A]; !ok {
			dist[i.A] = map[T]int{}
		}
		dist[i.A][i.B] = i.C
	}
	for _, i := range nodes {
		if _, ok := dist[i]; !ok {
			dist[i] = map[T]int{}
		}
		dist[i][i] = 0
	}
	for _, k := range nodes {
		for _, i := range nodes {
			for _, j := range nodes {
				cik, ok := dist.EdgeCost(i, k)
				if !ok {
					continue
				}
				ckj, ok := dist.EdgeCost(k, j)
				if !ok {
					continue
				}
				ck := cik + ckj
				if cij, ok := dist.EdgeCost(i, j); !ok || cij > ck {
					dist[i][j] = ck
				}
			}
		}
	}
	return dist
}
