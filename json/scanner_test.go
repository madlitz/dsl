package json

import (
	"github.com/madlitz/go-dsl"
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
		Options: dsl.ScanOptions{Multiple: true, Optional: true},
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
		Options: dsl.ScanOptions{Optional: true},
	})

	return s.Exit()
}

func skipWhitespace(s *dsl.Scanner) {
	s.SkipRune()
}

func stringLiteral(s *dsl.Scanner) {
	s.SkipRune() // Skip the opening quote
	s.Expect(dsl.ExpectRune{
		Branches: []dsl.Branch{
			{Rn: '"', Fn: nil}, // Closing quote
		},
		Options: dsl.ScanOptions{Multiple: true, Invert: true},
	})

	s.Match([]dsl.Match{{Literal: "", ID: TOKEN_STRING}})

	s.Expect(dsl.ExpectRune{
		Branches: []dsl.Branch{
			{Rn: '"', Fn: nil}, // Closing quote
		},
	})
	s.SkipRune() // Skip the closing quote
}

func number(s *dsl.Scanner) {
	s.Expect(dsl.ExpectRune{
		BranchRanges: []dsl.BranchRange{
			{StartRn: '0', EndRn: '9', Fn: nil},
		},
		Options: dsl.ScanOptions{Multiple: true, Optional: true},
	})
	s.Match([]dsl.Match{{Literal: "", ID: TOKEN_NUMBER}})
}

func literal(s *dsl.Scanner) {
	s.Expect(dsl.ExpectRune{
		BranchRanges: []dsl.BranchRange{
			{StartRn: 'a', EndRn: 'z', Fn: nil},
			{StartRn: 'A', EndRn: 'Z', Fn: nil},
		},
		Options: dsl.ScanOptions{Multiple: true, Optional: true},
	})

	s.Match([]dsl.Match{
		{Literal: "true", ID: TOKEN_TRUE},
		{Literal: "false", ID: TOKEN_FALSE},
		{Literal: "null", ID: TOKEN_NULL},
	})
}
