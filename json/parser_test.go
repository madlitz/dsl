package json

import(
    "github.com/deslittle/go-dsl"
)
func Parse(p *dsl.Parser) (dsl.AST, []dsl.Error) {
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"OPEN_OBJECT", nil},
		Options: dsl.ParseOptions{Skip: true}})
	p.Peek(dsl.PeekToken{
		Branches: []dsl.PeekToken{
			[]string{"STRING"}, fnObject}})
	p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"CLOSE_OBJECT", nil},
		Options: dsl.ParseOptions{Skip: true}})
    p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"EOF", nil},
		Options: dsl.ParseOptions{Skip: true}})
	return p.Exit()
}

// parse -> assignmentOrCall
func fnObject(p *dsl.Parser) {
    p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"STRING", fnMemberNoKey},
		Options: dsl.ParseOptions{Optional: true}}})
    p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"CLOSE_OBJECT", nil},
		Options: dsl.ParseOptions{Skip: true}}})
}

func fnMemberNoKey(p *dsl.Parser) {
    p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{":", nil},
		Options: dsl.ParseOptions{Skip: true}})
    p.Expect(dsl.ExpectToken{    
        Branches: []dsl.BranchToken{
			{"STRING", nil},
            {"NUMBER", nil},
            {"OPEN_OBJECT", fnObject},
            {"OPEN_ARRAY", fnArray},
            {"TRUE", nil},
            {"FALSE", nil},
            {"NULL", nil}},
		Options: dsl.ParseOptions{Multiple: true}})
    p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{",", fnMember},
		Options: dsl.ParseOptions{Multiple: true, Optional: true, Skip: true}})
}

func fnMember(p *dsl.Parser) {
    p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{"STRING", fnMember}}})
    p.Expect(dsl.ExpectToken{
		Branches: []dsl.BranchToken{
			{":", nil},
		Options: dsl.ParseOptions{Skip: true}})
    p.Expect(dsl.ExpectToken{    
        Branches: []dsl.BranchToken{
			{"STRING", nil},
            {"NUMBER", nil},
            {"OPEN_OBJECT", fnObject},
            {"OPEN_ARRAY", fnArray},
            {"TRUE", nil},
            {"FALSE", nil},
            {"NULL", nil}}})
}

func fnArray(p *dsl.Parser) {
	p.Expect(dsl.ExpectToken{    
        Branches: []dsl.BranchToken{
			{"STRING", nil},
            {"NUMBER", nil},
            {"OPEN_OBJECT", fnObject},
            {"OPEN_ARRAY", fnArray},
            {"TRUE", nil},
            {"FALSE", nil},
            {"NULL", nil},
		Options: dsl.ParseOptions{Optional: true}})
    p.Expect(dsl.ExpectToken{    
        Branches: []dsl.BranchToken{
			{"COMMA", fnArrayValue},
		Options: dsl.ParseOptions{Multiple: true, Optional: true}})
    p.Expect(dsl.ExpectToken{    
        Branches: []dsl.BranchToken{
			{"CLOSE_ARRAY", nil}}})
}

// parse -> assignmentOrCall -> [expression] -> addcomment
func fnArrayValue(p *dsl.Parser) {
	p.Expect(dsl.ExpectToken{    
        Branches: []dsl.BranchToken{
			{"STRING", nil},
            {"NUMBER", nil},
            {"OPEN_OBJECT", fnObject},
            {"OPEN_ARRAY", fnArray},
            {"TRUE", nil},
            {"FALSE", nil},
            {"NULL", nil}}})
}


func fnObject(p *dsl.Parser) {
	p.AddNode("COMMENT")
	p.AddTokens()
	p.WalkUp()
}