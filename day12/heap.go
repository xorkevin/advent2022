package main

import (
	"container/heap"
)

type (
	Heap[T any] struct {
		tree []T
		less func(a, b T) bool
	}
)

func NewHeap[T any](less func(a, b T) bool) *Heap[T] {
	return &Heap[T]{
		tree: nil,
		less: less,
	}
}

func (h Heap[T]) Len() int {
	return len(h.tree)
}

func (h Heap[T]) Less(i, j int) bool {
	return h.less(h.tree[i], h.tree[j])
}

func (h Heap[T]) Swap(i, j int) {
	h.tree[i], h.tree[j] = h.tree[j], h.tree[i]
}

func (h *Heap[T]) Push(x any) {
	k := x.(T)
	h.tree = append(h.tree, k)
}

func (h *Heap[T]) Pop() any {
	n := len(h.tree)
	k := h.tree[n-1]
	h.tree = h.tree[:n-1]
	return k
}

func (h *Heap[T]) Add(v T) {
	heap.Push(h, v)
}

func (h *Heap[T]) Remove() (T, bool) {
	if len(h.tree) == 0 {
		var k T
		return k, false
	}
	k := heap.Pop(h).(T)
	return k, true
}
