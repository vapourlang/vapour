package token

type ItemType int

type Item struct {
	Class ItemType
	Value string
	Line  int
	Pos   int
	Char  int
	File  string
}

type Items []Item

const (
	ItemError ItemType = iota

	// end of file
	ItemEOF

	// identifiers
	ItemIdent

	// quotes
	ItemDoubleQuote
	ItemSingleQuote

	// $
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

	// =
	ItemAssign
	ItemAssignParent

	// NULL
	ItemNULL

	// NA
	ItemNA
	ItemNan
	ItemNAString
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

	// New line \n
	ItemNewLine

	// + - / * ^
	ItemPlus
	ItemMinus
	ItemDivide
	ItemMultiply
	ItemPower
	ItemModulus

	// bang!
	ItemBang

	// comment
	ItemComment

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

	// func
	ItemFunction

	// arrow function =>
	ItemArrow

	// declare
	ItemLet
	ItemConst

	// types
	ItemTypes
	ItemTypesNew
	ItemTypesOr
	ItemTypesList
	ItemTypesDecl

	// range..
	ItemRange

	// objects
	ItemVector
	ItemDataframe
	ItemList   // list()
	ItemObject // named list()

	// @decorators
	ItemDecorator

	// attribute
	ItemAttribute
)

const EOF = -1
