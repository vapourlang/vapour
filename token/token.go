package token

type ItemType int

type Item struct {
	Class ItemType
	Value string
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

	// bang!
	ItemBang

	// comment
	ItemComment

	// roxygen comments
	ItemSpecialComment

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

	// function call()
	ItemFunctioCall

	// end of line \n or ;
	ItemEOL

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
)

const EOF = -1
