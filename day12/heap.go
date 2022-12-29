package main

import (
	"container/heap"
)

type (
	heapItem[K comparable, V any] struct {
		key   K
		value V
		index int
	}

	Heap[K comparable, V any] struct {
		tree  []K
		items map[K]*heapItem[K, V]
		less  func(a, b V) bool
	}
)

func NewHeap[K comparable, V any](less func(a, b V) bool) *Heap[K, V] {
	return &Heap[K, V]{
		tree:  nil,
		items: map[K]*heapItem[K, V]{},
		less:  less,
	}
}

func (h Heap[K, V]) Len() int {
	return len(h.tree)
}

func (h Heap[K, V]) Less(i, j int) bool {
	return h.less(h.items[h.tree[i]].value, h.items[h.tree[j]].value)
}

func (h Heap[K, V]) Swap(i, j int) {
	a := h.tree[i]
	b := h.tree[j]
	h.tree[i], h.tree[j] = b, a
	h.items[b].index = i
	h.items[a].index = j
}

func (h *Heap[K, V]) Push(x any) {
	n := len(h.tree)
	item := x.(*heapItem[K, V])
	item.index = n
	h.tree = append(h.tree, item.key)
	h.items[item.key] = item
}

func (h *Heap[K, V]) Pop() any {
	n := len(h.tree)
	k := h.tree[n-1]
	h.tree = h.tree[:n-1]
	item := h.items[k]
	item.index = -1
	delete(h.items, k)
	return item
}

func (h *Heap[K, V]) Upsert(k K, v V) {
	item, ok := h.items[k]
	if !ok {
		heap.Push(h, &heapItem[K, V]{
			key:   k,
			value: v,
		})
		return
	}
	item.value = v
	heap.Fix(h, item.index)
}

func (h *Heap[K, V]) Remove() (K, V, bool) {
	if len(h.tree) == 0 {
		var k K
		var v V
		return k, v, false
	}
	item := heap.Pop(h).(*heapItem[K, V])
	return item.key, item.value, true
}

func (h Heap[K, V]) Has(k K) bool {
	_, ok := h.items[k]
	return ok
}

type (
	Set[T comparable] struct {
		vals map[T]struct{}
	}
)

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		vals: map[T]struct{}{},
	}
}

func (s Set[T]) Has(v T) bool {
	_, ok := s.vals[v]
	return ok
}

func (s Set[T]) Add(v T) {
	s.vals[v] = struct{}{}
}

func (s Set[T]) Remove(v T) {
	delete(s.vals, v)
}
