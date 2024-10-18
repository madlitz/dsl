package json

import (
	"github.com/dezlitz/dsl"
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

func Scan(s *dsl.Scanner) dsl.Token {
	// Skip all whitespace at the beginning of the input
	s.Expect(dsl.ExpectRune{
		Branches: []dsl.Branch{
			{Rn: ' ', Fn: skipWhitespace},
			{Rn: '\t', Fn: skipWhitespace},
			{Rn: '\n', Fn: skipWhitespace},
			{Rn: '\r', Fn: skipWhitespace},
		},
		Options: dsl.ExpectRuneOptions{Multiple: true, Optional: true},
	})

	s.Expect(dsl.ExpectRune{
		Branches: []dsl.Branch{
			{Rn: '{', Fn: func(s *dsl.Scanner) { s.Match([]dsl.Match{{Literal: "{", ID: TOKEN_LBRACE}}) }},
			{Rn: '}', Fn: func(s *dsl.Scanner) { s.Match([]dsl.Match{{Literal: "}", ID: TOKEN_RBRACE}}) }},
			{Rn: '[', Fn: func(s *dsl.Scanner) { s.Match([]dsl.Match{{Literal: "[", ID: TOKEN_LBRACKET}}) }},
			{Rn: ']', Fn: func(s *dsl.Scanner) { s.Match([]dsl.Match{{Literal: "]", ID: TOKEN_RBRACKET}}) }},
			{Rn: ':', Fn: func(s *dsl.Scanner) { s.Match([]dsl.Match{{Literal: ":", ID: TOKEN_COLON}}) }},
			{Rn: ',', Fn: func(s *dsl.Scanner) { s.Match([]dsl.Match{{Literal: ",", ID: TOKEN_COMMA}}) }},
			{Rn: '"', Fn: stringLiteral},
		},
		BranchRanges: []dsl.BranchRange{
			{StartRn: '0', EndRn: '9', Fn: number},
			{StartRn: 'a', EndRn: 'z', Fn: literal},
			{StartRn: 'A', EndRn: 'Z', Fn: literal},
		},
		Options: dsl.ExpectRuneOptions{Optional: true},
	})

	return s.Exit()
}

func skipWhitespace(s *dsl.Scanner) {
	s.SkipRune()
}

func stringLiteral(s *dsl.Scanner) {
	s.SkipRune() // Skip the opening quote
	s.ExpectNot(dsl.ExpectNotRune{
		Runes: []rune{
			'"',
		},
		Fn:      nil,
		Options: dsl.ExpectRuneOptions{Multiple: true, Optional: true},
	})

	s.Expect(dsl.ExpectRune{
		Branches: []dsl.Branch{
			{Rn: '"', Fn: nil}, // Closing quote
		},
	})
	s.SkipRune() // Skip the closing quote
	s.Match([]dsl.Match{{Literal: "", ID: TOKEN_STRING}})

}

func number(s *dsl.Scanner) {
	s.Expect(dsl.ExpectRune{
		BranchRanges: []dsl.BranchRange{
			{StartRn: '0', EndRn: '9', Fn: nil},
		},
		Options: dsl.ExpectRuneOptions{Multiple: true, Optional: true},
	})
	s.Match([]dsl.Match{{Literal: "", ID: TOKEN_NUMBER}})
}

func literal(s *dsl.Scanner) {
	s.Expect(dsl.ExpectRune{
		BranchRanges: []dsl.BranchRange{
			{StartRn: 'a', EndRn: 'z', Fn: nil},
			{StartRn: 'A', EndRn: 'Z', Fn: nil},
		},
		Options: dsl.ExpectRuneOptions{Multiple: true, Optional: true},
	})

	s.Match([]dsl.Match{
		{Literal: "true", ID: TOKEN_TRUE},
		{Literal: "false", ID: TOKEN_FALSE},
		{Literal: "null", ID: TOKEN_NULL},
	})
}
