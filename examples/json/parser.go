package json

import (
	"github.com/dezlitz/dsl"
)

// NodeType represents the type of a node in the AST.
const (
	NODE_OBJECT dsl.NodeType = "OBJECT"
	NODE_ARRAY  dsl.NodeType = "ARRAY"
	NODE_MEMBER dsl.NodeType = "MEMBER"
	NODE_VALUE  dsl.NodeType = "VALUE"
)

func Parse(p *dsl.Parser) (dsl.AST, []dsl.Error) {
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_LBRACE, Fn: parseObject},
			{Id: TOKEN_LBRACKET, Fn: parseArray},
		},
	})

	return p.Exit()
}

func parseObject(p *dsl.Parser) {
	p.SkipToken()
	p.AddNode(NODE_OBJECT)

	// First member
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_STRING, Fn: func(p *dsl.Parser) {
				p.AddNode(NODE_MEMBER)
				p.AddTokens()
				p.Expect(dsl.ExpectToken{
					Branches: []dsl.BranchToken{
						{Id: TOKEN_COLON, Fn: nil},
					},
					Options: dsl.ParseOptions{Skip: true},
				})
				p.Call(parseValue)
			}},
		},
		Options: dsl.ParseOptions{Optional: true},
	})

	p.WalkUp()

	// Subsequent members
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_COMMA, Fn: parseKeyAndValue},
		},
		Options: dsl.ParseOptions{Multiple: true, Optional: true, Skip: true},
	})

	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_RBRACE, Fn: closeNode},
		},
	})
}

func parseKeyAndValue(p *dsl.Parser) {

	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_STRING, Fn: nil},
		},
	})

	p.AddNode(NODE_MEMBER)
	p.AddTokens()

	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_COLON, Fn: nil},
		},
		Options: dsl.ParseOptions{Skip: true},
	})

	p.Call(parseValue)
	p.WalkUp()
}

func parseValue(p *dsl.Parser) {

	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_STRING, Fn: addValue},
			{Id: TOKEN_NUMBER, Fn: addValue},
			{Id: TOKEN_TRUE, Fn: addValue},
			{Id: TOKEN_FALSE, Fn: addValue},
			{Id: TOKEN_NULL, Fn: addValue},
			{Id: TOKEN_LBRACE, Fn: parseObject},
			{Id: TOKEN_LBRACKET, Fn: parseArray},
		},
	})

}

func parseArray(p *dsl.Parser) {
	p.SkipToken()
	p.AddNode(NODE_ARRAY)

	// First value
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_STRING, Fn: addValue},
			{Id: TOKEN_NUMBER, Fn: addValue},
			{Id: TOKEN_TRUE, Fn: addValue},
			{Id: TOKEN_FALSE, Fn: addValue},
			{Id: TOKEN_NULL, Fn: addValue},
			{Id: TOKEN_LBRACE, Fn: parseObject},
			{Id: TOKEN_LBRACKET, Fn: parseArray},
		},
		Options: dsl.ParseOptions{Optional: true},
	})

	// Subsequent values
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_COMMA, Fn: parseValue},
		},
		Options: dsl.ParseOptions{Multiple: true, Optional: true, Skip: true},
	})

	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_RBRACKET, Fn: closeNode},
		},
	})

	p.WalkUp()

}

func addValue(p *dsl.Parser) {
	p.AddNode(NODE_VALUE)
	p.AddTokens()
	p.WalkUp()
}

func closeNode(p *dsl.Parser) {
	p.SkipToken()
	p.WalkUp()
}
