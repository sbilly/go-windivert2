package utils

// IPTree is a simple IP tree implementation for IP lookup
type IPTree struct {
	root *node
}

type node struct {
	left   *node
	right  *node
	isLeaf bool
}

// NewIPTree creates a new IP tree
func NewIPTree() *IPTree {
	return &IPTree{
		root: &node{},
	}
}

// Insert inserts an IP address into the tree
func (t *IPTree) Insert(ip []byte) {
	current := t.root
	for _, b := range ip {
		for i := 7; i >= 0; i-- {
			bit := (b >> uint(i)) & 1
			if bit == 0 {
				if current.left == nil {
					current.left = &node{}
				}
				current = current.left
			} else {
				if current.right == nil {
					current.right = &node{}
				}
				current = current.right
			}
		}
	}
	current.isLeaf = true
}

// Lookup checks if an IP address exists in the tree
func (t *IPTree) Lookup(ip []byte) bool {
	current := t.root
	for _, b := range ip {
		for i := 7; i >= 0; i-- {
			if current == nil {
				return false
			}
			if current.isLeaf {
				return true
			}
			bit := (b >> uint(i)) & 1
			if bit == 0 {
				current = current.left
			} else {
				current = current.right
			}
		}
	}
	return current != nil && current.isLeaf
}
