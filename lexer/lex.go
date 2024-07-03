package lexer

import (
	"fmt"
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

	// identifiers
	itemIdent

	// quotes
	itemDoubleQuote
	itemSingleQuote

	// dollar $ign
	itemDollar

	// backtick
	itemBacktick

	// infix %>%
	itemInfix

	// comma,
	itemComma

	// question mark?
	itemQuestion

	// boolean
	itemBool

	// boolean
	itemReturn

	// ...
	itemThreeDot

	// native pipe
	itemPipe

	// = <-
	itemAssign

	// .Call
	itemCCall

	// NULL
	itemNULL

	// NA
	itemNA
	itemNACharacter
	itemNAReal
	itemNAComplex
	itemNAInteger

	// parens and brackets
	itemLeftCurly
	itemRightCurly
	itemLeftParen
	itemRightParen
	itemLeftSquare
	itemRightSquare
	itemDoubleLeftSquare
	itemDoubleRightSquare

	// "strings"
	itemString

	// numbers
	itemInteger
	itemFloat

	// namespace::
	itemNamespace
	// namespace:::
	itemNamespaceInternal

	// colon
	itemColon

	// + - / * ^
	itemMathOperation
	itemModulus

	// comment
	itemComment

	// roxygen comments
	itemSpecialComment
	itemRoxygenTagAt
	itemRoxygenTag
	itemRoxygenTagContent

	// doctor tags
	itemTypeDef
	itemTypeVar

	// compare
	itemDoubleEqual
	itemLessThan
	itemGreaterThan
	itemNotEqual
	itemLessOrEqual
	itemGreaterOrEqual

	// if else
	itemIf
	itemElse
	itemAnd
	itemOr
)

const stringNumber = "0123456789"
const stringAlpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const stringAlphaNum = stringAlpha + stringNumber
const stringMathOp = "+\\-*^"

const eof = -1

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items = append(l.items, item{itemError, fmt.Sprintf(format, args...)})
	return nil
}

func (l *lexer) emit(t itemType) {
	// skip empty tokens
	if l.start == l.pos {
		return
	}

	l.items = append(l.items, item{t, l.input[l.start:l.pos]})
	l.start = l.pos
}

// returns currently accepted token
func (l *lexer) token() string {
	return l.input[l.start:l.pos]
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
		return l.lexString('"')
	}

	if r1 == '\'' {
		l.next()
		l.emit(itemSingleQuote)
		return l.lexString('\'')
	}

	if r1 == '#' {
		return lexComment
	}

	// we parsed strings: we skip spaces and new lines
	if r1 == ' ' || r1 == '\t' || r1 == '\n' {
		l.next()
		l.ignore()
		return lexDefault
	}

	// peek one more rune
	r2 := l.peek(2)

	// if it's not %% it's an infix
	if r1 == '%' && r2 != '%' {
		return lexInfix
	}

	// it's a modulus
	if r1 == '%' && r2 == '%' {
		l.next()
		l.next()
		l.emit(itemModulus)
		return lexDefault
	}

	if r1 == '=' && r2 == '=' {
		l.next()
		l.next()
		l.emit(itemDoubleEqual)
		return lexDefault
	}

	if r1 == '!' && r2 == '=' {
		l.next()
		l.next()
		l.emit(itemNotEqual)
		return lexDefault
	}

	if r1 == '>' && r2 == '=' {
		l.next()
		l.next()
		l.emit(itemGreaterOrEqual)
		return lexDefault
	}

	if r1 == '<' && r2 == '=' {
		l.next()
		l.next()
		l.emit(itemLessOrEqual)
		return lexDefault
	}

	if r1 == '<' && r2 == ' ' {
		l.next()
		l.emit(itemLessThan)
		return lexDefault
	}

	if r1 == '>' && r2 == ' ' {
		l.next()
		l.emit(itemGreaterThan)
		return lexDefault
	}

	if r1 == '<' && r2 == '-' {
		l.next()
		l.next()
		l.emit(itemAssign)
		return lexDefault
	}

	if r1 == ':' && r2 == ':' && l.peek(3) == ':' {
		l.next()
		l.next()
		l.next()
		l.emit(itemNamespaceInternal)
		return lexIdentifier
	}

	if r1 == ':' && r2 == ':' {
		l.next()
		l.next()
		l.emit(itemNamespace)
		return lexIdentifier
	}

	if r1 == '.' && r2 == '.' && l.peek(3) == '.' {
		l.next()
		l.next()
		l.next()
		l.emit(itemThreeDot)
		return lexDefault
	}

	// we also emit namespace:: (above)
	// so we can assume this is not
	if r1 == ':' {
		l.next()
		l.emit(itemColon)
		return lexDefault
	}

	if r1 == '&' {
		l.next()
		l.emit(itemAnd)
		return lexDefault
	}

	if r1 == '|' && r2 == '>' {
		l.next()
		l.next()
		l.emit(itemPipe)
		return lexDefault
	}

	if r1 == '|' {
		l.next()
		l.emit(itemOr)
		return lexDefault
	}

	if r1 == '$' {
		l.next()
		l.emit(itemDollar)
		return lexDefault
	}

	if r1 == ',' {
		l.next()
		l.emit(itemComma)
		return lexDefault
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

	if r1 == '[' && r2 == '[' {
		l.next()
		l.emit(itemDoubleLeftSquare)
		return lexDefault
	}

	if r1 == '[' {
		l.next()
		l.emit(itemLeftSquare)
		return lexDefault
	}

	if r1 == ']' && r2 == ']' {
		l.next()
		l.emit(itemDoubleRightSquare)
		return lexDefault
	}

	if r1 == ']' {
		l.next()
		l.emit(itemRightSquare)
		return lexDefault
	}

	if r1 == '?' {
		l.next()
		l.emit(itemQuestion)
		return lexDefault
	}

	if r1 == '`' {
		l.next()
		l.emit(itemBacktick)
		return lexDefault
	}

	if l.acceptNumber() {
		return lexNumber
	}

	if l.acceptMathOp() {
		return lexMathOp
	}

	if l.acceptAlphaNumeric() {
		return lexIdentifier
	}

	l.next()
	return lexDefault
}

func lexMathOp(l *lexer) stateFn {
	l.acceptRun(stringMathOp)
	l.emit(itemMathOperation)
	return lexDefault
}

func lexNumber(l *lexer) stateFn {
	l.acceptRun(stringNumber)

	r := l.peek(1)

	if r == 'e' {
		l.next()
		l.acceptRun(stringNumber)
	}

	if l.accept(".") {
		l.acceptRun(stringNumber)
		l.emit(itemFloat)
		return lexDefault
	}

	l.emit(itemInteger)
	return lexDefault
}

func lexComment(l *lexer) stateFn {
	r2 := l.peek(2)

	if r2 == '\'' {
		l.next() // #
		l.next() // '

		l.emit(itemSpecialComment)
		return lexSpecialComment
	}

	r := l.peek(1)
	for r != '\n' && r != eof {
		l.next()
		r = l.peek(1)
	}

	l.emit(itemComment)

	return lexDefault
}

func lexSpecialComment(l *lexer) stateFn {
	r := l.peek(1)
	r2 := l.peek(2)

	// not entirely certain we need
	// #'[space], e.g.: #' @param
	// @#', e.g.: #'@param
	// perhaps legal too
	if r == ' ' {
		l.next()
		l.ignore()
	}

	if r == '@' || r2 == '@' {
		l.next()
		l.emit(itemRoxygenTagAt)
		return lexRoxygen
	}

	for r != '\n' && r != eof {
		l.next()
		r = l.peek(1)
	}

	l.emit(itemSpecialComment)

	return lexDefault
}

func lexRoxygen(l *lexer) stateFn {
	r := l.peek(1)
	for r != ' ' && r != '\t' && r != '\n' && r != eof {
		l.next()
		r = l.peek(1)
	}

	token := l.token()

	l.emit(itemRoxygenTag)

	if token == "type" {
		return lexTypeTag
	}

	if token == "yield" {
		return lexTypes
	}

	return lexRoxygenTagContent
}

func lexRoxygenTagContent(l *lexer) stateFn {
	r := l.peek(1)

	// we ignore space
	// e.g.: @param x Definition
	// skip space between x and Definition
	if r == ' ' {
		l.next()
		l.ignore()
	}

	for r != '\n' && r != eof {
		l.next()
		r = l.peek(1)
	}

	l.emit(itemRoxygenTagContent)

	return lexDefault
}

func lexTypeTag(l *lexer) stateFn {
	r := l.peek(1)
	for r != ':' && r != '\n' && r != eof {
		l.next()
		r = l.peek(1)
	}

	if r != ':' {
		l.next()
		return l.errorf("expects `:`, found %v [@type variable: type]", l.token())
	}

	l.emit(itemTypeVar)

	// ignore colon
	// e.g.: @type x: numeric
	l.next()
	l.ignore()

	return lexTypes
}

func lexTypes(l *lexer) stateFn {
	r := l.peek(1)

	if r == eof {
		return nil
	}

	if r == ' ' {
		l.next()
		l.ignore()
	}

	if r == '|' {
		l.next()
		l.ignore()
	}

	if r == '\n' {
		l.next()
		l.ignore()
		return lexDefault
	}

	r = l.peek(1)
	for r != '|' && r != ' ' && r != '\n' && r != eof {
		l.next()
		r = l.peek(1)
	}

	l.emit(itemTypeDef)

	return lexTypes
}

func (l *lexer) lexString(closing rune) func(l *lexer) stateFn {
	return func(l *lexer) stateFn {
		var c rune
		r := l.peek(1)
		for r != closing && r != eof {
			c = l.next()
			r = l.peek(1)
		}

		// this means the closing is escaped so
		// it's not in fact closing:
		// we move the cursor and keep parsing string
		// e.g.: "hello \"world\""
		if c == '\\' && r == closing {
			l.next()
			return l.lexString(closing)
		}

		if r == eof {
			l.next()
			return l.errorf("expecting closing quote, got %v", l.token())
		}

		l.emit(itemString)

		r = l.next()

		if r == '"' {
			l.emit(itemDoubleQuote)
		}

		if r == '\'' {
			l.emit(itemSingleQuote)
		}

		return lexDefault
	}
}

func lexInfix(l *lexer) stateFn {
	l.next()
	r := l.peek(1)
	for r != '%' && r != eof {
		l.next()
		r = l.peek(1)
	}

	if r == eof {
		l.next()
		return l.errorf("expecting closing %%, got %v", l.token())
	}

	l.next()

	l.emit(itemInfix)

	return lexDefault
}

func lexIdentifier(l *lexer) stateFn {
	l.acceptRun(stringAlphaNum + "_.")

	token := l.token()

	if token == "TRUE" || token == "FALSE" {
		l.emit(itemBool)
		return lexDefault
	}

	if token == "if" {
		l.emit(itemIf)
		return lexDefault
	}

	if token == "else" {
		l.emit(itemElse)
		return lexDefault
	}

	if token == "return" {
		l.emit(itemReturn)
		return lexDefault
	}

	if token == ".Call" {
		l.emit(itemCCall)
		return lexDefault
	}

	if token == "NULL" {
		l.emit(itemNULL)
		return lexDefault
	}

	if token == "NA" {
		l.emit(itemNA)
		return lexDefault
	}

	if token == "NA_integer_" {
		l.emit(itemNAInteger)
		return lexDefault
	}

	if token == "NA_character_" {
		l.emit(itemNACharacter)
		return lexDefault
	}

	if token == "NA_real_" {
		l.emit(itemNAReal)
		return lexDefault
	}

	if token == "NA_complex_" {
		l.emit(itemNAComplex)
		return lexDefault
	}

	l.emit(itemIdent)
	return lexDefault
}

func (l *lexer) acceptSpace() bool {
	return l.accept(" \\t")
}

func (l *lexer) acceptAlpha() bool {
	return l.accept("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

func (l *lexer) acceptNumber() bool {
	return l.accept(stringNumber)
}

func (l *lexer) acceptMathOp() bool {
	return l.accept(stringMathOp)
}

func (l *lexer) acceptAlphaNumeric() bool {
	return l.accept(stringAlphaNum)
}

func (l *lexer) accept(rs string) bool {
	for strings.IndexRune(rs, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}

	l.backup()
}
