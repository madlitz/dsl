package mydsl

import(
    "github.com/Autoblocks/go-dsl"
)
var recover bool

func Parse(p *dsl.Parser) (dsl.AST, []dsl.Error) {
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"VARIABLE", assignmentOrCall},
			{"EOF", nil}},
		Options: dsl.ParseOptions{Multiple: true}})

	return p.Exit()
}

// parse -> assignmentOrCall
func assignmentOrCall(p *dsl.Parser) {
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"ASSIGN", assignment},
			{"OPEN_PAREN", call}}})
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"COMMENT", addcomment}},
		Options: dsl.ParseOptions{Optional: true}})
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"NL", nil},
			{"EOF", nil}},
		Options: dsl.ParseOptions{Skip: true}})
}

// parse -> assignmentOrCall -> assignment
func assignment(p *dsl.Parser) {
	p.SkipToken()
	p.AddNode("ASSIGNMENT")
	p.AddTokens()
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"VARIABLE", operator},
			{"LITERAL", operator},
			{"OPEN_PAREN", paren_expression}}})
}

// parse -> assignmentOrCall -> call
func call(p *dsl.Parser) {
	p.SkipToken()
	p.AddNode("CALL")
	p.AddTokens()
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"VARIABLE", operator},
			{"LITERAL", operator},
			{"OPEN_PAREN", paren_expression}}})
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"CLOSE_PAREN", closecall}}})
}

// parse -> assignmentOrCall -> call -> [expression] -> closecall
func closecall(p *dsl.Parser) {
	p.SkipToken()
	p.WalkUp()
}

// parse -> assignmentOrCall -> assignment -> [operator, expression]
// parse -> assignmentOrCall -> call -> [operator, expression]
func expression(p *dsl.Parser) {
	p.AddNode("EXPRESSION")
	p.AddTokens()
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"VARIABLE", operator},
			{"LITERAL", operator},
			{"OPEN_PAREN", paren_expression},
			{"CLOSE_PAREN", operator}}})

}

// parse -> assignmentOrCall -> assignment -> [expression, operator]
// parse -> assignmentOrCall -> call -> [expression, operator]
func operator(p *dsl.Parser) {
	p.AddNode("TERMINAL")
	p.AddTokens()
	p.WalkUp()
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"PLUS", expression},
			{"MINUS", expression},
			{"DIVIDE", expression},
			{"MULTIPLY", expression},
			{"OPEN_PAREN", paren_expression}},
		Options: dsl.ParseOptions{Optional: true}})
	p.WalkUp()
}

// parse -> assignmentOrCall -> assignment -> [expression, operator] -> paren_expression
// parse -> assignmentOrCall -> assignment -> paren_expression
// parse -> assignmentOrCall -> call -> [expression, operator] -> paren_expression
// parse -> assignmentOrCall -> call -> paren_expression
func paren_expression(p *dsl.Parser) {
	p.Peek([]dsl.PeekToken{
		{[]string{}, expression}})
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"CLOSE_PAREN", operator}}})
	p.Recover(skipUntilLineBreak)
}

// parse -> assignmentOrCall -> [expression] -> addcomment
func addcomment(p *dsl.Parser) {
	p.AddNode("COMMENT")
	p.AddTokens()
	p.WalkUp()
}

func skipUntilLineBreak(p *dsl.Parser) {
	recover = true
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"UNKNOWN", nil}}})
	recover = false
}