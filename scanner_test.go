package dsl

import (
	"bufio"
	"bytes"
	"testing"
)

func TestScan(t *testing.T) {
	input := "a,b,c"
	scanFn := func(s *Scanner) Token {

		s.Expect(ExpectRune{
			Branches: []Branch{
				{Rn: ',', Fn: func(s *Scanner) { s.SkipRune() }},
			},
			Options: ScanOptions{Optional: true},
		})

		s.Expect(ExpectRune{
			Branches: []Branch{
				{Rn: 'a', Fn: func(s *Scanner) { s.Match([]Match{{Literal: "a", ID: "A"}}) }},
				{Rn: 'b', Fn: func(s *Scanner) { s.Match([]Match{{Literal: "b", ID: "B"}}) }},
				{Rn: 'c', Fn: func(s *Scanner) { s.Match([]Match{{Literal: "c", ID: "C"}}) }},
			},
		})

		s.Match([]Match{{Literal: "", ID: TOKEN_EOF}})
		return s.Exit()
	}
	s := newScanner(scanFn, bufio.NewReader(bytes.NewBufferString(input)), &dslNoLogger{})

	expectedTokens := []Token{
		{ID: "A", Literal: "a", Line: 1, Position: 1},
		{ID: "B", Literal: "b", Line: 1, Position: 3},
		{ID: "C", Literal: "c", Line: 1, Position: 5},
	}

	for i, expected := range expectedTokens {
		token, _, _ := s.scan()
		if token != expected {
			t.Errorf("Token %d: expected %v, got %v", i+1, expected, token)
		}
	}

	if token, _, _ := s.scan(); token.ID != TOKEN_EOF {
		t.Errorf("Expected EOF, got %v", token)
	}
}
