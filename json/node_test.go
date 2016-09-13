package json

import(
    "github.com/deslittle/go-dsl"
)

func NewNodeSet() dsl.NodeSet{
    return dsl.NewNodeSet(
		"OBJECT",
		"ARRAY",
		"KEY",
		"VALUE",
	)
}