// Package dsl implements a set of helper and wrapper functions to allow an
// end user to create a parser for their Domain Specific Language. The user
// provides the Scan and Parse functions along with the input source. The
// package sets up and runs the parser, returning an AST.The output is an
// abstract syntax tree (AST) representing the Go source. The parser is invoked
// through one of the Parse* functions.
//
// This file contains the exported entry points for invoking the parser.
// The parser accepts either a bufio.Reader (Parse) or a string representing a
// file name (ParseFile). Along with the AST, any errors will be returned in a
// slice containing the line number and column. A log will also be produced to
// track any mistakes in the Scan or Parse function logic.

package dsl

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	logenb bool
)

// ParseFile sets up the bufio.Reader by accessing a file with the corresponding
// filename. If the file cannot be found an error will be returned at Error[0] and
// the AST will be nil.
//
// If the bufio.Reader was set up correctly it will be passed to the Parse function.
func ParseFile(pf ParseFunc, sf ScanFunc, ts TokenSet, ns NodeSet, inputfilename string) (AST, []Error) {
	inputfile, err := os.Open(inputfilename)
	if err != nil {
		return AST{}, []Error{{Code: ERROR_FILE_NOT_FOUND, Message: fmt.Sprintf("file '%v', not found", inputfilename)}}
	}
	r := bufio.NewReader(inputfile)
	return Parse(pf, sf, ts, ns, r)
}

// Parse sets up the parser, scanner and AST ready to accept input from
// the bufio.Reader and launches into the users entry parsing function.
//
// The function returns the AST, Errors and the Log. The user should check
// len(Errors) > 0 to determine if the input was correctly formed. The
// log is provided to diagnose errors in the parsing/scanning logic and can
// be ignored once the parse/scan functions have been proven correct.
func Parse(pf ParseFunc, sf ScanFunc, ts TokenSet, ns NodeSet, r *bufio.Reader) (AST, []Error) {
	s := newScanner(sf, r, nil)
	a := newAST(ns)
	p := newParser(pf, ts, s, a)
	return execute(p)
}

func execute(p *Parser) (AST, []Error) {
	pf := p.fn
	if logenb {
		p.log("Line 1: ", NO_PREFIX)
		p.log("Parsing: "+getFuncName(pf), NEWLINE)
	}
	ast, errors := pf(p)
	if logenb {
		p.log("Returning: "+getFuncName(pf), DECREMENT)
	}

	return ast, errors
}

func ParseFileAndLog(pf ParseFunc, sf ScanFunc, ts TokenSet, ns NodeSet, inputfilename string, logfilename string) (AST, []Error) {
	logfile, err := os.Create(logfilename)
	if err != nil {
		return AST{}, []Error{{Code: ERROR_COULD_NOT_CREATE_FILE, Message: fmt.Sprintf("could not create file '%v'", logfilename)}}
	}
	defer logfile.Close()
	//input, err := ioutil.ReadFile(inputfilename)
	inputfile, err := os.Open(inputfilename)
	if err != nil {
		return AST{}, []Error{{Code: ERROR_FILE_NOT_FOUND, Message: fmt.Sprintf("file '%v', not found", inputfilename)}}
	}
	r := bufio.NewReader(inputfile)
	return ParseAndLog(pf, sf, ts, ns, r, logfile)
}

func ParseAndLog(pf ParseFunc, sf ScanFunc, ts TokenSet, ns NodeSet, r *bufio.Reader, logfile io.Writer) (AST, []Error) {
	logenb = true
	l := log.New(logfile, "", 0)
	s := newScanner(sf, r, l)
	a := newAST(ns)
	p := newParser(pf, ts, s, a)
	return execute(p)
}
