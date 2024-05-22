package mydsl

import (
	"github.com/madlitz/go-dsl"
)

func NewNodeSet() dsl.NodeSet {
	return dsl.NewNodeSet(
		"COMMENT",
		"EXPRESSION",
		"ASSIGNMENT",
		"TERMINAL",
		"CALL",
	)
}
