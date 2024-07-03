package lexer

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type ItemType int

type Item struct {
	class ItemType
	val   string
}

type Lexer struct {
	Input string
	start int
	pos   int
	width int
	line  int
	Items []Item
}

const (
	ItemError ItemType = iota

	// identifiers
	ItemIdent

	// quotes
	ItemDoubleQuote
	ItemSingleQuote

	// dollar $ign
	ItemDollar

	// backtick
	ItemBacktick

	// infix %>%
	ItemInfix

	// comma,
	ItemComma

	// question mark?
	ItemQuestion

	// boolean
	ItemBool

	// boolean
	ItemReturn

	// ...
	ItemThreeDot

	// native pipe
	ItemPipe

	// = <-
	ItemAssign

	// .Call
	ItemC
	ItemCall
	ItemFortran

	// NULL
	ItemNULL

	// NA
	ItemNA
	ItemNan
	ItemNACharacter
	ItemNAReal
	ItemNAComplex
	ItemNAInteger

	// parens and brackets
	ItemLeftCurly
	ItemRightCurly
	ItemLeftParen
	ItemRightParen
	ItemLeftSquare
	ItemRightSquare
	ItemDoubleLeftSquare
	ItemDoubleRightSquare

	// "strings"
	ItemString

	// numbers
	ItemInteger
	ItemFloat

	// namespace::
	ItemNamespace
	// namespace:::
	ItemNamespaceInternal

	// colon
	ItemColon

	// semicolon;
	ItemSemiColon

	// + - / * ^
	ItemPlus
	ItemMinus
	ItemDivide
	ItemMultiply
	ItemPower
	ItemModulus

	// comment
	ItemComment

	// roxygen comments
	ItemSpecialComment
	ItemRoxygenTagAt
	ItemRoxygenTag
	ItemRoxygenTagContent

	// doctor tags
	ItemTypeDef
	ItemTypeVar

	// compare
	ItemDoubleEqual
	ItemLessThan
	ItemGreaterThan
	ItemNotEqual
	ItemLessOrEqual
	ItemGreaterOrEqual

	// if else
	ItemIf
	ItemElse
	ItemAnd
	ItemOr
	ItemBreak

	// Infinite
	ItemInf

	// loop
	ItemFor
	ItemRepeat
	ItemWhile
	ItemNext
	ItemIn

	// function()
	ItemFunction

	// end of line \n or ;
	ItemEOL
)

const stringNumber = "0123456789"
const stringAlpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const stringAlphaNum = stringAlpha + stringNumber
const stringMathOp = "+-*/^"

const eof = -1

func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.Items = append(l.Items, Item{ItemError, fmt.Sprintf(format, args...)})
	return nil
}

func (l *Lexer) emit(t ItemType) {
	// skip empty tokens
	if l.start == l.pos {
		return
	}

	l.Items = append(l.Items, Item{t, l.Input[l.start:l.pos]})
	l.start = l.pos
}

// returns currently accepted token
func (l *Lexer) token() string {
	return l.Input[l.start:l.pos]
}

// next returns the next rune in the input.
func (l *Lexer) next() rune {
	if l.pos >= len(l.Input) {
		l.width = 0
		return eof
	}

	r, w := utf8.DecodeRuneInString(l.Input[l.pos:])
	l.width = w
	l.pos += l.width

	if r == '\n' {
		l.line++
	}

	return r
}

func (l *Lexer) skipLine() {
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

func (l *Lexer) ignore() {
	l.start = l.pos
}

func (l *Lexer) backup() {
	l.pos -= l.width
}

func (l *Lexer) peek(n int) rune {
	var r rune
	for i := 0; i < n; i++ {
		r = l.next()
	}

	for i := 0; i < n; i++ {
		l.backup()
	}
	return r
}

type stateFn func(*Lexer) stateFn

func (l *Lexer) run() {
	for state := lexDefault; state != nil; {
		state = state(l)
	}
}

func lexDefault(l *Lexer) stateFn {
	r1 := l.peek(1)

	if r1 == eof {
		return nil
	}

	if r1 == '"' {
		l.next()
		l.emit(ItemDoubleQuote)
		return l.lexString('"')
	}

	if r1 == '\'' {
		l.next()
		l.emit(ItemSingleQuote)
		return l.lexString('\'')
	}

	if r1 == '#' {
		return lexComment
	}

	// we parsed strings: we skip spaces and tabs
	if r1 == ' ' || r1 == '\t' {
		l.next()
		l.ignore()
		return lexDefault
	}

	if r1 == '\n' || r1 == ';' {
		l.next()
		l.emit(ItemEOL)
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
		l.emit(ItemModulus)
		return lexDefault
	}

	if r1 == '=' && r2 == '=' {
		l.next()
		l.next()
		l.emit(ItemDoubleEqual)
		return lexDefault
	}

	if r1 == '!' && r2 == '=' {
		l.next()
		l.next()
		l.emit(ItemNotEqual)
		return lexDefault
	}

	if r1 == '>' && r2 == '=' {
		l.next()
		l.next()
		l.emit(ItemGreaterOrEqual)
		return lexDefault
	}

	if r1 == '<' && r2 == '=' {
		l.next()
		l.next()
		l.emit(ItemLessOrEqual)
		return lexDefault
	}

	if r1 == '<' && r2 == ' ' {
		l.next()
		l.emit(ItemLessThan)
		return lexDefault
	}

	if r1 == '>' && r2 == ' ' {
		l.next()
		l.emit(ItemGreaterThan)
		return lexDefault
	}

	if r1 == '<' && r2 == '-' {
		l.next()
		l.next()
		l.emit(ItemAssign)
		return lexDefault
	}

	if r1 == ':' && r2 == ':' && l.peek(3) == ':' {
		l.next()
		l.next()
		l.next()
		l.emit(ItemNamespaceInternal)
		return lexIdentifier
	}

	if r1 == ':' && r2 == ':' {
		l.next()
		l.next()
		l.emit(ItemNamespace)
		return lexIdentifier
	}

	if r1 == '.' && r2 == '.' && l.peek(3) == '.' {
		l.next()
		l.next()
		l.next()
		l.emit(ItemThreeDot)
		return lexDefault
	}

	// we also emit namespace:: (above)
	// so we can assume this is not
	if r1 == ':' {
		l.next()
		l.emit(ItemColon)
		return lexDefault
	}

	if r1 == ';' {
		l.next()
		l.emit(ItemSemiColon)
		return lexDefault
	}

	if r1 == '&' {
		l.next()
		l.emit(ItemAnd)
		return lexDefault
	}

	if r1 == '|' && r2 == '>' {
		l.next()
		l.next()
		l.emit(ItemPipe)
		return lexDefault
	}

	if r1 == '|' {
		l.next()
		l.emit(ItemOr)
		return lexDefault
	}

	if r1 == '$' {
		l.next()
		l.emit(ItemDollar)
		return lexDefault
	}

	if r1 == ',' {
		l.next()
		l.emit(ItemComma)
		return lexDefault
	}

	if r1 == '=' {
		l.next()
		l.emit(ItemAssign)
		return lexDefault
	}

	if r1 == '(' {
		l.next()
		l.emit(ItemLeftParen)
		return lexDefault
	}

	if r1 == ')' {
		l.next()
		l.emit(ItemLeftParen)
		return lexDefault
	}

	if r1 == '{' {
		l.next()
		l.emit(ItemLeftCurly)
		return lexDefault
	}

	if r1 == '}' {
		l.next()
		l.emit(ItemRightCurly)
		return lexDefault
	}

	if r1 == '[' && r2 == '[' {
		l.next()
		l.emit(ItemDoubleLeftSquare)
		return lexDefault
	}

	if r1 == '[' {
		l.next()
		l.emit(ItemLeftSquare)
		return lexDefault
	}

	if r1 == ']' && r2 == ']' {
		l.next()
		l.emit(ItemDoubleRightSquare)
		return lexDefault
	}

	if r1 == ']' {
		l.next()
		l.emit(ItemRightSquare)
		return lexDefault
	}

	if r1 == '?' {
		l.next()
		l.emit(ItemQuestion)
		return lexDefault
	}

	if r1 == '`' {
		l.next()
		l.emit(ItemBacktick)
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

func lexMathOp(l *Lexer) stateFn {
	l.acceptRun(stringMathOp)

	token := l.token()

	if token == "+" {
		l.emit(ItemPlus)
	}

	if token == "-" {
		l.emit(ItemMinus)
	}

	if token == "*" {
		l.emit(ItemMultiply)
	}

	if token == "/" {
		l.emit(ItemDivide)
	}

	if token == "^" {
		l.emit(ItemPower)
	}

	return lexDefault
}

func lexNumber(l *Lexer) stateFn {
	l.acceptRun(stringNumber)

	r := l.peek(1)

	if r == 'e' {
		l.next()
		l.acceptRun(stringNumber)
	}

	if l.accept(".") {
		l.acceptRun(stringNumber)
		l.emit(ItemFloat)
		return lexDefault
	}

	l.emit(ItemInteger)
	return lexDefault
}

func lexComment(l *Lexer) stateFn {
	r2 := l.peek(2)

	if r2 == '\'' {
		l.next() // #
		l.next() // '

		l.emit(ItemSpecialComment)
		return lexSpecialComment
	}

	r := l.peek(1)
	for r != '\n' && r != eof {
		l.next()
		r = l.peek(1)
	}

	l.emit(ItemComment)

	return lexDefault
}

func lexSpecialComment(l *Lexer) stateFn {
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
		l.emit(ItemRoxygenTagAt)
		return lexRoxygen
	}

	for r != '\n' && r != eof {
		l.next()
		r = l.peek(1)
	}

	l.emit(ItemSpecialComment)

	return lexDefault
}

func lexRoxygen(l *Lexer) stateFn {
	r := l.peek(1)
	for r != ' ' && r != '\t' && r != '\n' && r != eof {
		l.next()
		r = l.peek(1)
	}

	token := l.token()

	l.emit(ItemRoxygenTag)

	if token == "type" {
		return lexTypeTag
	}

	if token == "yield" {
		return lexTypes
	}

	return lexRoxygenTagContent
}

func lexRoxygenTagContent(l *Lexer) stateFn {
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

	l.emit(ItemRoxygenTagContent)

	return lexDefault
}

func lexTypeTag(l *Lexer) stateFn {
	r := l.peek(1)
	for r != ':' && r != '\n' && r != eof {
		l.next()
		r = l.peek(1)
	}

	if r != ':' {
		l.next()
		return l.errorf("expects `:`, found %v [@type variable: type]", l.token())
	}

	l.emit(ItemTypeVar)

	// ignore colon
	// e.g.: @type x: numeric
	l.next()
	l.ignore()

	return lexTypes
}

func lexTypes(l *Lexer) stateFn {
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

	l.emit(ItemTypeDef)

	return lexTypes
}

func (l *Lexer) lexString(closing rune) func(l *Lexer) stateFn {
	return func(l *Lexer) stateFn {
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

		l.emit(ItemString)

		r = l.next()

		if r == '"' {
			l.emit(ItemDoubleQuote)
		}

		if r == '\'' {
			l.emit(ItemSingleQuote)
		}

		return lexDefault
	}
}

func lexInfix(l *Lexer) stateFn {
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

	l.emit(ItemInfix)

	return lexDefault
}

func lexIdentifier(l *Lexer) stateFn {
	l.acceptRun(stringAlphaNum + "_.")

	token := l.token()

	if token == "TRUE" || token == "FALSE" {
		l.emit(ItemBool)
		return lexDefault
	}

	if token == "if" {
		l.emit(ItemIf)
		return lexDefault
	}

	if token == "else" {
		l.emit(ItemElse)
		return lexDefault
	}

	if token == "return" {
		l.emit(ItemReturn)
		return lexDefault
	}

	if token == ".Call" {
		l.emit(ItemCall)
		return lexDefault
	}

	if token == ".C" {
		l.emit(ItemC)
		return lexDefault
	}

	if token == ".Fortran" {
		l.emit(ItemFortran)
		return lexDefault
	}

	if token == "NULL" {
		l.emit(ItemNULL)
		return lexDefault
	}

	if token == "NA" {
		l.emit(ItemNA)
		return lexDefault
	}

	if token == "NA_integer_" {
		l.emit(ItemNAInteger)
		return lexDefault
	}

	if token == "NA_character_" {
		l.emit(ItemNACharacter)
		return lexDefault
	}

	if token == "NA_real_" {
		l.emit(ItemNAReal)
		return lexDefault
	}

	if token == "NA_complex_" {
		l.emit(ItemNAComplex)
		return lexDefault
	}

	if token == "Inf" {
		l.emit(ItemInf)
		return lexDefault
	}

	if token == "while" {
		l.emit(ItemWhile)
		return lexDefault
	}

	if token == "for" {
		l.emit(ItemFor)
		return lexDefault
	}

	if token == "repeat" {
		l.emit(ItemRepeat)
		return lexDefault
	}

	if token == "next" {
		l.emit(ItemNext)
		return lexDefault
	}

	if token == "break" {
		l.emit(ItemBreak)
		return lexDefault
	}

	if token == "function" {
		l.emit(ItemFunction)
		return lexDefault
	}

	if token == "NaN" {
		l.emit(ItemNan)
		return lexDefault
	}

	if token == "in" {
		l.emit(ItemIn)
		return lexDefault
	}

	l.emit(ItemIdent)
	return lexDefault
}

func (l *Lexer) acceptSpace() bool {
	return l.accept(" \\t")
}

func (l *Lexer) acceptAlpha() bool {
	return l.accept("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

func (l *Lexer) acceptNumber() bool {
	return l.accept(stringNumber)
}

func (l *Lexer) acceptMathOp() bool {
	return l.accept(stringMathOp)
}

func (l *Lexer) acceptAlphaNumeric() bool {
	return l.accept(stringAlphaNum)
}

func (l *Lexer) accept(rs string) bool {
	for strings.IndexRune(rs, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *Lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}

	l.backup()
}
