package json

import (
	"bufio"
	"bytes"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/madlitz/go-dsl"
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

	expected := []dsl.Node{
		{
			Type: NODE_OBJECT,
			Children: []dsl.Node{
				{
					Type:   NODE_MEMBER,
					Tokens: []dsl.Token{{ID: TOKEN_STRING, Literal: "key1", Line: 2, Position: 2}},
					Children: []dsl.Node{
						{Type: NODE_VALUE, Tokens: []dsl.Token{{ID: TOKEN_STRING, Literal: "value1", Line: 2, Position: 10}}},
					},
				},
				{
					Type:   NODE_MEMBER,
					Tokens: []dsl.Token{{ID: TOKEN_STRING, Literal: "key2", Line: 3, Position: 3}},
					Children: []dsl.Node{
						{Type: NODE_VALUE, Tokens: []dsl.Token{{ID: TOKEN_NUMBER, Literal: "42", Line: 3, Position: 10}}},
					},
				},
				{
					Type:   NODE_MEMBER,
					Tokens: []dsl.Token{{ID: TOKEN_STRING, Literal: "key3", Line: 4, Position: 2}},
					Children: []dsl.Node{
						{Type: NODE_VALUE, Tokens: []dsl.Token{{ID: TOKEN_TRUE, Literal: "true", Line: 4, Position: 9}}},
					},
				},
				{
					Type:   NODE_MEMBER,
					Tokens: []dsl.Token{{ID: TOKEN_STRING, Literal: "key4", Line: 5, Position: 3}},
					Children: []dsl.Node{
						{Type: NODE_VALUE, Tokens: []dsl.Token{{ID: TOKEN_NULL, Literal: "null", Line: 5, Position: 10}}},
					},
				},
				{
					Type:   NODE_MEMBER,
					Tokens: []dsl.Token{{ID: TOKEN_STRING, Literal: "key5", Line: 6, Position: 3}},
					Children: []dsl.Node{
						{
							Type: NODE_OBJECT,
							Children: []dsl.Node{
								{
									Type:   NODE_MEMBER,
									Tokens: []dsl.Token{{ID: TOKEN_STRING, Literal: "nestedKey", Line: 7, Position: 4}},
									Children: []dsl.Node{
										{Type: NODE_VALUE, Tokens: []dsl.Token{{ID: TOKEN_STRING, Literal: "nestedValue", Line: 7, Position: 17}}},
									},
								},
							},
						},
					},
				},
				{
					Type:   NODE_MEMBER,
					Tokens: []dsl.Token{{ID: TOKEN_STRING, Literal: "key6", Line: 9, Position: 3}},
					Children: []dsl.Node{
						{
							Type: NODE_ARRAY,
							Children: []dsl.Node{
								{Type: NODE_VALUE, Tokens: []dsl.Token{{ID: TOKEN_NUMBER, Literal: "1", Line: 9, Position: 11}}},
								{Type: NODE_VALUE, Tokens: []dsl.Token{{ID: TOKEN_NUMBER, Literal: "2", Line: 9, Position: 14}}},
								{Type: NODE_VALUE, Tokens: []dsl.Token{{ID: TOKEN_NUMBER, Literal: "3", Line: 9, Position: 17}}},
								{Type: NODE_VALUE, Tokens: []dsl.Token{{ID: TOKEN_STRING, Literal: "four", Line: 9, Position: 21}}},
							},
						},
					},
				},
			},
		},
	}

	if diff := cmp.Diff(expected, ast.RootNode.Children, cmpopts.IgnoreFields(dsl.Node{}, "Parent")); diff != "" {
		t.Errorf("AST mismatch (-got +want):\n%s", diff)
	}

}
