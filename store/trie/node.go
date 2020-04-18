package trie

import (
	"fmt"
	"sync"
)

// Node is a trie node.
type Node struct {
	lock   sync.RWMutex
	unique bool
	leaf   bool
	childs map[rune]*Node
	parent *Node
	key    string
	value  interface{}
}

func (n *Node) String() string {
	return fmt.Sprintf("(key: %s, value: %v)", n.key, n.value)
}

// Add adds a key and value pair to the trie.
func (n *Node) Add(key string, value interface{}, unique bool) error {
	n.lock.Lock()
	defer n.lock.Unlock()

	var parent *Node
	var node *Node
	var ok bool

	for _, r := range []rune(key) {
		node, ok = n.childs[r]
		if !ok {
			node = NewNode()
			node.parent = n
			node.key = string(r)
			n.childs[r] = node
		}
		parent = n
		n = node
	}

	if n.unique {
		return index.NewConstraintError(fmt.Sprintf("Unique constraint error on %s, failed to add %q at key %q", n, value, key))
	}

	n.parent = parent
	n.key = key
	n.value = value
	n.leaf = true
	n.unique = unique
	return nil
}

// Get returns the value found at key.
func (n *Node) Get(key string) (interface{}, error) {
	n.lock.RLock()
	defer n.lock.RUnlock()

	var ok bool

	for _, r := range []rune(key) {
		n, ok = n.childs[r]
		if !ok {
			return nil, index.NewNotFound(fmt.Sprintf("%q not found", key))
		}
	}

	if !n.leaf {
		return nil, index.NewNotFound(fmt.Sprintf("%q not found", key))
	}

	return n.value, nil
}

// fetchAll takes a key returns all the nodes in
// order till the node matching key is found.
func (n *Node) fetchAll(key string) []*Node {
	n.lock.RLock()
	defer n.lock.RUnlock()

	items := make([]*Node, len(key))
	for i, r := range []rune(key) {
		node, ok := n.childs[r]
		if !ok {
			return []*Node{}
		}
		items[i] = node
		n = node
	}

	return items
}

// Remove removes the key and value pair from the trie.
func (n *Node) Remove(key string) error {
	nodes := n.fetchAll(key)
	if len(nodes) == 0 {
		return fmt.Errorf("Did not find %q", key)
	}

	n.lock.Lock()
	defer n.lock.Unlock()

	runes := []rune(key)

	for i, node := range nodes {
		if len(node.childs) >= 2 {
			if (i == len(nodes)-1) && node.leaf {
				node.leaf = false
				node.unique = false
			}
			continue
		} else {
			parent := node.parent
			delete(parent.childs, runes[i])
			continue
		}
	}

	return nil
}

// NewNode returns a new empty Trie node.
func NewNode() *Node {
	return &Node{
		childs: make(map[rune]*Node),
	}
}
