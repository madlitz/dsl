package mydsl

import (
	"github.com/madlitz/dsl"
)

// NodeType represents the type of a node in the AST.
const (
	NODE_ASSIGNMENT dsl.NodeType = "ASSIGNMENT"
	NODE_CALL       dsl.NodeType = "CALL"
	NODE_EXPRESSION dsl.NodeType = "EXPRESSION"
	NODE_TERMINAL   dsl.NodeType = "TERMINAL"
	NODE_COMMENT    dsl.NodeType = "COMMENT"
)

var recover bool

func Parse(p *dsl.Parser) (dsl.AST, []dsl.Error) {
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_WS, Fn: skipWhitespace},
			{Id: TOKEN_VARIABLE, Fn: assignmentOrCall},
			{Id: TOKEN_EOF, Fn: nil},
		},
		Options: dsl.ParseOptions{Multiple: true},
	})

	return p.Exit()
}

func skipWhitespace(p *dsl.Parser) {
	p.SkipToken()
}

// parse -> assignmentOrCall
func assignmentOrCall(p *dsl.Parser) {
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_WS, Fn: skipWhitespace},
		},
		Options: dsl.ParseOptions{Optional: true},
	})
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_ASSIGN, Fn: assignment},
			{Id: TOKEN_OPEN_PAREN, Fn: call},
		},
	})
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_WS, Fn: skipWhitespace},
		},
		Options: dsl.ParseOptions{Optional: true},
	})
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_COMMENT, Fn: addcomment},
		},
		Options: dsl.ParseOptions{Optional: true},
	})
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_NL, Fn: skipWhitespace},
			{Id: TOKEN_EOF, Fn: nil},
		},
	})
}

// parse -> assignmentOrCall -> assignment
func assignment(p *dsl.Parser) {
	p.SkipToken()
	p.AddNode(NODE_ASSIGNMENT)
	p.AddTokens()
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_WS, Fn: skipWhitespace},
		},
		Options: dsl.ParseOptions{Optional: true},
	})
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_VARIABLE, Fn: operator},
			{Id: TOKEN_LITERAL, Fn: operator},
			{Id: TOKEN_OPEN_PAREN, Fn: parenExpression},
		},
	})
}

// parse -> assignmentOrCall -> call
func call(p *dsl.Parser) {
	p.SkipToken()
	p.AddNode(NODE_CALL)
	p.AddTokens()
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_WS, Fn: skipWhitespace},
		},
		Options: dsl.ParseOptions{Optional: true},
	})
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_VARIABLE, Fn: operator},
			{Id: TOKEN_LITERAL, Fn: operator},
			{Id: TOKEN_OPEN_PAREN, Fn: parenExpression},
		},
	})
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_WS, Fn: skipWhitespace},
		},
		Options: dsl.ParseOptions{Optional: true},
	})
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_CLOSE_PAREN, Fn: closecall},
		},
	})
}

// parse -> assignmentOrCall -> call -> [expression] -> closecall
func closecall(p *dsl.Parser) {
	p.SkipToken()
	p.WalkUp()
}

// parse -> assignmentOrCall -> assignment -> [operator, expression]
// parse -> assignmentOrCall -> call -> [operator, expression]
func expression(p *dsl.Parser) {
	p.AddNode(NODE_EXPRESSION)
	p.AddTokens()
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_WS, Fn: skipWhitespace},
		},
		Options: dsl.ParseOptions{Optional: true},
	})
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_VARIABLE, Fn: operator},
			{Id: TOKEN_LITERAL, Fn: operator},
			{Id: TOKEN_OPEN_PAREN, Fn: parenExpression},
			{Id: TOKEN_CLOSE_PAREN, Fn: operator},
		},
	})
}

// parse -> assignmentOrCall -> assignment -> [expression, operator]
// parse -> assignmentOrCall -> call -> [expression, operator]
func operator(p *dsl.Parser) {
	p.AddNode(NODE_TERMINAL)
	p.AddTokens()
	p.WalkUp()
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_WS, Fn: skipWhitespace},
		},
		Options: dsl.ParseOptions{Optional: true},
	})
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_PLUS, Fn: expression},
			{Id: TOKEN_MINUS, Fn: expression},
			{Id: TOKEN_DIVIDE, Fn: expression},
			{Id: TOKEN_MULTIPLY, Fn: expression},
			{Id: TOKEN_OPEN_PAREN, Fn: parenExpression},
		},
		Options: dsl.ParseOptions{Optional: true},
	})
	p.WalkUp()
}

// parse -> assignmentOrCall -> assignment -> [expression, operator] -> paren_expression
// parse -> assignmentOrCall -> assignment -> paren_expression
// parse -> assignmentOrCall -> call -> [expression, operator] -> paren_expression
// parse -> assignmentOrCall -> call -> paren_expression
func parenExpression(p *dsl.Parser) {
	p.Peek([]dsl.PeekToken{
		{IDs: []dsl.TokenType{}, Fn: expression},
	})
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_WS, Fn: skipWhitespace},
		},
		Options: dsl.ParseOptions{Optional: true},
	})
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_CLOSE_PAREN, Fn: operator},
		},
	})
	p.Recover(skipUntilLineBreak)
}

// parse -> assignmentOrCall -> [expression] -> addcomment
func addcomment(p *dsl.Parser) {
	p.AddNode(NODE_COMMENT)
	p.AddTokens()
	p.WalkUp()
}

func skipUntilLineBreak(p *dsl.Parser) {
	recover = true
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: dsl.TOKEN_UNKNOWN, Fn: nil},
		},
	})
	recover = false
}
