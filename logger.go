package dsl

import (
	"bytes"
	"log"
	"reflect"
	"runtime"
)

type logger interface {
	log(msg string, indent indent)
}

// dslLogger is a simple logger that uses the standard log package to print
// messages.
type dslLogger struct {
	indent int
	logger *log.Logger
	buf    []interface{}
}

type indent int

const (
	prefixNone indent = iota // 0
	prefixNewline
	prefixIncrement
	prefixDecrement
	prefixStartLine
	prefixError
)

// log logs the message to the io.Writer. It keeps a buffer of messages to print.
// The buffer is printed when a new line is requested as the std logger always
// prints a newline.
// The logger also allows the user to manage the indentation level of the log.
func (l *dslLogger) log(msg string, indent indent) {
	switch indent {
	case prefixNone:
	case prefixNewline:
		l.printLogBuffer()
	case prefixIncrement:
		l.printLogBuffer()
		l.indent++
		prefix := ""
		for i := 0; i < l.indent; i++ {
			prefix += "\t"
		}
		l.logger.SetPrefix(prefix)
	case prefixDecrement:
		l.printLogBuffer()
		l.indent--
		prefix := ""
		for i := 0; i < l.indent; i++ {
			prefix += "\t"
		}
		l.logger.Print(msg)
		l.logger.SetPrefix(prefix)
		return
	case prefixStartLine:
		l.printLogBuffer()
		prefix := l.logger.Prefix()
		l.logger.SetPrefix("")
		l.logger.Print(msg)
		l.logger.SetPrefix(prefix)
		return
	case prefixError:
		l.printLogBuffer()
		prefix := l.logger.Prefix()
		l.logger.SetPrefix("***")
		l.logger.Print(msg)
		l.logger.SetPrefix(prefix)
		return
	}
	l.buf = append(l.buf, msg)
}

func (l *dslLogger) printLogBuffer() {
	if l.buf != nil {
		l.logger.Print(l.buf...)
		l.buf = nil
	}
}

// dslNoLogger is a logger that does nothing. It is used when the user does not
// want to log anything.
type dslNoLogger struct{}

func (l *dslNoLogger) log(msg string, indent indent) {}

// ---------------------------------------------------------------------------------------------------------

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
