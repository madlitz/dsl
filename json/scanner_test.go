package json

// import(
//     "github.com/madlitz/go-dsl"
// )
// func Scan(s *dsl.Scanner) dsl.Token {
// 	s.Expect(dsl.ExpectRune{
// 		Branches: []dsl.Branch{
// 			{' ', nil},
// 			{'\t', nil},
//             {'\n', nil},
//             {'\r', nil}},
// 		Options: dsl.ScanOptions{Multiple: true, Optional: true, Skip: true}})
//     s.Expect(dsl.ExpectRune{
//         Branches: []dsl.Branch{
// 			{'"', fnString},
//             {'-', fnInteger},
//             {':', nil},
//             {',', nil},
//             {'{', nil},
//             {'}', nil},
//             {'[', nil},
//             {']', nil},
//             {'0', fnFraction},
//             {rune(0), fnEOF}},
//         BranchRanges: []dsl.BranchRange{
// 			{'1', '9', fnInteger},
//             {'a', 'z', fnKeyword},
//     s.Match([]dsl.Match{
//         {"{", "OPEN_OBJECT"},
//         {"}", "CLOSE_OBJECT"},
//         {"[", "OPEN_ARRAY"},
//         {"]", "CLOSE_ARRAY"},
//         {":", "COLON"},
//         {",", "COMMA"},
//         {"", "ILLEGAL"}})
// 	return s.Exit()
// }

// // Scan -> fnKeyword
// func fnKeyword(s *dsl.Scanner) {
//     s.Expect(dsl.ExpectRune{
//         BranchRanges: []dsl.BranchRange{
// 			{'a', 'z', nil}},
// 		Options: dsl.ScanOptions{Multiple: true, Optional: true}})
//     s.Match([]dsl.Match{
//         {"true", "TRUE"},
//         {"false", "FALSE"},
//         {"null", "NULL"},
//         {"", "ILLEGAL"}})
// }

// // Scan -> fnInteger
// func fnInteger(s *dsl.Scanner) {
//     s.Expect(dsl.ExpectRune{
//         BranchRanges: []dsl.BranchRange{
// 			{'0', '9', nil}},
// 		Options: dsl.ScanOptions{Multiple: true, Optional: true}})
//     s.Call(fnFraction)
// }

// // Scan -> fnInteger -> fnFraction
// // Scan -> fnFraction
// func fnFraction(s *dsl.Scanner) {
//     s.Expect(dsl.ExpectRune{
//         Branches: []dsl.Branch{
// 			{'.', fnFractionalPart},
//         Options: dsl.ScanOptions{Optional: true}})
//     s.Expect(dsl.ExpectRune{
//         Branches: []dsl.Branch{
// 			{'e', fnExponent},
//             {'E', fnExponent},
//         Options: dsl.ScanOptions{Optional: true}})
//     s.Match([]dsl.Match{{"", "NUMBER"}})
// }

// // Scan -> fnInteger -> fnFraction -> fnFractionalPart
// func fnFractionalPart(s *dsl.Scanner) {
//     s.Expect(dsl.ExpectRune{
//         BranchRanges: []dsl.BranchRange{
// 			{'0', '9', nil}},
//         Options: dsl.ScanOptions{Multiple: true})
// }

// // Scan -> fnInteger, fnFraction -> fnExponent
// func fnExponent(s *dsl.Scanner) {
//     s.Expect(dsl.ExpectRune{
//         Branches: []dsl.BranchToken{
// 			{'+', nil},
//             {'-', nil}},
//         Options: dsl.ScanOptions{Optional: true})
//     s.Expect(dsl.ExpectRune{
//         BranchRanges: []dsl.BranchRange{
// 			{'0', '9', nil}},
//         Options: dsl.ScanOptions{Multiple: true})
// }

// // Scan -> fnString
// func fnString(s *dsl.Scanner) {
//     s.Expect(dsl.ExpectRune{
//         Branches: []dsl.Branch{
// 			{'\', fnControl},
//         BranchRanges: []dsl.BranchRange{
// 			{rune(32), rune(127), nil},
//             {rune(160), rune(0x7FFFFFFF), nil}},
//         Options: dsl.ScanOptions{Multiple: true, Optional: true}})
//     s.Expect(dsl.ExpectRune{
//         Branches: []dsl.Branch{
// 			{'"', nil}}})
//     s.Match([]dsl.Match{{"", "STRING"}})
// }

// // Scan -> fnString -> fnControl
// func fnControl(s *dsl.Scanner) {
//      s.Expect(dsl.ExpectRune{
//         Branches: []dsl.Branch{
// 			{'"', nil},
//             {'\', nil},
//             {'/', nil},
//             {'b', nil},
//             {'f', nil},
//             {'n', nil},
//             {'r', nil},
//             {'t', nil},
//             {'u', fnHex}}})
// }

// // Scan -> fnString -> fnHex
// func fnHex(s *dsl.Scanner) {
//     s.Expect(dsl.ExpectRune{
//         BranchRanges: []dsl.BranchRange{
// 			{'0', '9', nil},
//             {'A', 'F', nil},
//             {'a', 'f', nil}}})
//     s.Expect(dsl.ExpectRune{
//         BranchRanges: []dsl.BranchRange{
// 			{'0', '9', nil},
//             {'A', 'F', nil},
//             {'a', 'f', nil}}})
//     s.Expect(dsl.ExpectRune{
//         BranchRanges: []dsl.BranchRange{
// 			{'0', '9', nil},
//             {'A', 'F', nil},
//             {'a', 'f', nil}}})
//     s.Expect(dsl.ExpectRune{
//         BranchRanges: []dsl.BranchRange{
// 			{'0', '9', nil},
//             {'A', 'F', nil},
//             {'a', 'f', nil}}})
// }

// // Scan -> fnEOF
// func fnEOF(s *dsl.Scanner) {
//     s.Match([]dsl.Match{{"", "EOF"}})
// }
