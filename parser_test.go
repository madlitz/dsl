package dsl

// import (
// 	"testing"
// )

// // MockScanner simulates the Scanner for testing purposes
// type MockScanner struct {
// 	tokens []Token
// 	index  int
// }

// func (m *MockScanner) scan() (Token, *Error) {
// 	if m.index >= len(m.tokens) {
// 		return Token{ID: TOKEN_EOF, Literal: "", Line: 0, Position: 0}, nil
// 	}
// 	token := m.tokens[m.index]
// 	m.index++
// 	return token, nil
// }

// func (m *MockScanner) newError(code ErrorCode, err error) *Error {
// 	return &Error{Code: code, Message: err.Error()}
// }

// // TestExpect tests the Expect method of the Parser
// func TestExpect(t *testing.T) {
// 	mockScanner := &MockScanner{
// 		tokens: []Token{
// 			{ID: "a", Literal: "a", Line: 1, Position: 1},
// 			{ID: "b", Literal: "b", Line: 1, Position: 2},
// 			{ID: "c", Literal: "c", Line: 1, Position: 3},
// 		},
// 	}

// 	parser := &Parser{
// 		s: mockScanner,
// 		l: &MockLogger{},
// 	}

// 	tests := []struct {
// 		name          string
// 		expectToken   ExpectToken
// 		expectedError bool
// 	}{
// 		{
// 			name: "Basic Expect",
// 			expectToken: ExpectToken{
// 				Branches: []BranchToken{
// 					{Id: "a", Fn: func(*Parser) {}},
// 				},
// 			},
// 			expectedError: false,
// 		},
// 		{
// 			name: "Expect with Multiple",
// 			expectToken: ExpectToken{
// 				Branches: []BranchToken{
// 					{Id: "a", Fn: func(*Parser) {}},
// 					{Id: "b", Fn: func(*Parser) {}},
// 				},
// 				Options: ParseOptions{Multiple: true},
// 			},
// 			expectedError: false,
// 		},
// 		{
// 			name: "Expect with Invert",
// 			expectToken: ExpectToken{
// 				Branches: []BranchToken{
// 					{Id: "d", Fn: func(*Parser) {}},
// 				},
// 				Options: ParseOptions{Invert: true},
// 			},
// 			expectedError: false,
// 		},
// 		{
// 			name: "Expect with Error",
// 			expectToken: ExpectToken{
// 				Branches: []BranchToken{
// 					{Id: "d", Fn: func(*Parser) {}},
// 				},
// 			},
// 			expectedError: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockScanner.index = 0 // Reset scanner for each test
// 			parser.Expect(tt.expectToken)
// 			if tt.expectedError && !parser.err {
// 				t.Errorf("Expected an error, but got none")
// 			}
// 			if !tt.expectedError && parser.err {
// 				t.Errorf("Did not expect an error, but got one")
// 			}
// 		})
// 	}
// }

// // TestPeek tests the Peek method of the Parser
// func TestPeek(t *testing.T) {
// 	mockScanner := &MockScanner{
// 		tokens: []Token{
// 			{ID: "a", Literal: "a", Line: 1, Position: 1},
// 			{ID: "b", Literal: "b", Line: 1, Position: 2},
// 			{ID: "c", Literal: "c", Line: 1, Position: 3},
// 		},
// 	}

// 	parser := &Parser{
// 		s: mockScanner,
// 		l: &MockLogger{},
// 	}

// 	called := false
// 	parser.Peek([]PeekToken{
// 		{
// 			IDs: []TokenType{"a", "b"},
// 			Fn:  func(*Parser) { called = true },
// 		},
// 	})

// 	if !called {
// 		t.Errorf("Peek function was not called when it should have been")
// 	}

// 	mockScanner.index = 0 // Reset scanner
// 	called = false
// 	parser.Peek([]PeekToken{
// 		{
// 			IDs: []TokenType{"a", "d"},
// 			Fn:  func(*Parser) { called = true },
// 		},
// 	})

// 	if called {
// 		t.Errorf("Peek function was called when it should not have been")
// 	}
// }

// // TestAddNode tests the AddNode method of the Parser
// func TestAddNode(t *testing.T) {
// 	mockAst := &mockAST{}
// 	parser := &Parser{
// 		ast: mockAst,
// 		l:   &MockLogger{},
// 	}

// 	parser.AddNode(NODE_ROOT)

// 	if len(mockAst.nodes) != 1 || mockAst.nodes[0] != NODE_ROOT {
// 		t.Errorf("AddNode did not correctly add the node to the AST")
// 	}
// }

// // TestAddTokens tests the AddTokens method of the Parser
// func TestAddTokens(t *testing.T) {
// 	mockAst := &mockAST{}
// 	parser := &Parser{
// 		ast:    mockAst,
// 		l:      &MockLogger{},
// 		tokens: []Token{{ID: "a", Literal: "a"}},
// 	}

// 	parser.AddTokens()

// 	if len(mockAst.tokens) != 1 || mockAst.tokens[0].ID != "a" {
// 		t.Errorf("AddTokens did not correctly add the tokens to the AST")
// 	}

// 	if len(parser.tokens) != 0 {
// 		t.Errorf("AddTokens did not clear the parser's token buffer")
// 	}
// }

// // TestRecover tests the Recover method of the Parser
// func TestRecover(t *testing.T) {
// 	parser := &Parser{
// 		err: true,
// 		l:   &MockLogger{},
// 	}

// 	called := false
// 	parser.Recover(func(*Parser) { called = true })

// 	if !called {
// 		t.Errorf("Recover function was not called when it should have been")
// 	}

// 	if parser.err {
// 		t.Errorf("Recover did not reset the error state")
// 	}
// }

// // Mock AST for testing
// type mockAST struct {
// 	nodes  []NodeType
// 	tokens []Token
// }

// func (m *mockAST) addNode(nt NodeType)      { m.nodes = append(m.nodes, nt) }
// func (m *mockAST) addToken(tokens []Token)  { m.tokens = append(m.tokens, tokens...) }
// func (m *mockAST) walkUp()                  {}
// func (m *mockAST) getRoot() *Node           { return nil }
// func (m *mockAST) getCurrent() *Node        { return nil }
// func (m *mockAST) getParent() *Node         { return nil }
// func (m *mockAST) getNode(id int) *Node     { return nil }
// func (m *mockAST) String() string           { return "" }
// func (m *mockAST) createNode(nt NodeType)   {}
// func (m *mockAST) addChild(child *Node)     {}
// func (m *mockAST) getChildren() []*Node     { return nil }
// func (m *mockAST) getTokens() []Token       { return nil }
// func (m *mockAST) setTokens(tokens []Token) {}

// type MockLogger struct{}

// func (m *MockLogger) log(msg string, indent indent) {}
