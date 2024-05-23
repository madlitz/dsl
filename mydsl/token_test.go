package mydsl

import (
	"github.com/madlitz/go-dsl"
)

const (
	TOKEN_LITERAL     dsl.TokenType = "LITERAL"
	TOKEN_PLUS        dsl.TokenType = "PLUS"
	TOKEN_MINUS       dsl.TokenType = "MINUS"
	TOKEN_MULTIPLY    dsl.TokenType = "MULTIPLY"
	TOKEN_DIVIDE      dsl.TokenType = "DIVIDE"
	TOKEN_OPEN_PAREN  dsl.TokenType = "OPEN_PAREN"
	TOKEN_CLOSE_PAREN dsl.TokenType = "CLOSE_PAREN"
	TOKEN_ASSIGN      dsl.TokenType = "ASSIGN"
	TOKEN_VARIABLE    dsl.TokenType = "VARIABLE"
	TOKEN_COMMENT     dsl.TokenType = "COMMENT"
	TOKEN_NL          dsl.TokenType = "NL"
	TOKEN_WS          dsl.TokenType = "WS"
	TOKEN_EOF         dsl.TokenType = "EOF"
)

func NewTokenSet() dsl.TokenSet {
	return dsl.NewTokenSet(
		TOKEN_LITERAL,
		TOKEN_PLUS,
		TOKEN_MINUS,
		TOKEN_MULTIPLY,
		TOKEN_DIVIDE,
		TOKEN_OPEN_PAREN,
		TOKEN_CLOSE_PAREN,
		TOKEN_ASSIGN,
		TOKEN_VARIABLE,
		TOKEN_COMMENT,
		TOKEN_NL,
		TOKEN_WS,
		TOKEN_EOF,
	)
}
