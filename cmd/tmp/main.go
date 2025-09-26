package main

import (
	"container/heap"
	"fmt"
)

// MinHeap is a generic min-heap that can store any type T.
type MinHeap[T any] struct {
	data []T
	less func(a, b T) bool
}

func (h *MinHeap[T]) Len() int           { return len(h.data) }
func (h *MinHeap[T]) Less(i, j int) bool { return h.less(h.data[i], h.data[j]) }
func (h *MinHeap[T]) Swap(i, j int)      { h.data[i], h.data[j] = h.data[j], h.data[i] }

// Push adds an element to the heap.
func (h *MinHeap[T]) Push(x any) {
	h.data = append(h.data, x.(T))
}

// Pop removes the highest-priority element.
func (h *MinHeap[T]) Pop() any {
	n := len(h.data)
	item := h.data[n-1]
	h.data = h.data[0 : n-1]
	return item
}

// An Item is something we manage in a priority queue.
type Item struct {
	value    string // The value of the item; arbitrary.
	priority int    // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
}

// This example creates a PriorityQueue with some items, adds and manipulates an item,
// and then removes the items in priority order.
func main() {
	// Some items and their priorities.
	pq := &MinHeap[Item]{
		data: []Item{},
		less: func(a, b Item) bool {
			// The priority queue will be based on the 'priority' field.
			return a.priority < b.priority
		},
	}
	heap.Init(pq)

	items := map[string]int{
		"banana": 3, "apple": 2, "pear": 4,
	}

	for value, priority := range items {
		item := Item{value, priority}
		heap.Push(pq, item)
	}

	// Insert a new item and then modify its priority.
	item := Item{
		value:    "orange",
		priority: 3,
	}
	heap.Push(pq, item)

	// Take the items out; they arrive in decreasing priority order.
	for pq.Len() > 0 {
		item := heap.Pop(pq).(Item)
		fmt.Printf("%.2d:%s ", item.priority, item.value)
	}
}
