package json

import (
	"github.com/madlitz/go-dsl"
)

func NewNodeSet() dsl.NodeSet {
	return dsl.NewNodeSet(
		"OBJECT",
		"ARRAY",
		"KEY",
		"VALUE",
	)
}
