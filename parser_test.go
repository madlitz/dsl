package dsl

import "testing"

// mockScanner simulates the Scanner for testing purposes
type mockScanner struct {
	tokens []Token
	index  int
}

func (m *mockScanner) scan() (Token, *Error) {
	if m.index >= len(m.tokens) {
		return Token{ID: TOKEN_EOF, Literal: "EOF", Line: 0, Position: 0}, nil
	}
	token := m.tokens[m.index]
	m.index++
	return token, nil
}

func (m *mockScanner) newError(code ErrorCode, err error) *Error {
	return &Error{Code: code, Message: err.Error()}
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
	parser := &Parser{
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

	parser.Expect(ExpectToken{
		Branches: []BranchToken{
			{Id: "a", Fn: func(p *Parser) {}},
			{Id: "b", Fn: func(p *Parser) {}},
		},
		Options: ParseOptions{Multiple: true},
	})

	parser.AddTokens()
	parser.AddNode(NODE_ROOT)

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
