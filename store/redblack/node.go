package redblack

// Node is a redblack node in the tree.
type Node struct {
	id     int64
	value  interface{}
	parent *Node
	left   *Node
	right  *Node
	color  Color
}

// IsRoot returns true if the node has no parents which makes it a root node.
func (n *Node) IsRoot() bool {
	return n.Parent() == nil
}

// IsLeaf returns true if the node is a leaf node.
func (n *Node) IsLeaf() bool {
	return n.right == nil && n.left == nil
}

// Parent returns the nodes parent.
//   (parent)
//       \
//      (child node)
func (n *Node) Parent() *Node {
	return n.parent
}

// Grandparent returns the nodes grandparent.
//   (grandparent)
//   /           \
//  ()           (parent)
//                  \
//                  (child node)
func (n *Node) Grandparent() *Node {
	parent := n.Parent()
	if parent != nil {
		return parent.parent
	}
	return nil
}

// Uncle returns the nodes uncle.
//   (grandparent)
//   /           \
// (uncle)      (parent)
//                  \
//                  (child node)
func (n *Node) Uncle() *Node {
	parent := n.Parent()
	grandparent := n.Grandparent()

	if parent != nil && grandparent != nil {
		switch parent {
		case grandparent.left:
			return grandparent.right
		case grandparent.right:
			return grandparent.left
		}
	}

	return nil
}

// Sibling returns the nodes sibling.
//           (parent)
//          /        \
//      (sister)     (child node)
func (n *Node) Sibling() *Node {
	parent := n.Parent()

	if parent != nil {
		switch n {
		case parent.left:
			return parent.right
		case parent.right:
			return parent.left
		}
	}

	return nil
}
