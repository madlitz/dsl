package mydsl

import (
	"github.com/madlitz/go-dsl"
)

// NodeType represents the type of a node in the AST.
const (
	NODE_ASSIGNMENT dsl.NodeType = "ASSIGNMENT"
	NODE_CALL       dsl.NodeType = "CALL"
	NODE_EXPRESSION dsl.NodeType = "EXPRESSION"
	NODE_TERMINAL   dsl.NodeType = "TERMINAL"
	NODE_COMMENT    dsl.NodeType = "COMMENT"
)

func NewNodeSet() dsl.NodeSet {
	return dsl.NewNodeSet(
		NODE_COMMENT,
		NODE_EXPRESSION,
		NODE_ASSIGNMENT,
		NODE_TERMINAL,
		NODE_CALL,
	)
}
