// Copyright (c) 2024 Dez Little <deslittle@gmail.com>
// All rights reserved. Use of this source code is governed by a LGPL v3
// license that can be found in the LICENSE file.

// parser.go implements the Parser for any DSL text. On construction it
// creates a new Scanner and AST and keeps a reference to them, then
// calls the users parser function.
//
// When the user calls functions such as p.Expect() in their user parse
// function, the Parser makes calls to the Scanner to ask for a token
// and to the AST to create nodes and store the tokens.
package dsl

import (
	"bytes"
	"fmt"
)

// The Parser type holds a reference to the user parse func, the Scanner,
// AST, Errors to return and other state variables.
type Parser struct {
	fn  ParseFunc
	ts  TokenSet
	s   *Scanner
	ast AST      // AST is just 2 pointers so no need to use another pointer
	buf struct { // buf holds unread tokens so we don't have to make repeat calls to the Scanner
		tokens []Token
		num    int
	}
	tokens []Token // Holds all tokens consumed until they are moved to the AST
	errors []Error
	eof    bool
	err    bool
}

// newParser returns an instance of a Parser
func newParser(pf ParseFunc, ts TokenSet, s *Scanner, a AST) *Parser {
	return &Parser{fn: pf, ts: ts, s: s, ast: a}
}

// The user implements the ParseFunc and passes it to the Parser in dsl.Parse()
type ParseFunc func(*Parser) (AST, []Error)

// Used as an input to Parser.Expect()
type ExpectToken struct {
	Branches []BranchToken
	Options  ParseOptions
}

// ID is an interface implemented by to user so they can use their own token ID's
type BranchToken struct {
	Id TokenType
	Fn func(*Parser)
}

type PeekToken struct {
	IDs []TokenType
	Fn  func(*Parser)
}

// If the Optional option is false and a match is not found, an error is returned to the
// parser.
//
// If the Multiple option is set to true the parser continues to read, consume tokens and
// take branches until a token is read that is not matched by any of the branches. If the
// Multiple option is set to false, only the first branch to be matched is taken and
// consumed.
//
// If the Invert option is set to true the parser consumes the token and takes a branch if
// it doesn't match any of the branches.
//
// If the the Skip option is set to true the parser will take the branch if a match is
// found but will not consume the token.
//
// If ParseOptions is omitted when creating {}ExpectToken, all options will be set to
// false.
type ParseOptions struct {
	Optional bool
	Multiple bool
	Invert   bool
	Skip     bool
}

//-----------------------------------------------------------------------------------

// Used as an option to log()
type indent int

const (
	NO_PREFIX indent = iota // 0
	NEWLINE
	INCREMENT
	DECREMENT
	STARTLINE
	ERROR
)

func (p *Parser) Expect(expect ExpectToken) {
	var found bool
	var tok Token
	var err *Error

	//If we have previously found an error but have not yet recovered with p.Recover, skip any call to p.Expect.
	if logenb {
		p.log(fmt.Sprintf("Expect Token %v: %v ", getParseOptions(expect.Options), branchTokensToStrings(expect.Branches)), NEWLINE)
	}
	if p.err {
		if logenb {
			p.log("Skipping Expect as error already found.", NEWLINE)
		}
		return
	}
	for {
		found = false
		tok, err = p.scan()
		if p.ts[tok.ID] == p.ts["EOF"] {
			p.eof = true
		}
		if err != nil {
			p.errors = append(p.errors, *err)
			p.err = true
			return
		}
		for _, branch := range expect.Branches {
			if p.ts[branch.Id] == 0 {
				if logenb {
					p.log(fmt.Sprintf("Error: Expected token [%v], not found in Token Set.", branch.Id), ERROR)
				}
				p.newError(ERROR_EXPECTED_TOKEN_NOT_IN_TOKENSET, fmt.Errorf("token %v, not found in Token Set", branch.Id))
				break
			}
			if p.ts[tok.ID] == 0 {
				if logenb {
					p.log(fmt.Sprintf("Error: Scanned token [%v], not found in Token Set.", tok.ID), ERROR)
				}
				p.newError(ERROR_SCANNED_TOKEN_NOT_IN_TOKENSET, fmt.Errorf("scanned token [%v], not found in Token Set", tok.ID))
				break
			}
			if p.ts[tok.ID] == p.ts[branch.Id] && !expect.Options.Invert {
				if !expect.Options.Invert {
					p.consume(tok, expect.Options.Skip)
				}
				found = true
				p.callFn(branch.Fn)
				break
			}
		}
		if (!expect.Options.Invert && !found) || (expect.Options.Invert && found) {
			p.unscan()
			break
		}
		if expect.Options.Invert && !found {
			p.consume(tok, expect.Options.Skip)
		}
		if !expect.Options.Multiple || p.eof || p.err {
			break
		}
	}
	if !found && !expect.Options.Optional && !expect.Options.Invert {
		p.newError(ERROR_TOKEN_EXPECTED_NOT_FOUND, fmt.Errorf("found [%v], expected any of %v", tok.ID, branchTokensToStrings(expect.Branches)))
	} else if !found && !expect.Options.Optional && expect.Options.Invert {
		p.newError(ERROR_TOKEN_EXPECTED_NOT_FOUND, fmt.Errorf("found [%v], expected any except %v", tok.ID, branchTokensToStrings(expect.Branches)))
	}
}

func (p *Parser) callFn(fn func(*Parser)) {
	if fn != nil && !p.eof {
		if logenb {
			p.log("Parsing: "+getFuncName(fn), INCREMENT)
		}
		fn(p)
		if logenb {
			p.log("Returning: "+getFuncName(fn), DECREMENT)
		}
	}
}

func (p *Parser) consume(tok Token, skip bool) {
	if logenb {
		p.log("Found: ", NEWLINE)
		p.log(tok.ID.String(), NO_PREFIX)
	}
	if !skip {
		p.tokens = append(p.tokens, tok)
	}
}

func (p *Parser) AddNode(nt NodeType) {
	if logenb {
		p.log("AST Add Node: "+nt.String(), NEWLINE)
	}
	if p.ast.ns[nt] == 0 {
		if logenb {
			p.log(fmt.Sprintf("Error: Node to add [%v], not found in NodeSet.", nt), ERROR)
		}
		p.newError(ERROR_NODE_NOT_IN_NODESET, fmt.Errorf("node to add [%v], not found in NodeSet", nt))
	}
	p.ast.addNode(nt)
}

func (p *Parser) AddTokens() {
	if logenb {
		p.log("AST Add Tokens: ", NEWLINE)
	}
	if len(p.tokens) > 0 {
		if logenb {
			for _, token := range p.tokens {
				p.log(token.ID.String()+" - ", NO_PREFIX)
				for _, rn := range token.Literal {
					p.log(sanitize(string(rn), false), NO_PREFIX)
				}
				p.log(", ", NO_PREFIX)
			}
		}
		p.ast.addToken(p.tokens)
		p.tokens = nil
	} else {
		if logenb {
			p.log("Warning: No Tokens to Add", ERROR)
		}
	}
}

func (p *Parser) SkipToken() {
	if logenb {
		p.log("AST Skip Token: ", NEWLINE)
	}
	if len(p.tokens) > 0 {
		token := p.tokens[len(p.tokens)-1]
		p.tokens = p.tokens[:len(p.tokens)-1]
		if logenb {
			p.log(token.ID.String()+" - ", NO_PREFIX)
			p.log(sanitize(token.Literal, true)+", ", NO_PREFIX)
		}
	} else {
		if logenb {
			p.log("Warning: No Tokens to Skip", ERROR)
		}
	}
}

func (p *Parser) GetToken() Token {
	if len(p.tokens) == 0 {
		p.log("Error: No tokens to get.", ERROR)
		return Token{"ERROR", "ERROR", p.s.curLine, p.s.curPos}
	}
	token := p.tokens[len(p.tokens)-1]
	if logenb {
		p.log("Get Last Token: ", NEWLINE)
		p.log(sanitize(token.Literal, true), NO_PREFIX)
	}
	return token
}

func (p *Parser) WalkUp() {
	if logenb {
		p.log("AST Walk Up", NEWLINE)
	}
	p.ast.walkUp()
}

func (p *Parser) Call(fn func(*Parser)) {
	if logenb {
		p.log("Call", NEWLINE)
	}
	p.callFn(fn)
}

func (p *Parser) Peek(branches []PeekToken) {
	if logenb {
		p.log(fmt.Sprintf("Peek: %v ", peekTokensToStrings(branches)), NEWLINE)
	}
	for _, branch := range branches {
		tokensLen := len(branch.IDs)
		bufLen := 0
		for i := 0; i < tokensLen; i++ {
			bufLen++
			tok, err := p.scan()
			if err != nil {
				p.errors = append(p.errors, *err)
			}
			if p.ts[tok.ID] != p.ts[branch.IDs[i]] {
				break
			}
		}
		for i := 0; i < bufLen; i++ {
			p.unscan()
		}
		if tokensLen == 0 || tokensLen == bufLen {
			p.callFn(branch.Fn)
			break
		}
	}
}

func (p *Parser) Exit() (AST, []Error) {
	return p.ast, p.errors
}

func (p *Parser) scan() (tok Token, err *Error) {
	// If we have a token on the buffer, then return it.
	if p.buf.num > 0 {
		tok = p.buf.tokens[len(p.buf.tokens)-p.buf.num]
		p.buf.num--
		return
	}

	// Otherwise read the next token from the scanner.
	tok, err = p.s.scan()
	// Save it to the buffer in case we unscan later.
	p.buf.tokens = append(p.buf.tokens, tok)

	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() {
	if p.buf.num < len(p.buf.tokens) {
		p.buf.num++
	}
}

func (p *Parser) newError(code ErrorCode, errMsg error) {
	p.err = true
	err := p.s.newError(code, errMsg)
	if err != nil {
		p.errors = append(p.errors, *err)
	}
}

func (p *Parser) Recover(Fn func(*Parser)) {
	if !p.err {
		return
	}
	if logenb {
		p.log("Recovering....", NEWLINE)
	}
	p.err = false
	p.callFn(Fn)
}

// -------------------------------- Parser Helper Functions---------------------------------------

// Used to log which options were used during a branch function call
func getParseOptions(opts ParseOptions) string {
	var buf bytes.Buffer
	buf.WriteRune('(')
	if opts.Invert {
		buf.WriteString("Invert ")
	}
	if opts.Optional {
		buf.WriteString("Optional ")
	}
	if opts.Multiple {
		buf.WriteString("Multiple ")
	}
	if opts.Skip {
		buf.WriteString("Skip ")
	}
	buf.WriteRune(')')
	return buf.String()
}

// Used to print tokens passed to Peek() to an Error or to the log
func peekTokensToStrings(branches []PeekToken) (literals []string) {
	for _, branch := range branches {
		literals = append(literals, fmt.Sprintf("%v, ", branch.IDs))
	}
	return
}

// Used to print tokens passed to Expect() to an Error or to the log
func branchTokensToStrings(branches []BranchToken) (literals []string) {
	for _, branch := range branches {
		literals = append(literals, branch.Id.String())
	}
	return
}

// Send log function to the scanner as the Scanner contains the log
func (p *Parser) log(msg string, indent indent) {
	p.s.log(msg, indent)
}
