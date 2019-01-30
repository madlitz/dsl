// Copyright (c) 2015 Des Little <deslittle@gmail.com>
// All rights reserved. Use of this source code is governed by a LGPL v3
// license that can be found in the LICENSE file.

// scanner.go implements a Scanner for any DSL source text.
// When a Scanner is created it takes a bufio.Reader as the source of
// the text, scans a number of characters (runes) as defined by the
// user ScanFunc and returns a token to the Parser.
//
package dsl

import (
	"bufio"
	"bytes"
	"fmt"
    "log"
)

type dslLogger struct{
    indent     int
    log        *log.Logger
    buf        []interface{}
}
// The Scanner contains a reference to the user scan function, the
// user input buffer, various state variables and the parser log.
//
type Scanner struct {
	fn  ScanFunc
	r   *bufio.Reader
	buf struct {
		runes     []rune
		unread    int
	}
	curLineBuffer bytes.Buffer
	curLine       int
	curPos        int
	options       ScanOptions
	expRunes      []rune
	tok           Token
	error         *Error
    logger        *dslLogger
	eof           bool
}

// NewScanner returns a new instance of Scanner.
func newScanner(sf ScanFunc, r *bufio.Reader, l *log.Logger) *Scanner {
    s := &Scanner{
            fn: sf,
            r: r, 
            curLine: 1,
            curPos: 1,
        }
    if l != nil{
        s.logger = &dslLogger{log: l}
    }
	return s
}

type ScanFunc func(*Scanner) Token

// If the Optional option is false and a match is not found, an error is returned to the
// parser.
//
// If the Multiple option is set to true the scanner continues to read, consume runes and
// take branches until a rune is read that is not matched by any of the branches or branch
// ranges. If the Multiple option is set to false, only the first branch (or branch range)
// to be matched is taken and consumed.
//
// If the Invert option is set to true the scanner consumes the rune and takes a branch if
// it doesn't match any of the branch or branch ranges.
//
// If the the Skip option is set to true the scanner will take the branch if a match is
// found but will not consume the rune.
//
// If ScanOptions is omitted when creating ExpectRune{}, all options will be set to
// false.
//
type ScanOptions struct {
	Optional bool
	Multiple bool
	Invert   bool
	Skip     bool
    Error    func(*Scanner)
}

type ExpectRune struct {
	Branches     []Branch
	BranchRanges []BranchRange
	Options      ScanOptions
}

type Branch struct {
	Rn rune
	Fn func(*Scanner)
}

type BranchRange struct {
	StartRn rune
	EndRn   rune
	Fn      func(*Scanner)
}

type BranchString struct {
	BranchString string
	Fn           func(*Scanner)
}

type Match struct {
	Literal string
	ID      string
}


// The user scan function should return the result of Exit()
//
func (s *Scanner) Exit() Token {
	if s.eof {
		return Token{"EOF", "EOF", s.curLine, s.curPos}
	}else if s.tok.ID == "" {
		return Token{"UNKNOWN", "UNKNOWN", s.curLine, s.curPos}
	}
	return s.tok
}

// Expect first reads a rune from s.read() and then tries to match it against input
// branches. If a match is found, the rune is 'consumed' (i.e. rune is put on the
// scanned buffer) and the branch is 'taken' (i.e. the branch function is called).
// If a match is not found, the read rune is then compared to each of the branch ranges.
//
// If a match is not found in either the branches or branch ranges and the Optional
// option is set to true, the scanner returns to the user scan function without calling
// any of the branch functions, but still consumes the rune.
//
// Any runes that are read but not consumed or skipped will be unread.
//
func (s *Scanner) Expect(expect ExpectRune) {
	var found1orMore bool
	var found1inverted bool
	var found bool
	var rn rune

    if logenb{
       s.log(fmt.Sprintf("Expect %v ", getScanOptions(expect.Options)), NEWLINE) //TODO Custom Print Function
	   s.log(fmt.Sprintf("Rune: %v ", branchesToStrings(expect.Branches)), NO_PREFIX)
	   s.log(fmt.Sprintf("Range: %v ", branchRangesToStrings(expect.BranchRanges)), NO_PREFIX)
    }
	for {
		found = false
		rn = s.read()
		for _, branch := range expect.Branches {
			if branch.Rn == rn {
				if !expect.Options.Invert{
					s.consume(rn, found1orMore, expect.Options.Skip)
				}
				found1orMore = true
				found = true
				s.callFn(branch.Fn)
				break
			}
		}
		if !found {
			for _, branch := range expect.BranchRanges {
				if branch.StartRn <= rn && rn <= branch.EndRn {
					if !expect.Options.Invert {
						s.consume(rn, found1orMore, expect.Options.Skip)
					}
					found1orMore = true
					found = true
					s.callFn(branch.Fn)
					break
				}
			}
		}
		if (!expect.Options.Invert && !found) || (expect.Options.Invert && found) {
			s.unread()
			break
		}
		if expect.Options.Invert && !found {
			s.consume(rn, found1inverted, expect.Options.Skip)
			found1inverted = true
		}
		if !expect.Options.Multiple || s.eof {
			break
		}
	}
	if !expect.Options.Invert && !found1orMore && !expect.Options.Optional {
		s.expRunes = append(s.expRunes, rn)
		strings := append(branchesToStrings(expect.Branches), branchRangesToStrings(expect.BranchRanges)...)
		s.error = s.newError(RUNE_EXPECTED_NOT_FOUND, fmt.Errorf("Found [%v], expected any of %v", string(rn), strings))
	} else if expect.Options.Invert && !found1inverted && !expect.Options.Optional {
		s.expRunes = append(s.expRunes, rn)
		strings := append(branchesToStrings(expect.Branches), branchRangesToStrings(expect.BranchRanges)...)
		s.error = s.newError(RUNE_EXPECTED_NOT_FOUND, fmt.Errorf("Found [%v], expected any except %v", string(rn), strings))
	}

	return
}

func (s *Scanner) consume(rn rune, found1orMore bool, skip bool) {
	if logenb{
        if !found1orMore {
			s.log(fmt.Sprintf("Pos:%v ", s.curPos), NO_PREFIX)
            s.log("Found: ", NO_PREFIX)
			s.log(sanitize(string(rn), true), NO_PREFIX)
        } else {
            s.log(", ", NO_PREFIX)
            s.log(sanitize(string(rn), true), NO_PREFIX)
        }
    }
	if !skip {
		s.expRunes = append(s.expRunes, rn)
	}
    s.curPos++
    if rn == '\n' {
       s.curLine++
	   s.curPos = 1
	   s.curLineBuffer.Reset()
	   if logenb{
		 s.log(fmt.Sprintf("Line %v:", s.curLine + 1), STARTLINE)
	   }
	}else {
		s.curLineBuffer.WriteRune(rn)
	}
}

func (s *Scanner) Call(fn func(*Scanner)) {
    s.log("Call", NEWLINE)
    s.callFn(fn)
}

func (s *Scanner) callFn(fn func(*Scanner)) {
	if fn != nil {
        if logenb{
            s.log("Scanning: "+getFuncName(fn), INCREMENT)
        }
		fn(s)
		if logenb{
            s.log("Returning: "+getFuncName(fn), DECREMENT)
        }
	}
}

// Match is required to be called by the user scan function before it
// returns to the user parse function, otherwise it will return the
// token NOT_MATCHED. Match will match every rune currently accepted by
// Expect() and not skipped (s.scanStr), against the input string.
//
// Once matched a Token is generated from the input ID, s.scanStr and
// the current line and position of the scanner. Once the user scan
// function has matched a token, any subsequent calls to Match will
// do nothing until the user scan function returns and is called again
// and reset (by s.init()) by the parser.
//
func (s *Scanner) Match(matches []Match) {
	if s.tok.ID != "" || s.eof {
		return
	}
    expString := runesToString(s.expRunes)
	for _, match := range matches {
		if expString == match.Literal || match.Literal == "" {
			if logenb{
                s.log("Matched: "+match.ID, NEWLINE)
                s.log(" - ", NO_PREFIX)
                s.log(sanitize(expString,true), NO_PREFIX)
            }
			s.tok = Token{match.ID, expString, s.curLine, s.curPos - len(expString)}
			break
		}
	}
}

func (s *Scanner) SkipRune() {
    if logenb {
        s.log("Skip Rune: ", NEWLINE)    
    }
	if len(s.expRunes) > 0 {
        rn := s.expRunes[len(s.expRunes)-1]
		s.expRunes = s.expRunes[:len(s.expRunes)-1]
        if logenb {
            s.log(sanitize(string(rn), true) + ", ", NO_PREFIX)
        }
	} else {
        if logenb {
		  s.log("Warning: No Runes to Skip", ERROR)
        }
	}
}

// Creates a new error and passes it to the parser. Only one error is generated by the
// scanner as it exits immediately after an error
//
func (s *Scanner) newError(code ErrorCode, err error) *Error {
	s.log(err.Error(), ERROR)
	errLength := len(s.expRunes)
	errStartPos := s.curPos - errLength
	errEndPos := s.curPos
	if s.error == nil {
		s.error = &Error{
			code,
			err,
			s.getLine(),
			s.curLine,
			errStartPos,
			errEndPos,
		}
		return s.error
	} 
	// else if s.error.Line != s.curLine {
	// 	s.error = &Error{
	// 		code,
	// 		err,
	// 		s.getLine(),
	// 		s.curLine,
	// 		s.curPos,
	// 	}
	// 	return s.error
	// }
	return nil
}

// scan is the entry point from the parser.
//
func (s *Scanner) scan() (Token, *Error) {
	s.init()
    if logenb{
	   s.log("Scanning: "+getFuncName(s.fn), INCREMENT) 
	   defer s.log("Returning: "+getFuncName(s.fn), DECREMENT) // use defer keyword to log after the fn has returned
    }
    return s.fn(s), s.error                                 // Call the user ScanFunc with a reference to the p.s scanner
}

// read reads the next rune from the bufferred reader. Only read from the
// bufio reader s.r if it hasn't already been read. Using another buffer
// s.buf means we can read and unread as many runes as we like.
//
func (s *Scanner) read() rune {

	if s.buf.unread > 0{
		rn := s.buf.runes[len(s.buf.runes)-s.buf.unread]
		s.buf.unread--
		return rn
	}
	
	rn, _, err := s.r.ReadRune() // We don't use s.r.UnreadRune as it can only be called once
    
	s.buf.runes = append(s.buf.runes, rn)
    
	// Assume an err means we have reached End of File
	if err != nil {
		s.eof = true // Used by p.Expect() to break out of Multiple + Inverted calls
		return rune(0)
	}
	return rn
}

// To unread simply increment the index to the rune buffer
func (s *Scanner) unread() {
	if s.buf.unread < len(s.buf.runes) { // Ensure we don't unread more runes than have been read
		s.buf.unread++
	}
}

// Reset the scanner after every s.callFn() call
//
func (s *Scanner) init() {
	s.tok.ID = ""
	s.expRunes = nil
}

// log is where all lines are added to the log.
// It is invoked with a number of indent options.
//
func (s *Scanner) log(msg string, indent indent) {
    if s.logger == nil {
		return
	}
    l := s.logger
	switch indent {
	case INCREMENT:
		{
			if l.buf != nil{
                l.log.Print(l.buf...)
                l.buf = nil
            }
            l.indent++
			prefix := ""
			for i := 0; i < l.indent; i++ {
				prefix += "\t"
			}
			l.log.SetPrefix(prefix)
		}
	case DECREMENT:
		{
            if l.buf != nil{
                l.log.Print(l.buf...)
                l.buf = nil
            }
			l.indent--
			prefix := ""
			for i := 0; i < l.indent; i++ {
				prefix += "\t"
			}
            l.log.Print(msg)
            l.log.SetPrefix(prefix)
            return
		}
	case NEWLINE:
		{
			if l.buf != nil{
                l.log.Print(l.buf...)
                l.buf = nil
            }
		}
	case NO_PREFIX:
		{
            //noop
		}
	case STARTLINE:
		{
            if l.buf != nil{
                l.log.Print(l.buf...)
                l.buf = nil
            }
            prefix := l.log.Prefix()
            l.log.SetPrefix("")
            l.log.Print(msg)
            l.log.SetPrefix(prefix)
            return
		}
	case ERROR:
		{           
            if l.buf != nil{
                l.log.Print(l.buf...)
                l.buf = nil
            }
            prefix := l.log.Prefix()
            l.log.SetPrefix("***")
            l.log.Print(msg)
            l.log.SetPrefix(prefix)
            return
		}
	}
	l.buf = append(l.buf,msg)
}

// getLine is used to scan the rest of the current line to display in the Error
//
func (s *Scanner) getLine() string {
	var numRunes int
	var tempBuffer bytes.Buffer
	
	tempBuffer.WriteString(s.curLineBuffer.String())

	for {
		rn := s.read()
		numRunes++
		if rn == '\n' || rn == rune(0) { // Make sure you break on New Line or End of File
			break
		} else {
			tempBuffer.WriteRune(rn)
		}
	}
	// Put all read runes back onto the buffer
	for i := 0; i < numRunes; i++ {
		s.unread()
	}
	s.unread()     // Unread the rune that was skipped

	return tempBuffer.String()
}

// -------------------------------- Scanner Helper Functions---------------------------------------

// Used to log which options were used during a branch function call
func getScanOptions(opts ScanOptions) string {
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

// Used to list possible scan branches in the log and errors in Expect()
func runesToString(runes []rune) (str string) {
	for _, rn := range runes {
		str = str + string(rn)
	}
	return
}

// Used to list possible scan branches in the log and errors in Expect()
func branchesToStrings(branches []Branch) (branchStrings []string) {
	for _, branch := range branches {
		branchStrings = append(branchStrings, sanitize(string(branch.Rn), true))
	}
	return
}

// Similar to runesToStrings except that it accepts the BranchString type used in Peek()
func branchStringsToStrings(branches []BranchString) (branchStrings []string) {
	for _, branch := range branches {
		branchStrings = append(branchStrings, branch.BranchString)
	}
	return
}

// Similar to runesToStrings except that it accepts the BranchRange type used in Expect()
func branchRangesToStrings(branches []BranchRange) (branchStrings []string) {
	for _, branch := range branches {
		branchStrings = append(branchStrings, sanitize(string(branch.StartRn), true)+"-"+sanitize(string(branch.EndRn), true))
	}
	return
}
