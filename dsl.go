// Package dsl implements a set of helper and wrapper functions to allow an
// end user to create a parser for their Domain Specific Language. The user
// provides the Scan and Parse functions along with the input source. The
// package sets up and runs the parser, returning an AST. The output is an
// abstract syntax tree (AST) representing the DSL source.
//
// This file contains the exported entry points for invoking the parser.
// Along with the AST, any errors will be returned in a slice containing the
// line number and column. A log will also be produced to track any mistakes
// in the Scan or Parse function logic including infinite loops.

package dsl

import (
	"bufio"
	"io"
	"log"
)

// ParseOption is a function type that modifies ParseConfig
type ParseOption func(*ParseConfig)

// ParseConfig holds the configuration for parsing
type ParseConfig struct {
	LogWriter io.Writer
	// Add other configuration options here as needed
}

// WithLogger returns a ParseOption that sets the log writer
func WithLogger(w io.Writer) ParseOption {
	return func(c *ParseConfig) {
		c.LogWriter = w
	}
}

// Parse sets up the parser, scanner and AST ready to accept input from
// the bufio.Reader and launches into the users entry parsing function.
//
// The function returns the AST, Errors and the Log. The user should check
// len(Errors) > 0 to determine if the input was correctly formed. The
// log is provided to diagnose errors in the parsing/scanning logic and can
// be ignored once the parse/scan functions have been proven correct.
func Parse(pf ParseFunc, sf ScanFunc, r *bufio.Reader, opts ...ParseOption) (AST, []Error) {

	config := &ParseConfig{}

	// Apply the options
	for _, opt := range opts {
		opt(config)
	}

	var logger logger
	logger = &dslNoLogger{}
	if config.LogWriter != nil {
		logger = &dslLogger{
			logger: log.New(config.LogWriter, "", 0),
		}
	}

	s := newScanner(sf, r, logger)
	a := newAST()
	p := newParser(pf, s, a, logger)
	return execute(p)
}

func execute(p *Parser) (AST, []Error) {
	pf := p.fn
	p.log("Line 1: ", prefixNone)
	p.log("Parsing: "+getFuncName(pf), prefixNewline)
	ast, errors := pf(p)
	p.log("Returning: "+getFuncName(pf), prefixDecrement)
	return ast, errors
}
