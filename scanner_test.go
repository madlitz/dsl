package dsl

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestScan(t *testing.T) {
	input := "a,b,c"
	scanFn := func(s *Scanner) Token {

		s.Expect(ExpectRune{
			Branches: []Branch{
				{Rn: ',', Fn: func(s *Scanner) { s.SkipRune() }},
			},
			Options: ExpectRuneOptions{Optional: true},
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

func TestPeekScan(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "Integer",
			input: "123",
			expected: []Token{
				{ID: "NUMBER", Literal: "123", Line: 1, Position: 1},
			},
		},
		{
			name:  "FloatingPoint",
			input: "123.456",
			expected: []Token{
				{ID: "NUMBER", Literal: "123.456", Line: 1, Position: 1},
			},
		},
		{
			name:  "FloatingPointSingleDecimal",
			input: "123.4",
			expected: []Token{
				{ID: "NUMBER", Literal: "123.4", Line: 1, Position: 1},
			},
		},
		{
			name:  "ArraySpread",
			input: "123..456",
			expected: []Token{
				{ID: "NUMBER", Literal: "123", Line: 1, Position: 1},
				{ID: "SPREAD", Literal: "..", Line: 1, Position: 4},
				{ID: "NUMBER", Literal: "456", Line: 1, Position: 6},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanFn := func(s *Scanner) Token {
				// At this point we could have a number, spread operator, or EOF
				s.Expect(ExpectRune{
					Branches: []Branch{
						{Rn: rune(0), Fn: func(s *Scanner) {
							s.Match([]Match{{Literal: "", ID: TOKEN_EOF}})
						}},
					},
					Options: ExpectRuneOptions{Optional: true},
				})

				// At this point if not EOF, we expect a number or spread operator
				s.Expect(ExpectRune{
					Branches: []Branch{
						{Rn: '.', Fn: func(s *Scanner) {
							s.Expect(ExpectRune{
								Branches: []Branch{
									{Rn: '.', Fn: func(s *Scanner) {
										s.Match([]Match{{Literal: "..", ID: "SPREAD"}})
									}},
								},
							})
						}},
					},
					Options: ExpectRuneOptions{Optional: true},
				})

				// At this point if not EOF or spread operator, we expect a number
				s.Expect(ExpectRune{
					BranchRanges: []BranchRange{
						{StartRn: '0', EndRn: '9', Fn: nil},
					},
					Options: ExpectRuneOptions{Multiple: true},
				})

				// At this point we have a number, so we need to check for a decimal point
				s.Expect(ExpectRune{
					Branches: []Branch{
						{Rn: '.', Fn: func(s *Scanner) {
							s.Expect(ExpectRune{
								Branches: []Branch{
									{Rn: '.', Fn: func(s *Scanner) {
										// If we have a second decimal point, we know the number is an integer
										s.Match([]Match{{Literal: "", ID: "NUMBER"}})
									}},
								},
								Options: ExpectRuneOptions{Peek: true, Optional: true},
							})
						}},
					},
					Options: ExpectRuneOptions{Peek: true, Optional: true},
				})

				// At this point we have a number and we know there is no spread operator
				s.Expect(ExpectRune{
					Branches: []Branch{
						{Rn: '.', Fn: func(s *Scanner) {
							// We have a decimal point, so we must have more digits
							s.Expect(ExpectRune{
								BranchRanges: []BranchRange{
									{StartRn: '0', EndRn: '9', Fn: nil},
								},
								Options: ExpectRuneOptions{Multiple: true},
							})
						}},
					},
					Options: ExpectRuneOptions{Optional: true},
				})

				s.Match([]Match{{Literal: "", ID: "NUMBER"}})
				return s.Exit()
			}

			logfilename := fmt.Sprintf("logs/TestPeekScan_%s.log", tt.name)
			logfile, fileErr := os.Create(logfilename)
			if fileErr != nil {
				t.Fatal("Error: Could not create log file " + logfilename + ": " + fileErr.Error())
			}
			defer logfile.Close()
			logger := &dslLogger{logger: log.New(logfile, "", 0)}

			s := newScanner(scanFn, bufio.NewReader(bytes.NewBufferString(tt.input)), logger)

			for i, expectedToken := range tt.expected {
				token, _, err := s.scan()
				if err != nil {
					t.Fatalf("Unexpected error at token %d: %v", i+1, err)
				}
				if token != expectedToken {
					t.Errorf("Token %d: expected %v, got %v", i+1, expectedToken, token)
				}
			}

			// Ensure EOF is reached
			token, _, err := s.scan()
			if err != nil {
				t.Fatalf("Unexpected error at EOF: %v", err)
			}
			if token.ID != TOKEN_EOF {
				t.Errorf("Expected EOF, got %v", token)
			}
		})
	}
}
