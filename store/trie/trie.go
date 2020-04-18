package trie

import (
	"sync"
)

// Trie is a trie datastructure used for storing and retrieving text.
type Trie struct {
	lock  sync.RWMutex
	root  *Node
	count int
}

// Add adds a item to the Trie raising an error if the key has been
// set to unique and already exists.
func (t *Trie) Add(key string, value interface{}, unique bool) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	root := t.root

	err := root.Add(key, value, unique)
	if err != nil {
		return err
	}

	t.count++
	return err
}

// Get returns the value found at key.
func (t *Trie) Get(key string) (interface{}, error) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	root := t.root
	return root.Get(key)
}

// Remove removes the key and value pair from the trie.
func (t *Trie) Remove(key string) {
	t.lock.Lock()
	defer t.lock.Unlock()

	root := t.root

	err := root.Remove(key)
	if err != nil {
		return
	}

	t.count--
}

// Has returns true if the key is found in the trie.
func (t *Trie) Has(key string) bool {
	t.lock.RLock()
	defer t.lock.RUnlock()
	root := t.root
	_, err := root.Get(key)
	return err == nil
}

// walk recursively walks the nodes till it finds a
// leaf node and returns it.
func walk(n *Node, items chan string) {
	if n.leaf {
		items <- n.key
	}

	for _, node := range n.childs {
		walk(node, items)
	}
}

// Keys returns all the keys in the trie.
func (t *Trie) Keys() []string {
	t.lock.RLock()
	defer t.lock.RUnlock()

	root := t.root
	items := []string{}

	var walker func(n *Node)

	walker = func(n *Node) {
		if n.leaf {
			items = append(items, n.key)
		}

		for _, node := range n.childs {
			walker(node)
		}
	}

	walker(root)
	return items
}

// Count returns the number of keys in the trie.
func (t *Trie) Count() int {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.count
}

// New returns a new empty Trie.
func New() *Trie {
	root := NewNode()
	root.key = "."
	return &Trie{root: root}
}
