// scanner.go implements a Scanner for any DSL source text.
// When a Scanner is created it takes a bufio.Reader as the source of
// the text, scans a number of characters (runes) as defined by the
// user ScanFunc and returns a token to the Parser.
package dsl

import (
	"bufio"
	"bytes"
	"fmt"
)

// The Scanner contains a reference to the user scan function, the
// user input buffer, various state variables and the parser log.
type Scanner struct {
	fn  ScanFunc
	r   *bufio.Reader
	l   logger
	buf struct {
		runes  []rune
		unread int
	}
	curLineBuffer bytes.Buffer
	peekBuffer    []rune
	startLine     int
	curLine       int
	startPos      int
	curPos        int
	options       ExpectRuneOptions
	expRunes      []rune
	tok           Token
	error         *Error
	eof           bool
}

type ScanFunc func(*Scanner) Token

type Branch struct {
	Rn rune
	Fn func(*Scanner)
}

type BranchRange struct {
	StartRn rune
	EndRn   rune
	Fn      func(*Scanner)
}

type Match struct {
	Literal string
	ID      TokenType
}

// NewScanner returns a new instance of Scanner.
func newScanner(sf ScanFunc, r *bufio.Reader, l logger) *Scanner {
	s := &Scanner{
		fn:      sf,
		r:       r,
		l:       l,
		curLine: 1,
		curPos:  1,
	}

	return s
}

// If the Optional option is false and a match is not found, an error is returned to the
// parser.
//
// If the Multiple option is set to true the scanner continues to read, consume runes and
// take branches until a rune is read that is not matched by any of the branches or branch
// ranges. If the Multiple option is set to false, only the first branch (or branch range)
// to be matched is taken and consumed.
//
// If the the Skip option is set to true the scanner will take the branch if a match is
// found but will not consume the rune.
//
// If ExpectOptions is omitted when creating ExpectRune{}, all options will be set to
// false.
type ExpectRuneOptions struct {
	Optional bool
	Multiple bool
	Skip     bool
	Peek     bool
}

type ExpectRune struct {
	Branches     []Branch
	BranchRanges []BranchRange
	Options      ExpectRuneOptions
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
func (s *Scanner) Expect(expect ExpectRune) {

	s.log(fmt.Sprintf("Expect %v ", getExpectRuneOptions(expect.Options)), prefixNewline)
	s.log(fmt.Sprintf("Rune: %v ", branchesToStrings(expect.Branches)), prefixNone)
	s.log(fmt.Sprintf("Range: %v ", branchRangesToStrings(expect.BranchRanges)), prefixNone)
	if s.tok.ID != "" {
		s.log("Already Matched. Skipping.", prefixNewline)
		return
	}

	var found1orMore bool
	var rn rune

	// Loop if Multiple is true
	for {
		var branchFn func(*Scanner)
		found := false
		rn = s.read()

		// Check if the current rune matches any of the branches
		for _, branch := range expect.Branches {
			if branch.Rn == rn {
				found = true
				branchFn = branch.Fn
				break
			}
		}

		// Check if the current rune matches any of the branch ranges
		if !found {
			for _, branch := range expect.BranchRanges {
				if branch.StartRn <= rn && rn <= branch.EndRn {
					found = true
					branchFn = branch.Fn
					break
				}
			}
		}
		if !found {
			// If no match was found, unread the rune and break out of the loop
			s.unread()
			break
		}

		// We found a match
		if expect.Options.Peek {
			// If we are peeking, remember each rune read
			s.peekBuffer = append(s.peekBuffer, rn)
		}
		if !expect.Options.Peek {
			// If we are not peeking, consume the rune
			if len(s.peekBuffer) > 0 {
				// If we have peeked but now no longer peeking, consume the peeked runes
				s.consumePeeked()
			}
			s.consume(rn, expect.Options.Skip)
		}
		s.logMatch(rn, found1orMore)
		found1orMore = true // Set to true only after logging the first match
		s.scanFn(branchFn)

		if !expect.Options.Multiple {
			// If Multiple is false break out of the loop
			break
		}
	}

	// If we have finished peeking and there are runes to unread, unread them
	if expect.Options.Peek && len(s.peekBuffer) > 0 {
		s.unreadPeeked()
	}

	if !found1orMore && !expect.Options.Optional {
		s.expRunes = append(s.expRunes, rn)
		strings := append(branchesToStrings(expect.Branches), branchRangesToStrings(expect.BranchRanges)...)
		s.error = s.newError(ErrorRuneExpectedNotFound, fmt.Errorf("found [%v], expected any of %v", string(rn), strings))
	}

}

type RuneRange struct {
	StartRn rune
	EndRn   rune
}

type ExpectNotRune struct {
	Runes      []rune
	RuneRanges []RuneRange
	Fn         func(*Scanner)
	Options    ExpectRuneOptions
}

// ExpectNot reads a rune from the buffer and then tries to match it against the input
// runes or ranges. If a match is found, the rune is 'consumed' (i.e. rune is put on the
// scanned buffer) and the branch is 'taken' (i.e. the branch function is called).
// If a match is not found, the read rune is then compared to each of the ranges.
//
// If a match is not found in either the runes or ranges and the Optional option is set
// to true, the scanner returns to the user scan function without calling any of the
// branch functions, but still consumes the rune.
//
// Any runes that are read but not consumed or skipped will be unread.
func (s *Scanner) ExpectNot(expect ExpectNotRune) {
	s.log(fmt.Sprintf("ExpectNot %v ", getExpectRuneOptions(expect.Options)), prefixNewline)
	s.log(fmt.Sprintf("Rune: %v ", runesToStrings(expect.Runes)), prefixNone)
	s.log(fmt.Sprintf("Range: %v ", runeRangesToStrings(expect.RuneRanges)), prefixNone)
	if s.tok.ID != "" {
		s.log("Already Matched. Skipping.", prefixNewline)
		return
	}

	var found1orMoreNot bool
	var rn rune

	for {
		found := false
		rn = s.read()
		for _, expectedRune := range expect.Runes {
			if expectedRune == rn {
				found = true
				break
			}
		}
		if !found {
			for _, expectedRuneRange := range expect.RuneRanges {
				if expectedRuneRange.StartRn <= rn && rn <= expectedRuneRange.EndRn {
					found = true
					break
				}
			}
		}

		// If a match was found, unread the rune and break out of the loop
		if found {
			s.unread()
			break
		}

		// Match not found, which is what we are expecting

		if expect.Options.Peek {
			// If we are peeking, remember each rune read
			s.peekBuffer = append(s.peekBuffer, rn)
		} else {
			// If we are not peeking, consume the rune
			if len(s.peekBuffer) > 0 {
				// If we have peeked but now no longer peeking, consume the peeked runes
				s.consumePeeked()
			}
			s.consume(rn, expect.Options.Skip)
		}

		s.logMatch(rn, found1orMoreNot)
		found1orMoreNot = true // Set to true only after logging the first match
		s.scanFn(expect.Fn)

		if !expect.Options.Multiple {
			break
		}
	}

	// If we have finished peeking and there are runes to unread, unread them
	if expect.Options.Peek && len(s.peekBuffer) > 0 {
		s.unreadPeeked()
	}

	if !found1orMoreNot && !expect.Options.Optional {
		s.expRunes = append(s.expRunes, rn)
		strings := append(runesToStrings(expect.Runes), runeRangesToStrings(expect.RuneRanges)...)
		s.error = s.newError(ErrorRuneExpectedNotFound, fmt.Errorf("found [%v], expected any except %v", string(rn), strings))
	}

}

func (s *Scanner) Call(fn func(*Scanner)) {
	s.log("Calling: "+getFuncName(fn), prefixIncrement)
	fn(s)
	s.log("Returning: "+getFuncName(fn), prefixDecrement)
}

// Match is required to be called by the user scan function before it
// returns to the user parse function, otherwise it will return the
// token UNKNOWN. Match will match every rune currently accepted by
// Expect() and not skipped (s.scanStr), against the input string.
//
// Once matched a Token is generated from the input ID, s.scanStr and
// the current line and position of the scanner. Once the user scan
// function has matched a token, any subsequent calls to Match will
// do nothing until the user scan function returns and is called again
// and reset (by s.init()) by the parser.
func (s *Scanner) Match(matches []Match) {
	if s.tok.ID != "" {
		return
	}
	expString := runesToString(s.expRunes)
	for _, match := range matches {
		if expString == match.Literal || match.Literal == "" {
			s.log("Matched: "+string(match.ID)+" - "+sanitize(expString, true), prefixNewline)
			s.tok = Token{match.ID, expString, s.curLine, s.curPos - len(s.expRunes)}
			break
		}
	}
}

// The user scan function should return the result of Exit()
func (s *Scanner) Exit() Token {
	if s.tok.ID == "" {
		return Token{
			ID:       TOKEN_UNKNOWN,
			Literal:  "UNKNOWN",
			Line:     s.curLine,
			Position: s.curPos,
		}
	}
	return s.tok
}

func (s *Scanner) SkipRune() {
	s.log("Skip Rune: ", prefixNewline)
	if len(s.expRunes) > 0 {
		rn := s.expRunes[len(s.expRunes)-1]
		s.expRunes = s.expRunes[:len(s.expRunes)-1]
		s.log(sanitize(string(rn), true)+", ", prefixNone)
	} else {
		s.log("Warning: No Runes to Skip", prefixError)
	}
}

// -------------------------------- scanner interface ---------------------------------------

type scanner interface {
	scan() (Token, string, *Error)
}

// scan is the entry point from the parser.
func (s *Scanner) scan() (Token, string, *Error) {
	s.init()
	s.log("Scanning: "+getFuncName(s.fn), prefixIncrement)
	defer s.log("Returning: "+getFuncName(s.fn), prefixDecrement) // use defer keyword to log after the fn has returned
	return s.fn(s), s.getLine(), s.error                          // Call the user ScanFunc with a reference to the p.s scanner
}

// -------------------------------- Scanner Core Functions---------------------------------------

func (s *Scanner) logMatch(rn rune, found1orMore bool) {
	if !found1orMore {
		s.log(fmt.Sprintf("Pos:%v ", s.curPos), prefixNone)
		s.log("Found: ", prefixNone)
		s.log(sanitize(string(rn), true), prefixNone)
	} else {
		s.log(", ", prefixNone)
		s.log(sanitize(string(rn), true), prefixNone)
	}
	if rn == '\n' {
		s.log(fmt.Sprintf("Line %v:", s.curLine), prefixStartLine)
	}
}

func (s *Scanner) consume(rn rune, skip bool) {
	if !skip {
		s.expRunes = append(s.expRunes, rn)
	}
	s.curPos++

	if rn == '\n' {
		s.curLine++
		s.curPos = 1
		s.curLineBuffer.Reset()
	} else {
		s.curLineBuffer.WriteRune(rn)
	}
}

// read reads the next rune from the bufferred reader. Only read from the
// bufio reader s.r if it hasn't already been read. Using another buffer
// s.buf means we can read and unread as many runes as we like.
func (s *Scanner) read() rune {

	if s.buf.unread > 0 {
		rn := s.buf.runes[len(s.buf.runes)-s.buf.unread]
		s.buf.unread--
		return rn
	}

	rn, _, err := s.r.ReadRune() // We don't use s.r.UnreadRune as it can only be called once

	s.buf.runes = append(s.buf.runes, rn)

	// Assume an err means we have reached End of File
	if err != nil {
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

// unreadPeeked is used to unread the peeked runes
func (s *Scanner) unreadPeeked() {
	for i := len(s.peekBuffer) - 1; i >= 0; i-- {
		s.unread()
	}
	s.peekBuffer = nil
}

// consumePeeked is used to consume the peeked runes
func (s *Scanner) consumePeeked() {
	for _, rn := range s.peekBuffer {
		s.consume(rn, false)
	}
	s.peekBuffer = nil
}

func (s *Scanner) scanFn(fn func(*Scanner)) {
	if fn != nil {
		if s.tok.ID != "" {
			s.log("Already Matched. Skipping.", prefixNewline)
			return
		}
		s.log("Scanning: "+getFuncName(fn), prefixIncrement)
		fn(s)
		s.log("Returning: "+getFuncName(fn), prefixDecrement)
	}
}

// -------------------------------- Scanner Helper Functions---------------------------------------

// Reset the scanner after every s.callFn() call
func (s *Scanner) init() {
	s.tok.ID = ""
	s.expRunes = nil
	s.error = nil
	s.startLine = s.curLine
	s.startPos = s.curPos
	s.peekBuffer = nil
	s.curLineBuffer.Reset()
}

// Creates a new error and passes it to the parser. Only one error is generated by the
// scanner as it exits immediately after an error
func (s *Scanner) newError(code ErrorCode, err error) *Error {
	s.log(err.Error(), prefixError)

	if s.error == nil {
		return &Error{
			Code:          code,
			Message:       err.Error(),
			LineString:    s.curLineBuffer.String(),
			StartLine:     s.startLine,
			StartPosition: s.startPos,
			EndLine:       s.curLine,
			EndPosition:   s.curPos,
		}
	}
	return nil
}

// log is where all lines are added to the log.
// It is invoked with a number of indent options.
func (s *Scanner) log(msg string, indent indent) {
	s.l.log(msg, indent)
}

// getLine is used to scan the rest of the current line to display in the Error
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

	return tempBuffer.String()
}

// Used to log which options were used during a branch function call
func getExpectRuneOptions(opts ExpectRuneOptions) string {
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

// Used to calculate token literal strings
func runesToString(runes []rune) (str string) {
	for _, rn := range runes {
		if rn != rune(0) {
			str = str + string(rn)
		}
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

// Similar to runesToStrings except that it accepts the BranchRange type used in Expect()
func branchRangesToStrings(branches []BranchRange) (branchStrings []string) {
	for _, branch := range branches {
		branchStrings = append(branchStrings, sanitize(string(branch.StartRn), true)+"-"+sanitize(string(branch.EndRn), true))
	}
	return
}

// Used to list possible scan branches in the log and errors
func runesToStrings(runes []rune) (runesStrings []string) {
	for _, rn := range runes {
		runesStrings = append(runesStrings, sanitize(string(rn), true))
	}
	return
}

// Used to list possible scan branches in the log and errors
func runeRangesToStrings(runes []RuneRange) (runesStrings []string) {
	for _, rn := range runes {
		runesStrings = append(runesStrings, sanitize(string(rn.StartRn), true)+"-"+sanitize(string(rn.EndRn), true))
	}
	return
}
