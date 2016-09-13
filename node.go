// Copyright (c) 2015 Des Little <deslittle@gmail.com>
// All rights reserved. Use of this source code is governed by a LGPL v3
// license that can be found in the LICENSE file.

// token.go defines what a Token is and the token ID interface
//
package dsl

// A Node can contain multiple Tokens which can be useful if the user knows how
// many Tokens belong to a particular Node type. Otherwise, the user should only
// add one token per node.
//
type Node struct {
	Type     string
	Tokens   []Token
	Parent   *Node
	Children []*Node
}

type NodeSet map[string]int

func NewNodeSet(userTypes ...string) NodeSet {
	ns := make(map[string]int)
	ns["ROOT"] = 1
	for i, id := range userTypes {
		ns[id] = i + 2
	}
	return ns
}
