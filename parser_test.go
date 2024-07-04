package dsl

import (
	"testing"
	"time"
)

// mockScanner simulates the Scanner for testing purposes
type mockScanner struct {
	tokens []Token
	index  int
}

func (m *mockScanner) scan() (Token, string, *Error) {
	if m.index >= len(m.tokens) {
		return Token{ID: TOKEN_EOF, Literal: "EOF", Line: 0, Position: 0}, "", nil
	}
	token := m.tokens[m.index]
	m.index++
	return token, tokensToLineString(m.tokens[0:m.index]), nil
}

// mockLogger simulates the Logger for testing purposes
type mockLogger struct{}

func (m *mockLogger) log(msg string, indent indent) {}

// TestExpect tests the Expect method of the Parser
func TestExpect(t *testing.T) {

	tests := []struct {
		name          string
		expectToken   ExpectToken
		expectedCount int
		expectedError bool
	}{
		{
			name: "Basic Expect",
			expectToken: ExpectToken{
				Branches: []BranchToken{
					{Id: "a", Fn: func(p *Parser) {}},
					{Id: "b", Fn: func(p *Parser) {}},
				},
			},
			expectedCount: 1,
			expectedError: false,
		},
		{
			name: "Expect with Multiple",
			expectToken: ExpectToken{
				Branches: []BranchToken{
					{Id: "a", Fn: func(p *Parser) {}},
					{Id: "b", Fn: func(p *Parser) {}},
				},
				Options: ParseOptions{Multiple: true},
			},
			expectedCount: 2,
			expectedError: false,
		},
		{
			name: "Expect with Invert",
			expectToken: ExpectToken{
				Branches: []BranchToken{
					{Id: "d", Fn: func(p *Parser) {}},
				},
				Options: ParseOptions{Invert: true},
			},
			expectedCount: 1,
			expectedError: false,
		},
		{
			name: "Expect with Error",
			expectToken: ExpectToken{
				Branches: []BranchToken{
					{Id: "d", Fn: func(p *Parser) {}},
				},
			},
			expectedCount: 0,
			expectedError: true,
		},
		{
			name: "Expect with Optional",
			expectToken: ExpectToken{
				Branches: []BranchToken{
					{Id: "d", Fn: func(p *Parser) {}},
				},
				Options: ParseOptions{Optional: true},
			},
			expectedCount: 0,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			parser := &Parser{
				s: &mockScanner{
					tokens: []Token{
						{ID: "a", Literal: "a", Line: 1, Position: 1},
						{ID: "b", Literal: "b", Line: 1, Position: 2},
					},
				},
				l: &mockLogger{},
			}

			parser.Expect(tt.expectToken)
			if len(parser.tokens) != tt.expectedCount {
				t.Errorf(
					"Unexpected token count: got %d, want %d",
					len(parser.tokens),
					tt.expectedCount,
				)
			}

			if parser.err && !tt.expectedError {
				t.Errorf("Expect returned an error when it should not have")
			}
			if !parser.err && tt.expectedError {
				t.Errorf("Expect did not return an error when it should have")
			}
		})
	}
}

func tokensToLineString(tokens []Token) string {
	// Create a line string from the tokens
	// If you dont have a token that covers the current position, add a space
	var lineString string
	var curPos = 1
	for _, token := range tokens {
		for i := curPos; i < token.Position; i++ {
			lineString += " "
		}
		lineString += token.Literal
		curPos = token.Position + len(token.Literal)
	}
	return lineString
}

// TestPeek tests the Peek method of the Parser
func TestPeek(t *testing.T) {
	s := &mockScanner{
		tokens: []Token{
			{ID: "a", Literal: "a", Line: 1, Position: 1},
			{ID: "b", Literal: "b", Line: 1, Position: 2},
			{ID: "c", Literal: "c", Line: 1, Position: 3},
		},
	}

	parser := &Parser{
		s: s,
		l: &mockLogger{},
	}

	called := false
	parser.Peek([]PeekToken{
		{
			IDs: []TokenType{"a", "b"},
			Fn:  func(*Parser) { called = true },
		},
	})

	if !called {
		t.Errorf("Peek function was not called when it should have been")
	}

	s.index = 0 // Reset scanner
	called = false
	parser.Peek([]PeekToken{
		{
			IDs: []TokenType{"a", "b", "c"},
			Fn:  func(*Parser) { called = true },
		},
	})

	if !called {
		t.Errorf("Peek function was not called when it should have been")
	}

	s.index = 0 // Reset scanner
	called = false
	parser.Peek([]PeekToken{
		{
			IDs: []TokenType{"a", "c", "b"},
			Fn:  func(*Parser) { called = true },
		},
	})

	if called {
		t.Errorf("Peek function was called when it should not have been")
	}
}

// TestAddNode tests the AddNode method of the Parser
func TestAddNode(t *testing.T) {
	ast := newAST()
	p := &Parser{
		ast: ast,
		s: &mockScanner{
			tokens: []Token{
				{ID: "a", Literal: "a", Line: 1, Position: 1},
				{ID: "b", Literal: "b", Line: 1, Position: 2},
				{ID: "c", Literal: "c", Line: 1, Position: 3},
			},
		},
		l: &mockLogger{},
	}

	p.Expect(ExpectToken{
		Branches: []BranchToken{
			{Id: "a", Fn: func(p *Parser) {}},
			{Id: "b", Fn: func(p *Parser) {}},
		},
		Options: ParseOptions{Multiple: true},
	})

	p.AddTokens()
	p.AddNode(NODE_ROOT)

	if len(ast.curNode.Children) != 1 || ast.curNode.Type != NODE_ROOT {
		t.Fatalf("Unexpected node in AST: got %v, want %v", ast.curNode.Type, NODE_ROOT)
	}

	if len(ast.curNode.Tokens) != 2 {
		t.Fatalf("Unexpected token count in node: got %d, want %d", len(ast.curNode.Tokens), 1)
	}

	if ast.curNode.Tokens[0].ID != "a" {
		t.Fatalf("Unexpected token in node: got %v, want %v", ast.curNode.Tokens[0].ID, "a")
	}

	if ast.curNode.Tokens[1].ID != "b" {
		t.Fatalf("Unexpected token in node: got %v, want %v", ast.curNode.Tokens[1].ID, "b")
	}

}

func TestParserInfiniteLoopDetection(t *testing.T) {

	s := &mockScanner{
		tokens: []Token{
			{ID: "a", Literal: "a", Line: 1, Position: 1},
			{ID: "b", Literal: "b", Line: 1, Position: 2},
			{ID: "c", Literal: "c", Line: 1, Position: 3},
		},
	}
	l := &mockLogger{}
	ast := newAST()

	var parseA, parseB func(*Parser)

	parseA = func(p *Parser) {
		p.Peek([]PeekToken{
			{
				IDs: []TokenType{"a", "b"},
				Fn:  parseB,
			},
		})
	}

	parseB = func(p *Parser) {
		p.Peek([]PeekToken{
			{
				IDs: []TokenType{"a", "b", "c"},
				Fn:  parseA,
			},
		})
	}

	parseFunc := func(p *Parser) (AST, []Error) {
		p.Call(parseA)
		return p.ast, p.errors
	}

	p := &Parser{
		fn:  parseFunc,
		s:   s,
		l:   l,
		ast: ast,
	}

	// Run the parser with a timeout
	done := make(chan bool)
	var errors []Error

	go func() {
		_, errors = p.fn(p)
		done <- true
	}()

	select {
	case <-done:
		// Check if the expected error was returned
		infiniteLoopErrorFound := false
		for _, err := range errors {
			if err.Code == ErrorInfiniteLoopDetected { // Assume this error code exists
				infiniteLoopErrorFound = true
				break
			}
		}
		if !infiniteLoopErrorFound {
			t.Error("Parser terminated without detecting the infinite loop")
		} else {
			t.Log("Infinite loop correctly detected and reported")
		}
	case <-time.After(2 * time.Second):
		t.Error("Test failed: parser entered an actual infinite loop")
	}

	// Additional checks can be performed here if needed
}
