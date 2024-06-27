package mydsl

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/madlitz/go-dsl"
)

func TestPrintAST(t *testing.T) {
	reader := bytes.NewBufferString(
		`a := 1 * 5 + 7
		b := 3.45 * 44.21 / (4 + a) 'A Simple Expression
		double(a + b)`)
	bufreader := bufio.NewReader(reader)
	logfilename := "logs/TestPrintAST.log"
	logfile, err := os.Create(logfilename)
	if err != nil {
		t.Fatal("Error: Could not create log file " + logfilename + ": " + err.Error())
	}
	ast, _ := dsl.Parse(Parse, Scan, bufreader, dsl.WithLogger(logfile))
	logfile.Close()

	astjson, _ := json.Marshal(ast)
	expectedJson := `{
		"root": {
			"type": "ROOT",
			"tokens": null,
			"children": [
				{
					"type": "ASSIGNMENT",
					"tokens": [
						{
							"ID": "VARIABLE",
							"Literal": "a",
							"Line": 1,
							"Position": 1
						}
					],
					"children": [
						{
							"type": "TERMINAL",
							"tokens": [
								{
									"ID": "LITERAL",
									"Literal": "1",
									"Line": 1,
									"Position": 6
								}
							],
							"children": null
						},
						{
							"type": "EXPRESSION",
							"tokens": [
								{
									"ID": "MULTIPLY",
									"Literal": "*",
									"Line": 1,
									"Position": 8
								}
							],
							"children": [
								{
									"type": "TERMINAL",
									"tokens": [
										{
											"ID": "LITERAL",
											"Literal": "5",
											"Line": 1,
											"Position": 10
										}
									],
									"children": null
								},
								{
									"type": "EXPRESSION",
									"tokens": [
										{
											"ID": "PLUS",
											"Literal": "+",
											"Line": 1,
											"Position": 12
										}
									],
									"children": [
										{
											"type": "TERMINAL",
											"tokens": [
												{
													"ID": "LITERAL",
													"Literal": "7",
													"Line": 1,
													"Position": 14
												}
											],
											"children": null
										}
									]
								}
							]
						}
					]
				},
				{
					"type": "ASSIGNMENT",
					"tokens": [
						{
							"ID": "VARIABLE",
							"Literal": "b",
							"Line": 2,
							"Position": 3
						}
					],
					"children": [
						{
							"type": "TERMINAL",
							"tokens": [
								{
									"ID": "LITERAL",
									"Literal": "3.45",
									"Line": 2,
									"Position": 8
								}
							],
							"children": null
						},
						{
							"type": "EXPRESSION",
							"tokens": [
								{
									"ID": "MULTIPLY",
									"Literal": "*",
									"Line": 2,
									"Position": 13
								}
							],
							"children": [
								{
									"type": "TERMINAL",
									"tokens": [
										{
											"ID": "LITERAL",
											"Literal": "44.21",
											"Line": 2,
											"Position": 15
										}
									],
									"children": null
								},
								{
									"type": "EXPRESSION",
									"tokens": [
										{
											"ID": "DIVIDE",
											"Literal": "/",
											"Line": 2,
											"Position": 21
										}
									],
									"children": [
										{
											"type": "EXPRESSION",
											"tokens": [
												{
													"ID": "OPEN_PAREN",
													"Literal": "(",
													"Line": 2,
													"Position": 23
												}
											],
											"children": [
												{
													"type": "TERMINAL",
													"tokens": [
														{
															"ID": "LITERAL",
															"Literal": "4",
															"Line": 2,
															"Position": 24
														}
													],
													"children": null
												},
												{
													"type": "EXPRESSION",
													"tokens": [
														{
															"ID": "PLUS",
															"Literal": "+",
															"Line": 2,
															"Position": 26
														}
													],
													"children": [
														{
															"type": "TERMINAL",
															"tokens": [
																{
																	"ID": "VARIABLE",
																	"Literal": "a",
																	"Line": 2,
																	"Position": 28
																}
															],
															"children": null
														}
													]
												}
											]
										},
										{
											"type": "TERMINAL",
											"tokens": [
												{
													"ID": "CLOSE_PAREN",
													"Literal": ")",
													"Line": 2,
													"Position": 29
												}
											],
											"children": null
										}
									]
								}
							]
						}
					]
				},
				{
					"type": "COMMENT",
					"tokens": [
						{
							"ID": "COMMENT",
							"Literal": "A Simple Expression",
							"Line": 2,
							"Position": 32
						}
					],
					"children": null
				},
				{
					"type": "CALL",
					"tokens": [
						{
							"ID": "VARIABLE",
							"Literal": "double",
							"Line": 3,
							"Position": 3
						}
					],
					"children": [
						{
							"type": "TERMINAL",
							"tokens": [
								{
									"ID": "VARIABLE",
									"Literal": "a",
									"Line": 3,
									"Position": 10
								}
							],
							"children": null
						},
						{
							"type": "EXPRESSION",
							"tokens": [
								{
									"ID": "PLUS",
									"Literal": "+",
									"Line": 3,
									"Position": 12
								}
							],
							"children": [
								{
									"type": "TERMINAL",
									"tokens": [
										{
											"ID": "VARIABLE",
											"Literal": "b",
											"Line": 3,
											"Position": 14
										}
									],
									"children": null
								}
							]
						}
					]
				}
			]
		}
	}`

	var actual interface{}
	if err := json.Unmarshal(astjson, &actual); err != nil {
		t.Fatal("Error: Could not unmarshal AST JSON: " + err.Error())
	}

	var expected interface{}
	if err := json.Unmarshal([]byte(expectedJson), &expected); err != nil {
		t.Fatal("Error: Could not unmarshal expected JSON: " + err.Error())
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("JSON mismatch (-expected +actual):\n%s", diff)
	}
}

func TestDSL(t *testing.T) {
	reader := bytes.NewBufferString(
		`a := 1 * 5 + 7
b := 3.45 * 44.21 / (4 + a) 'A Simple Expression
double(a + b)`)
	bufreader := bufio.NewReader(reader)
	logfilename := "logs/TestDSL.log"
	logfile, fileErr := os.Create(logfilename)
	if fileErr != nil {
		t.Fatal("Error: Could not create log file " + logfilename + ": " + fileErr.Error())
	}
	ast, errs := dsl.Parse(Parse, Scan, bufreader, dsl.WithLogger(logfile))
	logfile.Close()
	if len(errs) != 0 {
		t.Fatalf("Should report exactly 0 errors: got %d", len(errs))
	}

	expectedNodes := dsl.Node{
		Type: dsl.NODE_ROOT,
		Children: []dsl.Node{
			{
				Type: NODE_ASSIGNMENT,
				Tokens: []dsl.Token{
					{ID: TOKEN_VARIABLE, Literal: "a", Line: 1, Position: 1},
				},
				Children: []dsl.Node{
					{
						Type: NODE_TERMINAL,
						Tokens: []dsl.Token{
							{ID: TOKEN_LITERAL, Literal: "1", Line: 1, Position: 6},
						},
					},
					{
						Type: NODE_EXPRESSION,
						Tokens: []dsl.Token{
							{ID: TOKEN_MULTIPLY, Literal: "*", Line: 1, Position: 8},
						},
						Children: []dsl.Node{
							{
								Type: NODE_TERMINAL,
								Tokens: []dsl.Token{
									{ID: TOKEN_LITERAL, Literal: "5", Line: 1, Position: 10},
								},
							},
							{
								Type: NODE_EXPRESSION,
								Tokens: []dsl.Token{
									{ID: TOKEN_PLUS, Literal: "+", Line: 1, Position: 12},
								},
								Children: []dsl.Node{
									{
										Type: NODE_TERMINAL,
										Tokens: []dsl.Token{
											{ID: TOKEN_LITERAL, Literal: "7", Line: 1, Position: 14},
										},
									},
								},
							},
						},
					},
				},
			},
			{
				Type: NODE_ASSIGNMENT,
				Tokens: []dsl.Token{
					{ID: TOKEN_VARIABLE, Literal: "b", Line: 2, Position: 1},
				},
				Children: []dsl.Node{
					{
						Type: NODE_TERMINAL,
						Tokens: []dsl.Token{
							{ID: TOKEN_LITERAL, Literal: "3.45", Line: 2, Position: 6},
						},
					},
					{
						Type: NODE_EXPRESSION,
						Tokens: []dsl.Token{
							{ID: TOKEN_MULTIPLY, Literal: "*", Line: 2, Position: 11},
						},
						Children: []dsl.Node{
							{
								Type: NODE_TERMINAL,
								Tokens: []dsl.Token{
									{ID: TOKEN_LITERAL, Literal: "44.21", Line: 2, Position: 13},
								},
							},
							{
								Type: NODE_EXPRESSION,
								Tokens: []dsl.Token{
									{ID: TOKEN_DIVIDE, Literal: "/", Line: 2, Position: 19},
								},
								Children: []dsl.Node{
									{
										Type: NODE_EXPRESSION,
										Tokens: []dsl.Token{
											{ID: TOKEN_OPEN_PAREN, Literal: "(", Line: 2, Position: 21},
										},
										Children: []dsl.Node{
											{
												Type: NODE_TERMINAL,
												Tokens: []dsl.Token{
													{ID: TOKEN_LITERAL, Literal: "4", Line: 2, Position: 22},
												},
											},
											{
												Type: NODE_EXPRESSION,
												Tokens: []dsl.Token{
													{ID: TOKEN_PLUS, Literal: "+", Line: 2, Position: 24},
												},
												Children: []dsl.Node{
													{
														Type: NODE_TERMINAL,
														Tokens: []dsl.Token{
															{ID: TOKEN_VARIABLE, Literal: "a", Line: 2, Position: 26},
														},
													},
												},
											},
										},
									},
									{
										Type: NODE_TERMINAL,
										Tokens: []dsl.Token{
											{ID: TOKEN_CLOSE_PAREN, Literal: ")", Line: 2, Position: 27},
										},
									},
								},
							},
						},
					},
				},
			},
			{
				Type: NODE_COMMENT,
				Tokens: []dsl.Token{
					{ID: TOKEN_COMMENT, Literal: "A Simple Expression", Line: 2, Position: 30},
				},
			},
			{
				Type: NODE_CALL,
				Tokens: []dsl.Token{
					{ID: TOKEN_VARIABLE, Literal: "double", Line: 3, Position: 1},
				},
				Children: []dsl.Node{
					{
						Type: NODE_TERMINAL,
						Tokens: []dsl.Token{
							{ID: TOKEN_VARIABLE, Literal: "a", Line: 3, Position: 8},
						},
					},
					{
						Type: NODE_EXPRESSION,
						Tokens: []dsl.Token{
							{ID: TOKEN_PLUS, Literal: "+", Line: 3, Position: 10},
						},
						Children: []dsl.Node{
							{
								Type: NODE_TERMINAL,
								Tokens: []dsl.Token{
									{ID: TOKEN_VARIABLE, Literal: "b", Line: 3, Position: 12},
								},
							},
						},
					},
				},
			},
		},
	}

	if diff := cmp.Diff(
		expectedNodes,
		*ast.RootNode,
		cmpopts.IgnoreFields(dsl.Node{}, "Parent"), // Ignore Parent field as this is a pointer
	); diff != "" {
		t.Errorf("AST mismatch (-expected +actual):\n%s", diff)
	}

}

func TestTokenExpectedButNotFoundError(t *testing.T) {
	reader := bytes.NewBufferString(
		`a error := 1 * 5 + 7
b := 3.45 * 44.21 / (4 + a) 'A Simple Expression
double(a + b)  `)
	bufreader := bufio.NewReader(reader)
	logfilename := "logs/TestTokenExpectedButNotFoundError.log"
	logfile, fileErr := os.Create(logfilename)
	if fileErr != nil {
		t.Fatal("Error: Could not create log file " + logfilename + ": " + fileErr.Error())
	}
	_, errs := dsl.Parse(Parse, Scan, bufreader, dsl.WithLogger(logfile))

	if len(errs) != 1 {
		t.Fatalf("Should report exactly 1 error: got %d", len(errs))
	}
	err := errs[0]
	if err.Code != dsl.ERROR_TOKEN_EXPECTED_NOT_FOUND {
		t.Fail()
		t.Errorf("Expected error code 'Token expected but not found'. Found error: '%v", err)
	}
	if err.StartLine != 1 {
		t.Fail()
		t.Errorf("Expected error line 1. Found line: %v", err.StartLine)
	}
	if err.StartPosition != 3 {
		t.Fail()
		t.Errorf("Expected error start position 3. Found position: %v", err.StartPosition)
	}
	if err.EndPosition != 9 {
		t.Fail()
		t.Errorf("Expected error end position 9. Found position: %v", err.EndPosition)
	}

}

func TestRuneExpectedButNotFoundError(t *testing.T) {
	reader := bytes.NewBufferString(
		`_ := 1 * 5 + 7
b := 3.45 * 44.21 / (4 + a) 'A Simple Expression
double(a + b)`)
	bufreader := bufio.NewReader(reader)
	logfilename := "logs/TestRuneExpectedButNotFoundError.log"
	logfile, fileErr := os.Create(logfilename)
	if fileErr != nil {
		t.Fatal("Error: Could not create log file " + logfilename + ": " + fileErr.Error())
	}
	_, errs := dsl.Parse(Parse, Scan, bufreader, dsl.WithLogger(logfile))

	if len(errs) != 1 {
		t.Fatalf("Should report exactly 1 error: got %d", len(errs))
	}
	err := errs[0]
	if err.Code != dsl.ERROR_RUNE_EXPECTED_NOT_FOUND {
		t.Fail()
		t.Errorf("Expected error code 'Rune expected but not found'. Found error: '%v", err)
	}
	if err.StartLine != 1 {
		t.Fail()
		t.Errorf("Expected error line 1. Found line: %v", err.StartLine)
	}
	if err.StartPosition != 1 {
		t.Fail()
		t.Errorf("Expected error start position 1. Found position: %v", err.StartPosition)
	}
	if err.EndPosition != 1 {
		t.Fail()
		t.Errorf("Expected error end position 1. Found position: %v", err.EndPosition)
	}

}

func TestErrorThenRecovery(t *testing.T) {
	reader := bytes.NewBufferString(
		`a := 1 * 5 + 7
b := 3.45 * 44.21 / (4; + a) 'A Simple Expression
double((a + b)`)
	bufreader := bufio.NewReader(reader)
	logfilename := "logs/TestErrorThenRecovery.log"
	logfile, fileErr := os.Create(logfilename)
	if fileErr != nil {
		t.Fatal("Error: Could not create log file " + logfilename + ": " + fileErr.Error())
	}
	_, errs := dsl.Parse(Parse, Scan, bufreader, dsl.WithLogger(logfile))

	if len(errs) != 2 {
		t.Fatalf("Should report exactly 2 errors: got %d", len(errs))
	}
	err := errs[0]
	if err.Code != dsl.ERROR_RUNE_EXPECTED_NOT_FOUND {
		t.Fail()
		t.Errorf("Expected error code 'Rune expected but not found'. Found error: '%v", err)
	}
	if err.StartLine != 2 {
		t.Fail()
		t.Errorf("Expected error line 2. Found line: %v", err.StartLine)
	}
	if err.StartPosition != 23 {
		t.Fail()
		t.Errorf("Expected error start position 23. Found position: %v", err.StartPosition)
	}
	if err.EndPosition != 23 {
		t.Fail()
		t.Errorf("Expected error end position 23. Found position: %v", err.EndPosition)
	}
	err = errs[1]
	if err.Code != dsl.ERROR_TOKEN_EXPECTED_NOT_FOUND {
		t.Fail()
		t.Errorf("Expected error code 'Token expected but not found'. Found error: '%v", err)
	}
	if err.StartLine != 3 {
		t.Fail()
		t.Errorf("Expected error line 3. Found line: %v", err.StartLine)
	}
	if err.StartPosition != 15 {
		t.Fail()
		t.Errorf("Expected error start position 15. Found position: %v", err.StartPosition)
	}
	if err.EndPosition != 16 {
		t.Fail()
		t.Errorf("Expected error end position 16. Found position: %v", err.EndPosition)
	}

}

func TestMultiLineError(t *testing.T) {
	reader := bytes.NewBufferString(
		`a := 1 * 5 + 7
b := 3.45 * 44.21 / 4" \ercec
gevhvrh  " + a) 'A Simple Expression
double(a + b)`)
	bufreader := bufio.NewReader(reader)
	logfilename := "logs/TestMultiLineError.log"
	logfile, fileErr := os.Create(logfilename)
	if fileErr != nil {
		t.Fatal("Error: Could not create log file " + logfilename + ": " + fileErr.Error())
	}
	_, errs := dsl.Parse(Parse, Scan, bufreader, dsl.WithLogger(logfile))

	if len(errs) != 1 {
		t.Fatalf("Should report exactly 1 error: got %d", len(errs))
	}
	err := errs[0]
	if err.Code != dsl.ERROR_TOKEN_EXPECTED_NOT_FOUND {
		t.Fail()
		t.Errorf("Expected error code 'Token expected but not found'. Found error: '%v", err)
	}
	if err.StartLine != 2 {
		t.Fail()
		t.Errorf("Expected error start line 3. Found line: %v", err.StartLine)
	}
	if err.StartPosition != 22 {
		t.Fail()
		t.Errorf("Expected error start position 22. Found position: %v", err.StartPosition)
	}
	if err.EndLine != 3 {
		t.Fail()
		t.Errorf("Expected error end line 1. Found line: %v", err.StartLine)
	}
	if err.EndPosition != 10 {
		t.Fail()
		t.Errorf("Expected error end position 10. Found position: %v", err.EndPosition)
	}

}
