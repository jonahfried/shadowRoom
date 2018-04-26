package priorityqueue

import "shadowRoom/boundry"

// PriorityQueue is a datastructure that can be pushed to to store elems
// and Popped from to receive the lowest stored element.
type PriorityQueue []elem

type elem struct {
	val    boundry.Tile
	weight float64
}

// Push is a PriorityQueue method, adding a given value to the PriorityQueue with a given weight.
func (prior *PriorityQueue) Push(tl *boundry.Tile, weight float64) {

}
