package mydsl

import(
    "github.com/deslittle/go-dsl"
)
func Scan(s *dsl.Scanner) dsl.Token {
	s.Expect(dsl.ExpectRune{
		Branches: []dsl.Branch{
			{' ', nil},
			{'\t', nil}},
		Options: dsl.ScanOptions{Multiple: true, Optional: true, Skip: true}})
	s.Expect(dsl.ExpectRune{
		Branches: []dsl.Branch{
			{'-', nil},
			{'+', nil},
			{'*', nil},
			{'/', nil},
			{'(', nil},
			{')', nil},
			{'\n', nil},
			{':', assign},
			{'\'', comment},
			{rune(0), nil}},
		BranchRanges: []dsl.BranchRange{
			{'0', '9', literal},
			{'A', 'Z', variable},
			{'a', 'z', variable}}})
	s.Match(([]dsl.Match{
		{"-", "MINUS"},
		{"+", "PLUS"},
		{"*", "MULTIPLY"},
		{"/", "DIVIDE"},
		{"(", "OPEN_PAREN"},
		{")", "CLOSE_PAREN"},
		{"\n", "NL"}}))
	return s.Exit()
}

func variable(s *dsl.Scanner) {
	s.Expect(dsl.ExpectRune{
		Branches: []dsl.Branch{
			{'_', nil}},
		BranchRanges: []dsl.BranchRange{
			{'A', 'Z', nil},
			{'a', 'z', nil}},
		Options: dsl.ScanOptions{Multiple: true, Optional: true}})
	s.Match([]dsl.Match{{"", "VARIABLE"}})
}

// ScanFn -> literal
func literal(s *dsl.Scanner) {
	s.Expect(dsl.ExpectRune{
		BranchRanges: []dsl.BranchRange{
			{'0', '9', nil}},
		Options: dsl.ScanOptions{Multiple: true, Optional: true}})
	s.Expect(dsl.ExpectRune{
		Branches: []dsl.Branch{
			{'.', fraction}},
		Options: dsl.ScanOptions{Optional: true}})
	s.Match([]dsl.Match{{"", "LITERAL"}})
}

// ScanFn -> number -> fraction
func fraction(s *dsl.Scanner) {
	s.Expect(dsl.ExpectRune{
		BranchRanges: []dsl.BranchRange{
			{'0', '9', nil}},
		Options: dsl.ScanOptions{Multiple: true}})
	s.Match([]dsl.Match{{"", "LITERAL"}})
}

// ScanFn -> assign
func assign(s *dsl.Scanner) {
	s.Expect(dsl.ExpectRune{
		Branches: []dsl.Branch{
			{'=', nil}}})
	s.Match([]dsl.Match{{":=", "ASSIGN"}})
}

// ScanFn -> comment
func comment(s *dsl.Scanner) {
    s.SkipRune()
	s.Expect(dsl.ExpectRune{
		Branches: []dsl.Branch{
			{rune(0), nil},
			{'\n', nil}},
		Options: dsl.ScanOptions{Multiple: true, Invert: true, Optional: true}})
	s.Match([]dsl.Match{{"", "COMMENT"}})
}