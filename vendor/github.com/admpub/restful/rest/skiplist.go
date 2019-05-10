package rest

import (
	"math/rand"
	"time"
)

// SkipList is a great DataStructure for creating ordered Sets/Maps/Lists
// In that way, it has the same performance/complexity as a B-Tree, O(lg n)
// for operations such as search, insert, remove, update.
// However, it is much simpler to implement than a B-Tree
// To learn more about Skip lists see: http://epaperpress.com/sortsearch/download/skiplist.pdf

// 32 levels, counting from 0
const maxHeight = 31

// The skiplist struct itself
type skipList struct {
	height int
	head   *skipListNode
}

// A node representation
type skipListNode struct {
	ttl  time.Time       // we will use this for comparing nodes, and setting order
	key  string          // important to remove elements from other structures in the cache (lru and map)
	next []*skipListNode // pointers to the next nodes. One for each level
}

func newSkipList() *skipList {
	head := &skipListNode{
		next: make([]*skipListNode, maxHeight),
	}

	return &skipList{
		height: 0,
		head:   head,
	}
}

// Insert a node to the Skip list
func (s *skipList) insert(key string, ttl time.Time) *skipListNode {

	level := 0

	// New random seed
	rand.Seed(time.Now().UnixNano())

	// Like flipping a coin up to the maximum height
	// Level will have a value between 0 and 31
	for level < maxHeight && rand.Intn(2) == 1 {

		level++

		if level > s.height {
			s.height = level
			break
		}
	}

	node := &skipListNode{
		ttl:  ttl,
		key:  key,
		next: make([]*skipListNode, level+1),
	}

	// Get the Head
	current := s.head

	// Start from the top as any search
	for i := s.height; i >= 0; i-- {

		for ; current.next[i] != nil; current = current.next[i] {

			// If the ttl of the next element is > than the element to be inserted,
			// go down one level
			if current.next[i].ttl.Sub(node.ttl) > 0 {
				break
			}

		} //end for

		// We just care if we are at the right level or less
		if i <= level {
			node.next[i] = current.next[i]
			current.next[i] = node
		}

	}

	return node

}

// Remove a node from the Skip list
func (s *skipList) remove(node *skipListNode) {

	if node == nil {
		return
	}

	current := s.head

	// Start from the top
	for i := s.height; i >= 0; i-- {

		// If next is nil, move to the next level
		for ; current.next[i] != nil; current = current.next[i] {

			// If the ttl of the next element is > than the element to be removed,
			// go down one level
			if current.next[i].ttl.Sub(node.ttl) > 0 {
				break
			}

			// If current next points to the node we are trying to remove,
			// change pointers, so current.next will point to node.next
			if current.next[i] == node {
				current.next[i] = node.next[i]
				break
			}

		} // end for

	}

}

/*
func (s *skipList) debug() {

	now := time.Now()
	fmt.Println("")

	current := s.head

	for i := s.height; i >= 0; i-- {

		for ; current.next[i] != nil; current = current.next[i] {

			diff := int64(current.ttl.Sub(now) / time.Millisecond)
			fmt.Print(diff)
			fmt.Print("-")
		}

		fmt.Println("")

	}

}
*/
