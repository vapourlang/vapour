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
	itemIdent
	itemDoubleQuote
	itemSingleQuote
	itemAssign
	itemLeftCurly
	itemRightCurly
	itemLeftParen
	itemRightParen
	itemString
	itemInteger
	itemFloat
	itemNamespace
	itemMathOperation
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

	if r1 == '"' {
		l.next()
		l.emit(itemDoubleQuote)
		return lexString
	}

	if r1 == '\'' {
		l.next()
		l.emit(itemSingleQuote)
		return lexString
	}

	// we parsed strings: we skip spaces
	if r1 == ' ' {
		l.ignore()
		return lexDefault
	}

	// peek one more rune
	r2 := l.peek(2)

	if r1 == '<' && r2 == '-' {
		l.next()
		l.next()
		l.emit(itemAssign)
		return lexDefault
	}

	if r1 == ':' && r2 == ':' {
		l.next()
		l.next()
		l.emit(itemNamespace)
		return lexIdentifier
	}

	if r1 == '=' {
		l.next()
		l.emit(itemAssign)
		return lexDefault
	}

	if r1 == '(' {
		l.next()
		l.emit(itemLeftParen)
		return lexDefault
	}

	if r1 == ')' {
		l.next()
		l.emit(itemLeftParen)
		return lexDefault
	}

	if r1 == '{' {
		l.next()
		l.emit(itemLeftCurly)
		return lexDefault
	}

	if r1 == '}' {
		l.next()
		l.emit(itemRightCurly)
		return lexDefault
	}

	if l.acceptNumber() {
		l.emit(itemInteger)
		return lexDefault
	}

	if l.acceptFloat() {
		l.emit(itemFloat)
		return lexDefault
	}

	if l.acceptAlphaNumeric() {
		l.emit(itemIdent)
		return lexDefault
	}

	if l.acceptMathOp() {
		l.emit(itemMathOperation)
		return lexDefault
	}

	l.next()
	return lexDefault
}

func lexString(l *lexer) stateFn {
	r := l.peek(1)
	for r != '"' && r != '\'' {
		l.next()
		r = l.peek(1)
	}
	l.emit(itemString)
	return lexDefault
}

func lexIdentifier(l *lexer) stateFn {
	l.acceptAlphaNumeric()
	return lexDefault
}

func (l *lexer) acceptSpace() bool {
	return l.accept(" \\t")
}

func (l *lexer) acceptAlpha() bool {
	return l.accept("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

func (l *lexer) acceptNumber() bool {
	return l.accept("0123456789")
}

func (l *lexer) acceptFloat() bool {
	return l.accept("0123456789.")
}

func (l *lexer) acceptMathOp() bool {
	return l.accept("+\\-*")
}

func (l *lexer) acceptAlphaNumeric() bool {
	return l.accept("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
}

func (l *lexer) accept(rs string) bool {
	for strings.IndexRune(rs, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}
