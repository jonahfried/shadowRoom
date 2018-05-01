package priorityqueue

import (
	"container/heap"
)

// PriorityQueue is a datastructure that can be pushed to to store elems
// and Popped from to receive the lowest stored element.
type PriorityQueue []*Elem

// An Elem is used in priorityQueue to store value, priority, and ind in priorityQueue
type Elem struct {
	Value    int
	Priority float64
	Index    int
}

// Len returns the length of the PriorityQueue
func (prior PriorityQueue) Len() int { return len(prior) }

// Less returns whether index i of PriorityQueue has less priority than index j
func (prior PriorityQueue) Less(i, j int) bool {
	return (prior[i].Priority < prior[j].Priority)
}

func (prior PriorityQueue) Swap(i, j int) {
	prior[i], prior[j] = prior[j], prior[i]
	prior[i].Index, prior[j].Index = i, j
}

// Push adds an elem to the queue
func (prior *PriorityQueue) Push(x interface{}) {
	n := len(*prior)
	// var item *Elem
	var item = x.(*Elem)
	// item.value = x
	// item.priority = priority
	item.Index = n
	*prior = append(*prior, item)
}

// Pop removes the next element of the PriorityQueue
func (prior *PriorityQueue) Pop() interface{} {
	old := *prior
	n := len(old)
	item := old[n-1]
	item.Index = -1 // for safety
	*prior = old[0 : n-1]
	return item.Value

}

func (prior *PriorityQueue) update(item *Elem, value int, priority float64) {
	item.Value = value
	item.Priority = priority
	heap.Fix(prior, item.Index)
}
