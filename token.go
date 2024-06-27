// token.go defines what a Token is and the token ID interface
package dsl

// Line is the line of the source text the Token was found. Position is the
// position (or column) the Token was found. This information is used when
// displaying errors but could also be useful to the user for things like
// syntax highlighting and debugging if they were to implement it.
type Token struct {
	ID       TokenType
	Literal  string
	Line     int
	Position int
}

type TokenType string

const (
	TOKEN_UNKNOWN TokenType = "UNKNOWN"
	TOKEN_EOF     TokenType = "EOF"
)
