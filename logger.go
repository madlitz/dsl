package dsl

import "log"

type logger interface {
	log(msg string, indent indent)
}

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

type dslNoLogger struct{}

func (l *dslNoLogger) log(msg string, indent indent) {}
