package dsl

import (
	"bytes"
	"fmt"
)

type ErrorCode int

const (
	ErrorFileNotFound ErrorCode = iota // 0
	ErrorCouldNotCreateFile
	ErrorTokenExpectedNotFound
	ErrorRuneExpectedNotFound
	ErrorNodeNotInNodeSet
	ErrorNoTokensToGet
)

// Error contains the error text, the line and positions the error occurred on, and
// a string containing the input text from that line.
type Error struct {
	Code          ErrorCode
	Message       string
	LineString    string
	StartLine     int
	StartPosition int
	EndLine       int
	EndPosition   int
}

// Error implements the error interface.
func (e *Error) Error() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("\nError Line:%v %v\n", e.StartLine, e.Message))
	buf.WriteString(e.LineString + "\n")
	for i := 0; i < e.StartPosition-1; i++ {
		if i < len(e.LineString) && e.LineString[i] == '\t' {
			buf.WriteString("\t")
		} else {
			buf.WriteString(" ")
		}
	}
	buf.WriteString("^")
	if e.StartLine != e.EndLine {
		for i := e.StartPosition; i < len(e.LineString); i++ {
			buf.WriteString("-")
		}
	} else {
		for i := e.StartPosition; i < e.EndPosition-1; i++ {
			buf.WriteString("-")
		}
	}
	buf.WriteString("^\n")
	return buf.String()
}

// NewError creates a new Error instance.
func NewError(code ErrorCode, message, lineString string, startLine, startPosition, endLine, endPosition int) *Error {
	return &Error{
		Code:          code,
		Message:       message,
		LineString:    lineString,
		StartLine:     startLine,
		StartPosition: startPosition,
		EndLine:       endLine,
		EndPosition:   endPosition,
	}
}
