// Copyright (c) 2015 Des Little <deslittle@gmail.com>
// All rights reserved. Use of this source code is governed by a LGPL v3
// license that can be found in the LICENSE file.

package dsl

import (
	"fmt"
	"bytes"
)

type ErrorCode int

const (
	FILE_NOT_FOUND ErrorCode = iota // 0
	COULD_NOT_CREATE_FILE
	TOKEN_EXPECTED_NOT_FOUND
	EXPECTED_TOKEN_NOT_IN_TOKENSET
	SCANNED_TOKEN_NOT_IN_TOKENSET
	RUNE_EXPECTED_NOT_FOUND
	NODE_NOT_IN_NODESET
)

// Errors contain the error text, the line and position the error occurred on, and
// a string containing the input text from that line
type Error struct {
	Code       ErrorCode
	Error      error
	LineString string
	Line       int
	Position   int
}

func (e *Error) String() string{
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("\nError Line:%v %v\n", e.Line, e.Error))
	buf.WriteString(e.LineString + "\n")
	for i := 0; i < e.Position - 1; i++ {
		if e.LineString[i] == '\t' {
			buf.WriteString("\t")
		} else {
			buf.WriteString(" ")
		}
	}
	buf.WriteString("^")
	return buf.String()
}
