package graph

type FilterType int

const (
	// LABEL is used for filtering by label.
	LABEL FilterType = iota
	PROPERTY
)

type ItemType int

const (
	NODE ItemType = iota
	EDGE
)

// Filter is used for filtering and searching.
type Filter struct {
	// Type is the filter typle, example label, properties
	Type FilterType
	// Value is the filter criterial to match.
	Value Value
}
