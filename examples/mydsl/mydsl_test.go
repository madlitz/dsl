package mydsl_test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/madlitz/dsl"
	. "github.com/madlitz/dsl/examples/mydsl"
)

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

	astJSON, _ := json.Marshal(ast.RootNode)

	expectedJSON := []byte(`
	{
		"type": "ROOT",
		"tokens": null,
		"children": [
			{
				"type": "ASSIGNMENT",
				"tokens": [
					{"ID": "VARIABLE", "Literal": "a", "Line": 1, "Position": 1}
				],
				"children": [
					{
						"type": "TERMINAL",
						"tokens": [
							{"ID": "LITERAL", "Literal": "1", "Line": 1, "Position": 6}
						],
						"children": null
					},
					{
						"type": "EXPRESSION",
						"tokens": [
							{"ID": "MULTIPLY", "Literal": "*", "Line": 1, "Position": 8}
						],
						"children": [
							{
								"type": "TERMINAL",
								"tokens": [
									{"ID": "LITERAL", "Literal": "5", "Line": 1, "Position": 10}
								],
								"children": null
							},
							{
								"type": "EXPRESSION",
								"tokens": [
									{"ID": "PLUS", "Literal": "+", "Line": 1, "Position": 12}
								],
								"children": [
									{
										"type": "TERMINAL",
										"tokens": [
											{"ID": "LITERAL", "Literal": "7", "Line": 1, "Position": 14}
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
					{"ID": "VARIABLE", "Literal": "b", "Line": 2, "Position": 1}
				],
				"children": [
					{
						"type": "TERMINAL",
						"tokens": [
							{"ID": "LITERAL", "Literal": "3.45", "Line": 2, "Position": 6}
						],
						"children": null
					},
					{
						"type": "EXPRESSION",
						"tokens": [
							{"ID": "MULTIPLY", "Literal": "*", "Line": 2, "Position": 11}
						],
						"children": [
							{
								"type": "TERMINAL",
								"tokens": [
									{"ID": "LITERAL", "Literal": "44.21", "Line": 2, "Position": 13}
								],
								"children": null
							},
							{
								"type": "EXPRESSION",
								"tokens": [
									{"ID": "DIVIDE", "Literal": "/", "Line": 2, "Position": 19}
								],
								"children": [
									{
										"type": "EXPRESSION",
										"tokens": [
											{"ID": "OPEN_PAREN", "Literal": "(", "Line": 2, "Position": 21}
										],
										"children": [
											{
												"type": "TERMINAL",
												"tokens": [
													{"ID": "LITERAL", "Literal": "4", "Line": 2, "Position": 22}
												],
												"children": null
											},
											{
												"type": "EXPRESSION",
												"tokens": [
													{"ID": "PLUS", "Literal": "+", "Line": 2, "Position": 24}
												],
												"children": [
													{
														"type": "TERMINAL",
														"tokens": [
															{"ID": "VARIABLE", "Literal": "a", "Line": 2, "Position": 26}
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
											{"ID": "CLOSE_PAREN", "Literal": ")", "Line": 2, "Position": 27}
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
					{"ID": "COMMENT", "Literal": "A Simple Expression", "Line": 2, "Position": 30}
				],
				"children": null
			},
			{
				"type": "CALL",
				"tokens": [
					{"ID": "VARIABLE", "Literal": "double", "Line": 3, "Position": 1}
				],
				"children": [
					{
						"type": "TERMINAL",
						"tokens": [
							{"ID": "VARIABLE", "Literal": "a", "Line": 3, "Position": 8}
						],
						"children": null
					},
					{
						"type": "EXPRESSION",
						"tokens": [
							{"ID": "PLUS", "Literal": "+", "Line": 3, "Position": 10}
						],
						"children": [
							{
								"type": "TERMINAL",
								"tokens": [
									{"ID": "VARIABLE", "Literal": "b", "Line": 3, "Position": 12}
								],
								"children": null
							}
						]
					}
				]
			}
		]
	}`)

	expectJSON(t, expectedJSON, astJSON)

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
	if err.Code != dsl.ErrorTokenExpectedNotFound {
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
	if err.EndPosition != 7 {
		t.Fail()
		t.Errorf("Expected error end position 7. Found position: %v", err.EndPosition)
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
	if err.Code != dsl.ErrorRuneExpectedNotFound {
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
	if err.Code != dsl.ErrorRuneExpectedNotFound {
		t.Fail()
		t.Log(err.Error())
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
	if err.Code != dsl.ErrorTokenExpectedNotFound {
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
	if err.EndPosition != 15 {
		t.Fail()
		t.Errorf("Expected error end position 15. Found position: %v", err.EndPosition)
	}

}

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

	astJSON, _ := json.Marshal(ast)
	expectedJSON := []byte(`{
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
	}`)

	expectJSON(t, expectedJSON, astJSON)

}

// expectJSON returns an assertion function that compares the expected and
// actual JSON payloads.
func expectJSON(t *testing.T, expected []byte, actual []byte) {

	t.Helper()

	var a, e map[string]any
	if err := json.Unmarshal(expected, &e); err != nil {
		t.Fatalf("error unmarshaling expected json payload: %v", err)
	}

	if err := json.Unmarshal(actual, &a); err != nil {
		t.Fatalf("error unmarshaling actual json payload: %v", err)
	}

	if diff := cmp.Diff(e, a); diff != "" {
		t.Errorf(diff)
	}

}
