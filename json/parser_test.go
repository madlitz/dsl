package json

import (
	"github.com/madlitz/go-dsl"
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
	p.Peek([]dsl.PeekToken{
		{IDs: []dsl.TokenType{TOKEN_STRING}, Fn: parseMember},
	})

	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_COMMA, Fn: parseMember},
		},
		Options: dsl.ParseOptions{Multiple: true, Optional: true, Skip: true},
	})

	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_RBRACE, Fn: closeNode},
		},
	})
}

func parseMember(p *dsl.Parser) {

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
	})
	p.SkipToken() // Skip the colon
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_STRING, Fn: parseValue},
			{Id: TOKEN_NUMBER, Fn: parseValue},
			{Id: TOKEN_TRUE, Fn: parseValue},
			{Id: TOKEN_FALSE, Fn: parseValue},
			{Id: TOKEN_NULL, Fn: parseValue},
			{Id: TOKEN_LBRACE, Fn: parseObject},
			{Id: TOKEN_LBRACKET, Fn: parseArray},
		},
	})
	p.WalkUp()
}

func parseArray(p *dsl.Parser) {
	p.SkipToken()
	p.AddNode(NODE_ARRAY)

	p.Peek([]dsl.PeekToken{
		{IDs: []dsl.TokenType{TOKEN_STRING}, Fn: parseArrayValue},
		{IDs: []dsl.TokenType{TOKEN_NUMBER}, Fn: parseArrayValue},
		{IDs: []dsl.TokenType{TOKEN_TRUE}, Fn: parseArrayValue},
		{IDs: []dsl.TokenType{TOKEN_FALSE}, Fn: parseArrayValue},
		{IDs: []dsl.TokenType{TOKEN_NULL}, Fn: parseArrayValue},
		{IDs: []dsl.TokenType{TOKEN_LBRACE}, Fn: parseArrayValue},
		{IDs: []dsl.TokenType{TOKEN_LBRACKET}, Fn: parseArrayValue},
	})

	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_RBRACKET, Fn: closeNode},
		},
	})

	p.WalkUp()

}

func parseArrayValue(p *dsl.Parser) {

	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_STRING, Fn: parseValue},
			{Id: TOKEN_NUMBER, Fn: parseValue},
			{Id: TOKEN_TRUE, Fn: parseValue},
			{Id: TOKEN_FALSE, Fn: parseValue},
			{Id: TOKEN_NULL, Fn: parseValue},
			{Id: TOKEN_LBRACE, Fn: parseObject},
			{Id: TOKEN_LBRACKET, Fn: parseArray},
			{Id: TOKEN_RBRACKET, Fn: closeNode},
		},
	})

	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{Id: TOKEN_COMMA, Fn: parseArrayValue},
		},
		Options: dsl.ParseOptions{Multiple: true, Optional: true, Skip: true},
	})
}

func parseValue(p *dsl.Parser) {
	p.AddNode(NODE_VALUE)
	p.AddTokens()
	p.WalkUp()
}

func closeNode(p *dsl.Parser) {
	p.SkipToken()
	p.WalkUp()
}
