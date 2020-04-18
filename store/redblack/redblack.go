package redblack

import (
	"fmt"
	"sync"
)

/* Rules
1. Every node in T is either red or black.
2. The root node of T is black.
3. Every NULL node is black.
   (
       NULL nodes are the leaf nodes.
       They do not contain any keys.
       When we search for a key that is not present in the tree, we reach the NULL node.
    )
4. If a node is red, both of its children are black. This means no two nodes on a path
   can be red nodes.
5. Every path from a root node to a NULL node has the same number of black nodes.
*/

// Store is the redblack tree.
type Store struct {
	lock  sync.RWMutex
	root  *Node
	count int
}

// Count returns the number of items in the store.
func (s *Store) Count() int {
	return s.count
}

// Root returns the root node of the tree.
func (s *Store) Root() *Node {
	return s.root
}

// rightRotate right rotate node the node.
func (s *Store) rightRotate(x *Node) {
	// noop
	if x == nil {
		return
	}

	// noop as there is no left subtree
	if x.left == nil {
		return
	}

	y := x.left
	x.left = y.right

	if y.right != nil {
		y.right.parent = x
	}

	y.parent = x.parent

	if x.parent == nil {
		s.root = y
	} else {
		if x == x.parent.left {
			x.parent.left = y
		} else {
			x.parent.right = y
		}
	}

	y.right = x
	x.parent = y
}

// leftRotate left rotate node the node.
func (s *Store) leftRotate(x *Node) {
	if x == nil {
		// noop
		return
	}

	if x.right == nil {
		// noop as there is no right subtree
		return
	}

	y := x.right
	x.right = y.left

	if y.left != nil {
		y.left.parent = x
	}

	y.parent = x.parent

	if x.parent == nil {
		s.root = y
	} else {
		if x == x.parent.left {
			x.parent.left = y
		} else {
			x.parent.right = y
		}
	}

	y.left = x
	x.parent = y
}

// repairInsert repairs the redblack tree.
func (s *Store) repairInsert(z *Node) {
loop:
	for {
		switch {
		case z.parent == nil:
			fallthrough
		case z.parent.color == Black:
			fallthrough
		default:
			break loop
		case z.parent.color == Red:
			grandparent := z.parent.parent
			if z.parent == grandparent.left {
				y := grandparent.right
				if y != nil && y.color == Red {
					z.parent.color = Black
					y.color = Black
					grandparent.color = Red
					z = grandparent
				} else {
					// case 2
					if z == z.parent.right {
						z = z.parent
						s.leftRotate(z)
					}

					// case 3
					z.parent.color = Black
					grandparent.color = Red
					s.rightRotate(grandparent)
				}
			} else {
				y := grandparent.left
				if y != nil && y.color == Red {
					z.parent.color = Black
					y.color = Black
					grandparent.color = Red
					z = grandparent
				} else {
					// case 2
					if z == z.parent.left {
						z = z.parent
						s.rightRotate(z)
					}

					// case 3
					z.parent.color = Black
					grandparent.color = Red
					s.leftRotate(grandparent)
				}
			}
		}
	}

	s.root.color = Black
}

// Add adds a items to the store.
func (s *Store) Add(item interface{}) error {
	node := new(Node)
	node.id = item.ID()
	node.value = item

	s.lock.Lock()
	defer s.lock.Unlock()

	if s.root == nil {
		s.root = node
		s.count++
		return nil
	}

	node = insert(s.root, node)
	s.count++

	if node.Parent() == nil {
		node.color = Black
		return nil
	}

	if node.Grandparent() == nil {
		return nil
	}

	s.repairInsert(node)
	return nil
}

// transplant is a process of swapping nodes during the delete process.
func (s *Store) transplant(u *Node, v *Node) {
	if u.parent == nil {
		s.root = v
	} else {
		if u == u.parent.left {
			u.parent.left = v
		} else {
			u.parent.right = v
		}
	}

	if v != nil {
		v.parent = u.parent
	}
}

// getMin returns the node with minimum key starting at the
// subtree rooted at node x.
func (s *Store) getMin(x *Node) *Node {
	for {
		if x.left != nil {
			x = x.left
		} else {
			return x
		}
	}
}

// getMax returns the node with maximum key starting at the
// subtree rooted at node x.
func (s *Store) getMax(x *Node) *Node {
	for {
		if x.right != nil {
			x = x.right
		} else {
			return x
		}
	}
}

func (s *Store) removeCase1(node *Node) {
	if node.parent == nil {
		return
	}
	s.removeCase2(node)
}

func (s *Store) removeCase2(node *Node) {
	sibling := node.Sibling()
	if NodeColor(sibling) == Red {
		node.parent.color = Red
		sibling.color = Black
		if node == node.parent.left {
			s.leftRotate(node.parent)
		} else {
			s.rightRotate(node.parent)
		}
	}
	s.removeCase3(node)
}

func (s *Store) removeCase3(node *Node) {
	sibling := node.Sibling()
	if NodeColor(node.parent) == Black && NodeColor(sibling) == Black && NodeColor(sibling.left) == Black && NodeColor(sibling.right) == Black {
		sibling.color = Red
		s.removeCase1(node.Parent())
	} else {
		s.removeCase4(node)
	}
}

func (s *Store) removeCase4(node *Node) {
	sibling := node.Sibling()
	if NodeColor(node.parent) == Red && NodeColor(sibling) == Black && NodeColor(sibling.left) == Black && NodeColor(sibling.right) == Black {
		sibling.color = Red
		node.parent.color = Black
	} else {
		s.removeCase5(node)
	}
}

func (s *Store) removeCase5(node *Node) {
	sibling := node.Sibling()
	if node == node.parent.left && NodeColor(sibling) == Black && NodeColor(sibling.left) == Red && NodeColor(sibling.right) == Black {
		sibling.color = Red
		sibling.left.color = Black
		s.rightRotate(sibling)
	} else if node == node.parent.right && NodeColor(sibling) == Black && NodeColor(sibling.right) == Red && NodeColor(sibling.left) == Black {
		sibling.color = Red
		sibling.right.color = Black
		s.leftRotate(sibling)
	}
	s.removeCase6(node)
}

func (s *Store) removeCase6(node *Node) {
	sibling := node.Sibling()
	sibling.color = NodeColor(node.parent)
	node.parent.color = Black
	if node == node.parent.left && NodeColor(sibling.right) == Red {
		sibling.right.color = Black
		s.leftRotate(node.parent)
	} else if NodeColor(sibling.left) == Red {
		sibling.left.color = Black
		s.rightRotate(node.parent)
	}
}

// Remove remove the item from the tree.
func (s *Store) Remove(id int64) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	var child *Node
	node := find(id, s.root)

	// if node is nil, it is not found and there is nothing to do.
	if node == nil {
		return fmt.Errorf("Remove did not found node with id %d", id)
	}

	if node.left != nil && node.right != nil {
		pred := s.getMax(node.left)
		node.id = pred.id
		node.value = pred.value
		node = pred
	}

	if node.left == nil || node.right == nil {
		if node.right == nil {
			child = node.left
		} else {
			child = node.right
		}

		if node.color == Black {
			node.color = NodeColor(child)
			s.removeCase1(node)
		}

		s.transplant(node, child)
		if node.parent == nil && child != nil {
			child.color = Black
		}
	}

	s.count--
	return nil
}

// Has returns true if there is a node in the tree with `id`.
func (s *Store) Has(id int64) bool {
	node, _ := s.Get(id)
	return node != nil
}

// Get gets a item from the store.
func (s *Store) Get(id int64) (interface{}, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	node := find(id, s.root)
	if node == nil {
		return nil, fmt.Errorf("Could not find %d", id)
	}

	return node.value, nil
}

func walk(out chan interface{}, node *Node) {
	if node == nil {
		return
	}

	out <- node.value

	if node.left != nil {
		walk(out, node.left)
	}

	if node.right != nil {
		walk(out, node.right)
	}
}

// Iterate iterates over all the items in the store.
func (s *Store) Iterate() <-chan interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()

	out := make(chan interface{}, s.count)

	go func() {
		walk(out, s.root)
		close(out)
	}()

	return out
}

// New returns a new redblack tree.
func New() *Store {
	return &Store{count: 0}
}
