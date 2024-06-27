package json

import (
	"github.com/madlitz/go-dsl"
)

// NodeType represents the type of a node in the AST.
const (
	NODE_OBJECT dsl.NodeType = "OBJECT"
	NODE_ARRAY  dsl.NodeType = "ARRAY"
	NODE_MEMBER dsl.NodeType = "MEMBER"
	NODE_VALUE  dsl.NodeType = "VALUE"
)

func NewNodeSet() dsl.NodeSet {
	return dsl.NewNodeSet(
		NODE_OBJECT,
		NODE_ARRAY,
		NODE_MEMBER,
		NODE_VALUE,
	)
}
