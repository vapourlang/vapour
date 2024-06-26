package lexer

import (
	"strings"
	"unicode/utf8"
)

type itemType int

type item struct {
	class itemType
	val   string
}

type lexer struct {
	input string
	start int
	pos   int
	width int
	line  int
	items []item
}

const (
	itemError itemType = iota
	itemDoubleQuote
	itemSingleQuote
	itemBackslash
	itemCurlyLeft
	itemCurlyRight
	itemSquareLeft
	itemSquareRight
	itemName
	itemContent
)

const eof = -1

func (l *lexer) emit(t itemType) {
	// skip empty tokens
	if l.start == l.pos {
		return
	}

	l.items = append(l.items, item{t, l.input[l.start:l.pos]})
	l.start = l.pos
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}

	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width

	if r == '\n' {
		l.line++
	}

	return r
}

func (l *lexer) skipLine() {
	currentLine := l.line
	for {
		newLine := l.line

		if newLine > currentLine {
			break
		}

		l.next()
		l.ignore()
	}
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek(n int) rune {
	var r rune
	for i := 0; i < n; i++ {
		r = l.next()
	}

	for i := 0; i < n; i++ {
		l.backup()
	}
	return r
}

type stateFn func(*lexer) stateFn

func (l *lexer) run() {
	for state := lexDefault; state != nil; {
		state = state(l)
	}
}

func lexDefault(l *lexer) stateFn {
	r1 := l.peek(1)

	if r1 == eof {
		return nil
	}

	switch r1 {
	case '"':
		l.emit(itemDoubleQuote)
	case '\'':
		l.emit(itemSingleQuote)
	}

	return lexDefault
}

func (l *lexer) acceptSpace() bool {
	return l.accept(" \\t")
}

func (l *lexer) acceptAlpha() bool {
	return l.accept("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

func (l *lexer) acceptNumeric() bool {
	return l.accept("0123456789")
}

func (l *lexer) accept(rs string) bool {
	for strings.IndexRune(rs, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}
