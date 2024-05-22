// Copyright (c) 2024 Dez Little <deslittle@gmail.com>
// All rights reserved. Use of this source code is governed by a LGPL v3
// license that can be found in the LICENSE file.

// token.go defines what a Token is and the token ID interface
package dsl

import "fmt"

// Line is the line of the source text the Token was found. Position is the
// position (or column) the Token was found. This information is used when
// displaying errors but could also be useful to the user for things like
// syntax highlighting and debugging if they were to implement it.
type Token struct {
	ID       string
	Literal  string
	Line     int
	Position int
}

type TokenSet map[string]int

func NewTokenSet(userIds ...string) TokenSet {
	ts := make(map[string]int)
	ts["UNKNOWN"] = 1
	for i, id := range userIds {
		if ts[id] != 0 {
			panic(fmt.Sprintf("Duplicate token ID found (%v)", id))
		}
		ts[id] = i + 2
	}
	return ts
}
