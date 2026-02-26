// btree.go
package main

import "fmt"

// B-Tree (CLRS) with minimum degree t.
// - Max keys per node: 2t-1
// - Max children per node: 2t
type BTree struct {
	t    int
	root *Node
}

type Node struct {
	leaf     bool
	keys     []int
	children []*Node
}

func NewBTree(t int) *BTree {
	if t < 2 {
		panic("t must be >= 2")
	}
	return &BTree{
		t: t,
		root: &Node{
			leaf:     true,
			keys:     make([]int, 0, 2*t-1),
			children: make([]*Node, 0, 2*t),
		},
	}
}

// Search returns true if key exists.
func (bt *BTree) Search(k int) bool {
	return search(bt.root, k)
}

func search(n *Node, k int) bool {
	i := 0
	for i < len(n.keys) && k > n.keys[i] {
		i++
	}
	if i < len(n.keys) && n.keys[i] == k {
		return true
	}
	if n.leaf {
		return false
	}
	return search(n.children[i], k)
}

// Insert inserts a key into the B-Tree.
func (bt *BTree) Insert(k int) {
	r := bt.root
	if len(r.keys) == 2*bt.t-1 {
		// Root is full: grow tree height
		s := &Node{
			leaf:     false,
			keys:     make([]int, 0, 2*bt.t-1),
			children: make([]*Node, 0, 2*bt.t),
		}
		s.children = append(s.children, r)
		bt.splitChild(s, 0)
		bt.root = s
		bt.insertNonFull(s, k)
		return
	}
	bt.insertNonFull(r, k)
}

func (bt *BTree) insertNonFull(x *Node, k int) {
	i := len(x.keys) - 1
	if x.leaf {
		// Insert into sorted keys slice
		x.keys = append(x.keys, 0)
		for i >= 0 && k < x.keys[i] {
			x.keys[i+1] = x.keys[i]
			i--
		}
		x.keys[i+1] = k
		return
	}

	// Find child index to descend
	for i >= 0 && k < x.keys[i] {
		i--
	}
	i++

	// If child is full, split first
	if len(x.children[i].keys) == 2*bt.t-1 {
		bt.splitChild(x, i)
		// After split, decide which child to descend
		if k > x.keys[i] {
			i++
		}
	}
	bt.insertNonFull(x.children[i], k)
}

// splitChild splits x.children[i] (which must be full) into two nodes.
// Median key moves up into parent x at position i.
func (bt *BTree) splitChild(x *Node, i int) {
	t := bt.t
	y := x.children[i] // full child
	z := &Node{
		leaf:     y.leaf,
		keys:     make([]int, 0, 2*t-1),
		children: make([]*Node, 0, 2*t),
	}

	// y has 2t-1 keys: [0 .. 2t-2]
	// median index = t-1
	median := y.keys[t-1]

	// z takes keys [t .. 2t-2] => t-1 keys
	z.keys = append(z.keys, y.keys[t:]...)
	// y keeps keys [0 .. t-2] => t-1 keys
	y.keys = y.keys[:t-1]

	// Move children if internal
	if !y.leaf {
		// y has 2t children: [0 .. 2t-1]
		z.children = append(z.children, y.children[t:]...)
		y.children = y.children[:t]
	}

	// Insert z as new child right after y
	x.children = append(x.children, nil)
	copy(x.children[i+2:], x.children[i+1:])
	x.children[i+1] = z

	// Insert median key into x.keys at position i
	x.keys = append(x.keys, 0)
	copy(x.keys[i+1:], x.keys[i:])
	x.keys[i] = median
}

// Traverse prints keys in sorted order.
func (bt *BTree) Traverse() {
	traverse(bt.root)
	fmt.Println()
}

func traverse(n *Node) {
	i := 0
	for i < len(n.keys) {
		if !n.leaf {
			traverse(n.children[i])
		}
		fmt.Print(n.keys[i], " ")
		i++
	}
	if !n.leaf {
		traverse(n.children[i])
	}
}

func main() {
	bt := NewBTree(3) // minimum degree t=3 => max keys per node = 5

	nums := []int{10, 20, 5, 6, 12, 30, 7, 17, 3, 2, 4, 25, 26, 27, 28}
	for _, v := range nums {
		bt.Insert(v)
	}

	bt.Traverse()
	fmt.Println("Search 12:", bt.Search(12))
	fmt.Println("Search 99:", bt.Search(99))
}