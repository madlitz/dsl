package dsl

import (
	"bytes"
	"reflect"
	"runtime"
)

// Uses reflection package to pull out the name of the currently executing user
// scan function or user parse function. Useful for debugging.
func getFuncName(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

// When printing scanned runes in a log or error it is useful to represent
// non text characters as their text equivalents.
func sanitize(str string, whitespace bool) string {
	var buf bytes.Buffer

	if whitespace && str == " " {
		return "WS" // Only if the string is a single whitespace character
	}
	for _, rn := range str {
		switch rn {
		case '\n':
			buf.WriteString("NL") // New Line
		case '\r':
			buf.WriteString("CR") // Carriage Return
		case '\t':
			buf.WriteString("TAB") // Tab
		case '\v':
			buf.WriteString("VTAB") // Vertical Tab
		case '\a':
			buf.WriteString("BELL") // Bell
		case rune(0):
			buf.WriteString("EOF") // End of File
		default:
			buf.WriteRune(rn)
		}
	}
	return buf.String()
}
