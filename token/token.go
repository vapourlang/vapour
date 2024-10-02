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
	ItemTypesPkg
	ItemTypesList
	ItemTypesDecl

	// range..
	ItemRange

	// objects
	ItemObjDataframe
	ItemObjList   // list()
	ItemObjObject // named list()
	ItemObjStruct
	ItemObjMatrix
	ItemObjFunc
	ItemObjFactor
	ItemObjEnvironment

	// @decorators
	ItemDecorator
	ItemDecoratorClass
	ItemDecoratorGeneric
	ItemDecoratorDefault
	ItemDecoratorMatrix
	ItemDecoratorFactor
	ItemDecoratorEnvironment

	// attribute
	ItemAttribute

	// defer
	ItemDefer

	// += and -=
	ItemAssignInc
	ItemAssignDec
)

const EOF = -1
