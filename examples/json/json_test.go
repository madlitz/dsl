package json_test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/madlitz/dsl"

	. "github.com/madlitz/dsl/examples/json"
)

func TestJSONParser(t *testing.T) {
	reader := bytes.NewBufferString(`{
"key1": "value1",
	"key2": 42,
"key3": true,
	"key4": null,
	"key5": {
		"nestedKey": "nestedValue"
	},
	"key6": [1, 2, 3, "four"]
}`)
	bufreader := bufio.NewReader(reader)
	logfilename := "logs/json.log"
	logfile, err := os.Create(logfilename)
	if err != nil {
		t.Fatal(err)
	}

	ast, errs := dsl.Parse(Parse, Scan, bufreader, dsl.WithLogger(logfile))
	if len(errs) > 0 {
		t.Fail()
		for _, err := range errs {
			t.Error(err.Error())
		}
	}

	astJSON, _ := json.Marshal(ast.RootNode)

	expectedJSON := []byte(`{
		"type": "ROOT",
		"tokens": null,
		"children": [
			{
				"type": "OBJECT",
				"tokens": null,
				"children": [
					{
						"type": "MEMBER",
						"tokens": [
							{
								"ID": "STRING",
								"Literal": "key1"
							}
						],
						"children": [
							{
								"type": "VALUE",
								"tokens": [
									{
										"ID": "STRING",
										"Literal": "value1"
									}
								],
								"children": null
							}
						]
					},
					{
						"type": "MEMBER",
						"tokens": [
							{
								"ID": "STRING",
								"Literal": "key2"
							}
						],
						"children": [
							{
								"type": "VALUE",
								"tokens": [
									{
										"ID": "NUMBER",
										"Literal": "42"
									}
								],
								"children": null
							}
						]
					},
					{
						"type": "MEMBER",
						"tokens": [
							{
								"ID": "STRING",
								"Literal": "key3"
							}
						],
						"children": [
							{
								"type": "VALUE",
								"tokens": [
									{
										"ID": "TRUE",
										"Literal": "true"
									}
								],
								"children": null
							}
						]
					},
					{
						"type": "MEMBER",
						"tokens": [
							{
								"ID": "STRING",
								"Literal": "key4"
							}
						],
						"children": [
							{
								"type": "VALUE",
								"tokens": [
									{
										"ID": "NULL",
										"Literal": "null"
									}
								],
								"children": null
							}
						]
					},
					{
						"type": "MEMBER",
						"tokens": [
							{
								"ID": "STRING",
								"Literal": "key5"
							}
						],
						"children": [
							{
								"type": "OBJECT",
								"tokens": null,
								"children": [
									{
										"type": "MEMBER",
										"tokens": [
											{
												"ID": "STRING",
												"Literal": "nestedKey"
											}
										],
										"children": [
											{
												"type": "VALUE",
												"tokens": [
													{
														"ID": "STRING",
														"Literal": "nestedValue"
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
						"type": "MEMBER",
						"tokens": [
							{
								"ID": "STRING",
								"Literal": "key6"
							}
						],
						"children": [
							{
								"type": "ARRAY",
								"tokens": null,
								"children": [
									{
										"type": "VALUE",
										"tokens": [
											{
												"ID": "NUMBER",
												"Literal": "1"
											}
										],
										"children": null
									},
									{
										"type": "VALUE",
										"tokens": [
											{
												"ID": "NUMBER",
												"Literal": "2"
											}
										],
										"children": null
									},
									{
										"type": "VALUE",
										"tokens": [
											{
												"ID": "NUMBER",
												"Literal": "3"
											}
										],
										"children": null
									},
									{
										"type": "VALUE",
										"tokens": [
											{
												"ID": "STRING",
												"Literal": "four"
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
	  	]
	}`)

	fmt.Println(string(astJSON))

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
