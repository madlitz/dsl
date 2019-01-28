// Copyright (c) 2015 Des Little <deslittle@gmail.com>
// All rights reserved. Use of this source code is governed by a LGPL v3
// license that can be found in the LICENSE file.

package mydsl

import (
	"bufio"
	"bytes"
	"testing"
	"os"
	"encoding/json"
	"fmt"
    "github.com/Autoblocks/go-dsl"
)

func TestPrintAST(t *testing.T) {
	reader := bytes.NewBufferString(
		`a := 1 * 5 + 7
		b := 3.45 * 44.21 / (4 + a) 'A Simple Expression
		double(a + b)`)
	bufreader := bufio.NewReader(reader)
	ts := NewTokenSet()
    ns := NewNodeSet()
    logfilename := "log.txt"
    logfile, err := os.Create(logfilename)
    if err != nil {
		t.Fatal("Error: Could not create log file " + logfilename + ": " + err.Error())
	}
	ast, _ := dsl.ParseAndLog(Parse, Scan, ts, ns, bufreader, logfile)
	logfile.Close()
	
	astjson, _ := json.Marshal(ast)
	fmt.Print(astjson)
	if(string(astjson) != `{"root": "1"}`){
		t.Errorf("JSON malformed.")
	}
}

func TestDSL(t *testing.T) {
	
	reader := bytes.NewBufferString(
`a := 1 * 5 + 7
b := 3.45 * 44.21 / (4 + a) 'A Simple Expression
double(a + b)`)
	bufreader := bufio.NewReader(reader)
    ts := NewTokenSet()
    ns := NewNodeSet()
    logfilename := "log.txt"
    logfile, err := os.Create(logfilename)
    if err != nil {
		t.Fatal("Error: Could not create log file " + logfilename + ": " + err.Error())
	}
	ast, errs := dsl.ParseAndLog(Parse, Scan, ts, ns, bufreader, logfile)
    logfile.Close()
	cases := []dsl.Node {
		{Type: "TERMINAL", Tokens: []dsl.Token{{"LITERAL", "1", 1, 6}}},
		{Type: "TERMINAL", Tokens: []dsl.Token{{"LITERAL", "5", 1, 10}}},
		{Type: "TERMINAL", Tokens: []dsl.Token{{"LITERAL", "7", 1, 14}}},
		{Type: "EXPRESSION", Tokens: []dsl.Token{{"PLUS", "+", 1, 12}}},
		{Type: "EXPRESSION", Tokens: []dsl.Token{{"MULTIPLY", "*", 1, 8}}},
		{Type: "ASSIGNMENT", Tokens: []dsl.Token{{"VARIABLE", "a", 1, 1}}},
		{Type: "TERMINAL", Tokens: []dsl.Token{{"LITERAL", "3.45", 2, 6}}},
		{Type: "TERMINAL", Tokens: []dsl.Token{{"LITERAL", "44.21", 2, 13}}},
		{Type: "TERMINAL", Tokens: []dsl.Token{{"LITERAL", "4", 2, 22}}},
		{Type: "TERMINAL", Tokens: []dsl.Token{{"VARIABLE", "a", 2, 26}}},
		{Type: "EXPRESSION", Tokens: []dsl.Token{{"PLUS", "+", 2, 24}}},
		{Type: "EXPRESSION", Tokens: []dsl.Token{{"OPEN_PAREN", "(", 2, 21}}},
		{Type: "TERMINAL", Tokens: []dsl.Token{{"CLOSE_PAREN", ")", 2, 27}}},
		{Type: "EXPRESSION", Tokens: []dsl.Token{{"DIVIDE", "/", 2, 19}}},
		{Type: "EXPRESSION", Tokens: []dsl.Token{{"MULTIPLY", "*", 2, 11}}},
		{Type: "ASSIGNMENT", Tokens: []dsl.Token{{"VARIABLE", "b", 2, 1}}},
		{Type: "COMMENT", Tokens: []dsl.Token{{"COMMENT", "A Simple Expression", 2, 30}}},
		{Type: "TERMINAL", Tokens: []dsl.Token{{"VARIABLE", "a", 3, 8}}},
		{Type: "TERMINAL", Tokens: []dsl.Token{{"VARIABLE", "b", 3, 12}}},
		{Type: "EXPRESSION", Tokens: []dsl.Token{{"PLUS", "+", 3, 10}}},
		{Type: "CALL", Tokens: []dsl.Token{{"VARIABLE", "double", 3, 1}}},
		{Type: "ROOT", Tokens: []dsl.Token{{"", "", 0, 0}}},
	}
	count := 0
	ast.Print()
	ast.Inspect(func(node *dsl.Node)(){
		if count > len(cases) - 1{
			t.Fatalf("Too many nodes.")
		}
		if cases[count].Type != node.Type{
			t.Errorf("Line: %v:%v Node: \"%v\" Wanted node type %v, found %v", cases[count].Tokens[0].Line, cases[count].Tokens[0].Position, 
                node.Type, cases[count].Type, node.Type)
		}
        for i, token := range node.Tokens{
           if cases[count].Tokens[i].ID != token.ID{
			 t.Errorf("Line: %v:%v Token: \"%v\" Wanted token ID %v, found %v", cases[count].Tokens[i].Line, cases[count].Tokens[i].Position, 
                token.Literal, cases[count].Tokens[i].ID, token.ID)
		   }
           if cases[count].Tokens[i].Literal != token.Literal{
			 t.Errorf("Line: %v:%v ID: \"%v\" Wanted token literal \"%v\", found \"%v\"", cases[count].Tokens[i].Line, cases[count].Tokens[i].Position, 
                token.ID, cases[count].Tokens[i].Literal, token.Literal)
		   }
           if cases[count].Tokens[i].Line != token.Line{
			 t.Errorf("Line: %v:%v Token: \"%v\" Wanted token line %v, found %v", cases[count].Tokens[i].Line, cases[count].Tokens[i].Position, 
                token.Literal, cases[count].Tokens[i].Line, token.Line)
		   }
           if cases[count].Tokens[i].Position != token.Position{
			 t.Errorf("Line: %v:%v Token: \"%v\" Wanted token position %v, found %v", cases[count].Tokens[i].Line, cases[count].Tokens[i].Position, 
                token.Literal, cases[count].Tokens[i].Position, token.Position)
		   } 
        }
		count++
	})
	if count != len(cases){
		t.Errorf("Not enough nodes.")
	}

	if errs != nil {
		t.Fail()
		for _, err := range errs {
			t.Error(err.String())
		}
	}

}



