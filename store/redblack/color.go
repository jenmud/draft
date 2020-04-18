package redblack

// Color is a redback node color.
type Color int

const (
	Black Color = iota
	Red
)

// String returns the color enum as a string.
func (c Color) String() string {
	return [...]string{"black", "red"}[c]
}

// NodeColor returns the nodes color or defaulting to black.
func NodeColor(node *Node) Color {
	if node == nil {
		return Black
	}
	return node.color
}
