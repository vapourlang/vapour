package lexer

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/devOpifex/vapour/token"
)

type Lexer struct {
	Input string
	start int
	pos   int
	width int
	line  int
	Items token.Items
}

const stringNumber = "0123456789"
const stringAlpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const stringAlphaNum = stringAlpha + stringNumber
const stringMathOp = "+-*/^"

func (l *Lexer) getItem(index int) token.Item {
	return l.Items[index]
}

func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.Items = append(l.Items, token.Item{
		Class: token.ItemError,
		Value: fmt.Sprintf(format, args...),
	})
	return nil
}

func (l *Lexer) emit(t token.ItemType) {
	// skip empty tokens
	if l.start == l.pos {
		return
	}

	l.Items = append(l.Items, token.Item{
		Class: t,
		Value: l.Input[l.start:l.pos],
	})
	l.start = l.pos
}

func (l *Lexer) emitEOF() {
	l.Items = append(l.Items, token.Item{Class: token.ItemEOF, Value: "EOF"})
}

// returns currently accepted token
func (l *Lexer) token() string {
	return l.Input[l.start:l.pos]
}

// next returns the next rune in the input.
func (l *Lexer) next() rune {
	if l.pos >= len(l.Input) {
		l.width = 0
		return token.EOF
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

func (l *Lexer) Run() {
	for state := lexDefault; state != nil; {
		state = state(l)
	}
}

func lexDefault(l *Lexer) stateFn {
	r1 := l.peek(1)

	if r1 == token.EOF {
		l.emitEOF()
		return nil
	}

	if r1 == '"' {
		l.next()
		l.emit(token.ItemDoubleQuote)
		return l.lexString('"')
	}

	if r1 == '\'' {
		l.next()
		l.emit(token.ItemSingleQuote)
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

	if r1 == '[' || r1 == ']' {
		l.next()
		l.next()
		l.emit(token.ItemTypesList)
		return lexIdentifier
	}

	if r1 == '\n' || r1 == ';' {
		l.next()
		l.ignore()
		//l.emit(token.ItemEOL)
		return lexDefault
	}

	// peek one more rune
	r2 := l.peek(2)

	if r1 == '.' && r2 == '.' {
		l.next()
		l.next()
		l.emit(token.ItemRange)
	}

	// if it's not %% it's an infix
	if r1 == '%' && r2 != '%' {
		return lexInfix
	}

	// it's a modulus
	if r1 == '%' && r2 == '%' {
		l.next()
		l.next()
		l.emit(token.ItemModulus)
		return lexDefault
	}

	if r1 == '=' && r2 == '=' {
		l.next()
		l.next()
		l.emit(token.ItemDoubleEqual)
		return lexDefault
	}

	if r1 == '!' && r2 == '=' {
		l.next()
		l.next()
		l.emit(token.ItemNotEqual)
		return lexDefault
	}

	if r1 == '!' {
		l.next()
		l.emit(token.ItemBang)
		return lexDefault
	}

	if r1 == '>' && r2 == '=' {
		l.next()
		l.next()
		l.emit(token.ItemGreaterOrEqual)
		return lexDefault
	}

	if r1 == '<' && r2 == '=' {
		l.next()
		l.next()
		l.emit(token.ItemLessOrEqual)
		return lexDefault
	}

	if r1 == '<' && r2 == ' ' {
		l.next()
		l.emit(token.ItemLessThan)
		return lexDefault
	}

	if r1 == '>' && r2 == ' ' {
		l.next()
		l.emit(token.ItemGreaterThan)
		return lexDefault
	}

	if r1 == '<' && r2 == '-' {
		l.next()
		l.next()
		l.emit(token.ItemAssign)
		return lexDefault
	}

	if r1 == ':' && r2 == ':' && l.peek(3) == ':' {
		l.next()
		l.next()
		l.next()
		l.emit(token.ItemNamespaceInternal)
		return lexIdentifier
	}

	if r1 == ':' && r2 == ':' {
		l.next()
		l.next()
		l.emit(token.ItemNamespace)
		return lexIdentifier
	}

	if r1 == '.' && r2 == '.' && l.peek(3) == '.' {
		l.next()
		l.next()
		l.next()
		l.emit(token.ItemThreeDot)
		return lexDefault
	}

	if r1 == '=' && r2 == '>' {
		l.next()
		l.next()
		l.emit(token.ItemArrow)
		return lexDefault
	}

	// we also emit namespace:: (above)
	// so we can assume this is not
	if r1 == ':' {
		l.next()
		l.emit(token.ItemColon)
		return lexType
	}

	if r1 == ';' {
		l.next()
		l.emit(token.ItemSemiColon)
		return lexDefault
	}

	if r1 == '&' {
		l.next()
		l.emit(token.ItemAnd)
		return lexDefault
	}

	if r1 == '|' && r2 == '>' {
		l.next()
		l.next()
		l.emit(token.ItemPipe)
		return lexDefault
	}

	if r1 == '|' {
		l.next()
		l.emit(token.ItemOr)
		return lexDefault
	}

	if r1 == '.' {
		l.next()
		l.emit(token.ItemDot)
		return lexDefault
	}

	if r1 == ',' {
		l.next()
		l.emit(token.ItemComma)
		return lexDefault
	}

	if r1 == '=' {
		l.next()
		l.emit(token.ItemAssign)
		return lexDefault
	}

	if r1 == '(' {
		l.next()
		l.emit(token.ItemLeftParen)
		return lexDefault
	}

	if r1 == ')' {
		l.next()
		l.emit(token.ItemRightParen)
		return lexType
	}

	if r1 == '{' {
		l.next()
		l.emit(token.ItemLeftCurly)
		return lexDefault
	}

	if r1 == '}' {
		l.next()
		l.emit(token.ItemRightCurly)
		return lexDefault
	}

	if r1 == '[' && r2 == '[' {
		l.next()
		l.emit(token.ItemDoubleLeftSquare)
		return lexDefault
	}

	if r1 == '[' {
		l.next()
		l.emit(token.ItemLeftSquare)
		return lexDefault
	}

	if r1 == ']' && r2 == ']' {
		l.next()
		l.emit(token.ItemDoubleRightSquare)
		return lexDefault
	}

	if r1 == ']' {
		l.next()
		l.emit(token.ItemRightSquare)
		return lexDefault
	}

	if r1 == '?' {
		l.next()
		l.emit(token.ItemQuestion)
		return lexDefault
	}

	if r1 == '`' {
		l.next()
		l.emit(token.ItemBacktick)
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

	tk := l.token()

	if tk == "+" {
		l.emit(token.ItemPlus)
	}

	if tk == "-" {
		l.emit(token.ItemMinus)
	}

	if tk == "*" {
		l.emit(token.ItemMultiply)
	}

	if tk == "/" {
		l.emit(token.ItemDivide)
	}

	if tk == "^" {
		l.emit(token.ItemPower)
	}

	return lexDefault
}

func lexNumber(l *Lexer) stateFn {
	l.acceptRun(stringNumber)

	r1 := l.peek(1)
	r2 := l.peek(2)

	if r1 == 'e' {
		l.next()
		l.acceptRun(stringNumber)
	}

	if r1 == '.' && r2 == '.' {
		l.emit(token.ItemInteger)
		l.next()
		l.next()
		l.emit(token.ItemRange)
		return lexNumber
	}

	if l.accept(".") {
		l.acceptRun(stringNumber)
		l.emit(token.ItemFloat)
		return lexDefault
	}

	l.emit(token.ItemInteger)
	return lexDefault
}

func lexComment(l *Lexer) stateFn {
	r2 := l.peek(2)

	if r2 == '\'' {
		l.next() // #
		l.next() // '

		l.emit(token.ItemSpecialComment)
		return lexSpecialComment
	}

	r := l.peek(1)
	for r != '\n' && r != token.EOF {
		l.next()
		r = l.peek(1)
	}

	l.emit(token.ItemComment)

	return lexDefault
}

func lexSpecialComment(l *Lexer) stateFn {
	r := l.peek(1)

	// not entirely certain we need
	// #'[space], e.g.: #' @param
	// @#', e.g.: #'@param
	// perhaps legal too
	if r == ' ' {
		l.next()
		l.ignore()
	}

	for r != '\n' && r != token.EOF {
		l.next()
		r = l.peek(1)
	}

	l.emit(token.ItemSpecialComment)

	return lexDefault
}

func (l *Lexer) lexString(closing rune) func(l *Lexer) stateFn {
	return func(l *Lexer) stateFn {
		var c rune
		r := l.peek(1)
		for r != closing && r != token.EOF {
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

		if r == token.EOF {
			l.next()
			return l.errorf("expecting closing quote, got %v", l.token())
		}

		l.emit(token.ItemString)

		r = l.next()

		if r == '"' {
			l.emit(token.ItemDoubleQuote)
		}

		if r == '\'' {
			l.emit(token.ItemSingleQuote)
		}

		return lexDefault
	}
}

func lexInfix(l *Lexer) stateFn {
	l.next()
	r := l.peek(1)
	for r != '%' && r != token.EOF {
		l.next()
		r = l.peek(1)
	}

	if r == token.EOF {
		l.next()
		return l.errorf("expecting closing %%, got %v", l.token())
	}

	l.next()

	l.emit(token.ItemInfix)

	return lexDefault
}

func lexIdentifier(l *Lexer) stateFn {
	l.acceptRun(stringAlphaNum + "_.")

	tk := l.token()

	if tk == "true" || tk == "false" {
		l.emit(token.ItemBool)
		return lexDefault
	}

	if tk == "if" {
		l.emit(token.ItemIf)
		return lexDefault
	}

	if tk == "else" {
		l.emit(token.ItemElse)
		return lexDefault
	}

	if tk == "return" {
		l.emit(token.ItemReturn)
		return lexDefault
	}

	if tk == "null" {
		l.emit(token.ItemNULL)
		return lexDefault
	}

	if tk == "NA" {
		l.emit(token.ItemNA)
		return lexDefault
	}

	if tk == "na_integer" {
		l.emit(token.ItemNAInteger)
		return lexDefault
	}

	if tk == "na_character" {
		l.emit(token.ItemNACharacter)
		return lexDefault
	}

	if tk == "na_real" {
		l.emit(token.ItemNAReal)
		return lexDefault
	}

	if tk == "na_complex" {
		l.emit(token.ItemNAComplex)
		return lexDefault
	}

	if tk == "inf" {
		l.emit(token.ItemInf)
		return lexDefault
	}

	if tk == "while" {
		l.emit(token.ItemWhile)
		return lexDefault
	}

	if tk == "for" {
		l.emit(token.ItemFor)
		return lexDefault
	}

	if tk == "repeat" {
		l.emit(token.ItemRepeat)
		return lexDefault
	}

	if tk == "next" {
		l.emit(token.ItemNext)
		return lexDefault
	}

	if tk == "break" {
		l.emit(token.ItemBreak)
		return lexDefault
	}

	if tk == "func" {
		l.emit(token.ItemFunction)
		return lexIdentifier
	}

	if tk == "nan" {
		l.emit(token.ItemNan)
		return lexDefault
	}

	if tk == "in" {
		l.emit(token.ItemIn)
		return lexDefault
	}

	if tk == "let" {
		l.emit(token.ItemLet)
		return lexLet
	}

	if tk == "const" {
		l.emit(token.ItemConst)
		return lexLet
	}

	if tk == "type" {
		l.emit(token.ItemTypesDecl)
		return lexTypeDeclaration
	}

	l.emit(token.ItemIdent)
	return lexDefault
}

func lexTypeDeclaration(l *Lexer) stateFn {
	// ignore space
	l.next()
	l.ignore()

	// emit custom type
	l.acceptRun(stringAlphaNum)
	l.emit(token.ItemTypes)

	// ignore space
	l.next()
	l.ignore()

	return lexType
}

func lexLet(l *Lexer) stateFn {
	r := l.peek(1)

	if r != ' ' {
		l.errorf("expecting a space, got %c", r)
		return lexDefault
	}

	// ignore the space
	l.next()
	l.ignore()

	l.acceptRun(stringAlphaNum)

	l.emit(token.ItemIdent)

	r = l.peek(1)

	if r != ':' {
		l.errorf("expecting :, got %c", r)
		return nil
	}

	// ignore the colon
	l.next()
	l.emit(token.ItemColon)

	return lexType
}

func lexType(l *Lexer) stateFn {
	r := l.peek(1)

	if r == ' ' {
		l.next()
		l.ignore()
	}

	if r == '|' {
		l.next()
		l.emit(token.ItemTypesOr)
	}

	r = l.peek(1)
	r2 := l.peek(2)
	if r == '[' && r2 == ']' {
		l.next()
		l.next()
		l.emit(token.ItemTypesList)
	}

	l.acceptRun(stringAlpha)

	l.emit(token.ItemTypes)

	r = l.peek(1)

	if r == '|' {
		l.next()
		l.emit(token.ItemTypesOr)
		return lexType
	}

	if r == ' ' {
		return lexType
	}

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
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

func itemIn(item string, items []string) bool {
	for _, v := range items {
		if v == item {
			return true
		}
	}

	return false
}
