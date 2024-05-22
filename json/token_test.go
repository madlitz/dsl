package json

import (
	"github.com/madlitz/go-dsl"
)

func NewTokenSet() dsl.TokenSet {
	return dsl.NewTokenSet(
		"NUMBER",
		"STRING",
		"TRUE",
		"FALSE",
		"NULL",
		"OPEN_ARRAY",
		"CLOSE_ARRAY",
		"OPEN_OBJECT",
		"CLOSE_OBJECT",
		"EOF",
		"ILLEGAL",
	)
}
