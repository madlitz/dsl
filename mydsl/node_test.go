package mydsl

import(
    "github.com/deslittle/go-dsl"
)

func NewNodeSet() dsl.NodeSet{
    return dsl.NewNodeSet(
		"COMMENT",
		"EXPRESSION",
		"ASSIGNMENT",
		"TERMINAL",
		"CALL",
	)
}