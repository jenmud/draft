package redblack

import "fmt"

// find walks the tree searching for the node which has `id`.
func find(id int64, node *Node) *Node {
	if node == nil {
		return nil
	}

	if node.id == id {
		return node
	}

	if id > node.id {
		return find(id, node.right)
	}

	return find(id, node.left)
}

// insert inserts the items using a binary search tree insert.
func insert(root, node *Node) *Node {
	switch {
	// Go right
	case node.id > root.id:
		if root.right == nil { // root.right is a leaf so set it to node
			root.right = node
			root.right.color = Red
			node.parent = root
			return node
		}

		return insert(root.right, node)

	// Go left
	case node.id < root.id:
		if root.left == nil { // root.left is a leaf so set it to node
			root.left = node
			root.left.color = Red
			node.parent = root
			return node
		}

		return insert(root.left, node)
	}

	return node
}

// TreePrint pretty prints the tree.
func TreePrint(node *Node, indent string, last bool) {
	if node != nil {
		fmt.Printf(indent)
		if last {
			fmt.Printf("R----")
			indent += "    "
		} else {
			fmt.Printf("L----")
			indent += "|   "
		}

		fmt.Printf("(%d %s)\n", node.id, node.color)
		TreePrint(node.left, indent, false)
		TreePrint(node.right, indent, true)
	}
}
