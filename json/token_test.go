package json

import (
	"github.com/madlitz/go-dsl"
)

const (
	TOKEN_STRING   dsl.TokenType = "STRING"
	TOKEN_NUMBER   dsl.TokenType = "NUMBER"
	TOKEN_TRUE     dsl.TokenType = "TRUE"
	TOKEN_FALSE    dsl.TokenType = "FALSE"
	TOKEN_NULL     dsl.TokenType = "NULL"
	TOKEN_LBRACE   dsl.TokenType = "LBRACE"
	TOKEN_RBRACE   dsl.TokenType = "RBRACE"
	TOKEN_LBRACKET dsl.TokenType = "LBRACKET"
	TOKEN_RBRACKET dsl.TokenType = "RBRACKET"
	TOKEN_COLON    dsl.TokenType = "COLON"
	TOKEN_COMMA    dsl.TokenType = "COMMA"
	TOKEN_WS       dsl.TokenType = "WS"
	TOKEN_EOF      dsl.TokenType = "EOF"
)

func NewTokenSet() dsl.TokenSet {
	return dsl.NewTokenSet(
		TOKEN_STRING,
		TOKEN_NUMBER,
		TOKEN_TRUE,
		TOKEN_FALSE,
		TOKEN_NULL,
		TOKEN_LBRACE,
		TOKEN_RBRACE,
		TOKEN_LBRACKET,
		TOKEN_RBRACKET,
		TOKEN_COLON,
		TOKEN_COMMA,
		TOKEN_WS,
		TOKEN_EOF,
	)
}
