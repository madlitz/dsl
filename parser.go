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
	fn   ParseFunc
	s    scanner
	ast  AST
	l    logger
	line string
	buf  struct {
		tokens []Token
		num    int
	} // Holds unread tokens so we don't have to make repeat calls to the Scanner
	tokens     []Token // Holds all tokens consumed until they are moved to the AST
	peekBuffer []Token // Holds all tokens peeked until they are consumed
	errors     []Error
	eof        bool
	err        bool
	loopCheck  struct {
		count int
		line  int
		pos   int
	}
}

// TokenType is a string that represents the type of token found in the source.
// String is used to represent the token type as it is easier to print in logs.
type TokenType string

// Token is used to identify the token type and value after being added to the AST.
// The extra information is used when displaying errors but could also be useful
// to the user for things like syntax highlighting and debugging if they were to
// implement it.
type Token struct {
	ID       TokenType
	Literal  string
	Line     int // Line is the line of the source text the Token was found.
	Position int // Position is the position (or column) the Token was found.
}

const (
	TOKEN_UNKNOWN TokenType = "UNKNOWN"
	TOKEN_ERROR   TokenType = "ERROR"
	TOKEN_EOF     TokenType = "EOF"
)

// newParser returns an instance of a Parser
func newParser(pf ParseFunc, s *Scanner, ast AST, l logger) *Parser {
	return &Parser{
		fn:  pf,
		s:   s,
		ast: ast,
		l:   l,
	}
}

// The user implements the ParseFunc and passes it to the Parser in dsl.Parse()
type ParseFunc func(*Parser) (AST, []Error)

// Used as an input to Parser.Expect()
type ExpectToken struct {
	Branches []BranchToken
	Options  ParseOptions
}

type ExpectNotToken struct {
	Tokens  []TokenType
	Fn      func(*Parser)
	Options ParseOptions
}

type BranchToken struct {
	Id TokenType
	Fn func(*Parser)
}

// If the Optional option is false and a match is not found, an error is returned to the
// parser.
//
// If the Multiple option is set to true the parser continues to read, consume tokens and
// take branches until a token is read that is not matched by any of the branches. If the
// Multiple option is set to false, only the first branch to be matched is taken and
// consumed.
//
// If the Peek option is set to true the parser will not consume the token that is matched
// and will push the token back onto the buffer if no subseqent match is found.
//
// If the the Skip option is set to true the parser will take the branch if a match is
// found but will not consume the token.
type ParseOptions struct {
	Optional bool
	Multiple bool
	Peek     bool
	Skip     bool
}

//-----------------------------------------------------------------------------------

func (p *Parser) Expect(expect ExpectToken) {
	var found1orMore bool
	var tok Token
	var err *Error

	//If we have previously found an error but have not yet recovered with p.Recover, skip any call to p.Expect.
	p.log(fmt.Sprintf("Expect Token %v: %v ", getParseOptions(expect.Options), branchTokensToStrings(expect.Branches)), prefixNewline)
	if p.err {
		p.log("Skipping Expect as error already found.", prefixNewline)
		return
	}
	for {
		var branchFn func(*Parser)
		found := false

		tok, err = p.scan()
		if err != nil {
			return
		}
		if tok.ID == TOKEN_EOF {
			p.eof = true
		}
		for _, branch := range expect.Branches {
			if tok.ID == branch.Id {
				found = true
				branchFn = branch.Fn
				break
			}
		}

		// If no match was found, unread the token and break out of the loop
		if !found {
			p.unscan()
			break
		}

		// Match found

		if expect.Options.Peek {
			// If we are peeking, remember each token read
			p.peekBuffer = append(p.peekBuffer, tok)
		} else {
			// If we are not peeking, consume the token
			if len(p.peekBuffer) > 0 {
				// If we have peeked but now no longer peeking, consume the peeked runes
				p.consumePeeked()
			}
			p.consume(tok, expect.Options.Skip)
		}

		p.logMatch(tok)
		found1orMore = true // Set to true only after logging the first match
		p.parseFn(branchFn)

		if !expect.Options.Multiple || p.eof || p.err {
			// If Multiple is false break out of the loop
			break
		}
	}

	// If we have finished peeking and there are tokens to unread, unread them
	if expect.Options.Peek && len(p.peekBuffer) > 0 {
		p.unreadPeeked()
	}

	if !found1orMore && !expect.Options.Optional {
		p.newError(ErrorTokenExpectedNotFound, fmt.Errorf("found [%v], expected any of %v", tok.ID, branchTokensToStrings(expect.Branches)), p.tokToErrLine(tok))
	}
}

func (p *Parser) ExpectNot(expect ExpectNotToken) {
	var found bool
	var found1orMoreNot bool
	var tok Token
	var err *Error

	p.log(fmt.Sprintf("Expect Not Token %v: %v ", getParseOptions(expect.Options), tokensToStrings(expect.Tokens)), prefixNewline)
	if p.err {
		p.log("Skipping Expect Not as error already found.", prefixNewline)
		return
	}

	for {
		found = false
		tok, err = p.scan()
		if err != nil {
			return
		}
		if tok.ID == TOKEN_EOF {
			p.eof = true
		}

		for _, token := range expect.Tokens {
			if tok.ID == token {
				found = true
				break
			}
		}

		// If a match was found, unread the token and break out of the loop
		if found {
			p.unscan()
			break
		}

		// Match not found, which is what we are expecting

		if expect.Options.Peek {
			// If we are peeking, remember each token read
			p.peekBuffer = append(p.peekBuffer, tok)
		} else {
			// If we are not peeking, consume the token
			if len(p.peekBuffer) > 0 {
				// If we have peeked but now no longer peeking, consume the peeked tokens
				p.consumePeeked()
			}
			p.consume(tok, expect.Options.Skip)
		}

		p.logMatch(tok)
		found1orMoreNot = true
		p.parseFn(expect.Fn)

		if !expect.Options.Multiple || p.eof || p.err {
			// If Multiple is false break out of the loop
			break
		}
	}

	// If we have finished peeking and there are tokens to unread, unread them
	if expect.Options.Peek && len(p.peekBuffer) > 0 {
		p.unreadPeeked()
	}

	if !found1orMoreNot && !expect.Options.Optional {
		p.newError(ErrorTokenExpectedNotFound, fmt.Errorf("found [%v], expected any except %v", tok.ID, tokensToStrings(expect.Tokens)), p.tokToErrLine(tok))
	}
}

func (p *Parser) parseFn(fn func(*Parser)) {
	if ok, tok := p.checkForInfiniteLoop(); ok {
		p.newError(ErrorInfiniteLoopDetected, fmt.Errorf("infinite loop detected: %v", getFuncName(fn)), p.tokToErrLine(tok))
		return
	}
	if fn != nil && !p.eof {
		p.log("Parsing: "+getFuncName(fn), prefixIncrement)
		fn(p)
		p.log("Returning: "+getFuncName(fn), prefixDecrement)
	}
}

func (p *Parser) consume(tok Token, skip bool) {
	if !skip {
		p.tokens = append(p.tokens, tok)
	}
}

func (p *Parser) logMatch(tok Token) {
	p.log("Found: ", prefixNone)
	p.log(string(tok.ID), prefixNone)
}

func (p *Parser) AddNode(nt NodeType) {
	p.log("AST Add Node: "+string(nt), prefixNewline)
	p.ast.addNode(nt)
}

func (p *Parser) AddTokens() {
	p.log("AST Add Tokens: ", prefixNewline)
	if len(p.tokens) > 0 {
		for _, token := range p.tokens {
			p.log(string(token.ID)+" - ", prefixNone)
			for _, rn := range token.Literal {
				p.log(sanitize(string(rn), false), prefixNone)
			}
			p.log(", ", prefixNone)
		}
		p.ast.addToken(p.tokens)
		p.tokens = nil
	} else {
		p.log("Warning: No Tokens to Add", prefixError)
	}
}

func (p *Parser) SkipToken() {
	p.log("AST Skip Token: ", prefixNewline)
	if len(p.tokens) > 0 {
		token := p.tokens[len(p.tokens)-1]
		p.tokens = p.tokens[:len(p.tokens)-1]
		p.log(string(token.ID)+" - ", prefixNone)
		p.log(sanitize(token.Literal, true)+", ", prefixNone)
	} else {
		p.log("Warning: No Tokens to Skip", prefixError)
	}
}

func (p *Parser) GetToken() Token {
	if len(p.tokens) == 0 {
		p.log("Error: No tokens to get.", prefixError)
		return Token{TOKEN_ERROR, "ERROR", 0, 0}
	}
	token := p.tokens[len(p.tokens)-1]
	p.log("Get Last Token: ", prefixNewline)
	p.log(sanitize(token.Literal, true), prefixNone)
	return token
}

func (p *Parser) WalkUp() {
	p.log("AST Walk Up", prefixNewline)
	p.ast.walkUp()
}

func (p *Parser) Call(fn func(*Parser)) {
	if fn != nil && !p.eof {
		p.log("Calling: "+getFuncName(fn), prefixIncrement)
		fn(p)
		p.log("Returning: "+getFuncName(fn), prefixDecrement)
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
	tok, line, err := p.s.scan()
	if err != nil {
		p.err = true
		p.errors = append(p.errors, *err)
	}
	p.line = line

	// Save it to the buffer in case we unscan later.
	p.buf.tokens = append(p.buf.tokens, tok)

	return
}

const MaxConsecutivePeeks = 10

func (p *Parser) checkForInfiniteLoop() (bool, Token) {
	tok := p.buf.tokens[len(p.buf.tokens)-1]

	if p.loopCheck.line == tok.Line && p.loopCheck.pos == tok.Position {
		p.loopCheck.count++
		if p.loopCheck.count > MaxConsecutivePeeks {
			return true, tok
		}
	} else {
		p.loopCheck.count = 0
	}
	p.loopCheck.line = tok.Line
	p.loopCheck.pos = tok.Position

	return false, Token{}

}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() {
	if p.buf.num < len(p.buf.tokens) {
		p.buf.num++
	}
}

// unreadPeeked is used to unread the peeked tokens
func (p *Parser) unreadPeeked() {
	for i := len(p.peekBuffer) - 1; i >= 0; i-- {
		p.unscan()
	}
	p.peekBuffer = nil
}

// consumePeeked is used to consume the peeked tokens
func (p *Parser) consumePeeked() {
	for _, tok := range p.peekBuffer {
		p.consume(tok, false)
	}
	p.peekBuffer = nil
}

type errorLine struct {
	line      string
	startLine int
	startPos  int
	endLine   int
	endPos    int
}

func (p *Parser) tokToErrLine(tok Token) errorLine {
	endPos := tok.Position
	if tok.ID == TOKEN_EOF {
		endPos++
	}
	return errorLine{
		line:      p.line,
		startLine: tok.Line,
		startPos:  tok.Position,
		endLine:   tok.Line,
		endPos:    endPos + len(tok.Literal) - 1,
	}
}

func (p *Parser) newError(code ErrorCode, errMsg error, el errorLine) {
	p.err = true
	p.errors = append(p.errors, Error{
		Code:          code,
		Message:       errMsg.Error(),
		LineString:    el.line,
		StartLine:     el.startLine,
		StartPosition: el.startPos,
		EndLine:       el.endLine,
		EndPosition:   el.endPos,
	})
	p.log(errMsg.Error(), prefixError)

}

func (p *Parser) Recover(fn func(*Parser)) {
	if !p.err {
		return
	}

	if fn != nil && !p.eof {
		p.log("Recovering: "+getFuncName(fn), prefixIncrement)
		p.err = false
		fn(p)
		p.log("Returning: "+getFuncName(fn), prefixDecrement)
	}
}

// -------------------------------- Parser Helper Functions---------------------------------------

// Used to log which options were used during a branch function call
func getParseOptions(opts ParseOptions) string {
	var buf bytes.Buffer
	buf.WriteRune('(')
	if opts.Optional {
		buf.WriteString("Optional ")
	}
	if opts.Multiple {
		buf.WriteString("Multiple ")
	}
	if opts.Skip {
		buf.WriteString("Skip ")
	}
	if opts.Peek {
		buf.WriteString("Peek ")
	}
	buf.WriteRune(')')
	return buf.String()
}

// Used to print tokens passed to Peek() to an Error or to the log
func tokensToStrings(tokens []TokenType) (literals []string) {
	for _, token := range tokens {
		literals = append(literals, fmt.Sprintf("%v, ", token))
	}
	return
}

// Used to print tokens passed to Expect() to an Error or to the log
func branchTokensToStrings(branches []BranchToken) (literals []string) {
	for _, branch := range branches {
		literals = append(literals, string(branch.Id))
	}
	return
}

// Send log function to the scanner as the Scanner contains the log
func (p *Parser) log(msg string, indent indent) {
	p.l.log(msg, indent)
}
