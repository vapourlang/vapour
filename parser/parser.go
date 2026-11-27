package parser

import (
	"strings"

	"fmt"

	"github.com/vapourlang/vapour/ast"
	"github.com/vapourlang/vapour/diagnostics"
	"github.com/vapourlang/vapour/lexer"
	"github.com/vapourlang/vapour/token"
)

// operator precedence, from lowest to highest, follows R's
// see ?Syntax in R
const (
	_ int = iota
	LOWEST
	ASSIGN      // = <- += -=
	OR          // | ||
	AND         // & &&
	NOT         // !X
	EQUALS      // == !=
	LESSGREATER // > < >= <=
	SUM         // + -
	PRODUCT     // * /
	SPECIAL     // %any% %% |>
	RANGE       // ..
	PREFIX      // -X
	POWER       // ^
	DOLLAR      // $
	NAMESPACE   // :: :::
	CALL        // call(X)
	INDEX       // x[i] x[[i]]
)

var precedences = map[token.ItemType]int{
	token.ItemAssign:            ASSIGN,
	token.ItemAssignInc:         ASSIGN,
	token.ItemAssignDec:         ASSIGN,
	token.ItemAssignParent:      ASSIGN,
	token.ItemOr:                OR,
	token.ItemAnd:               AND,
	token.ItemDoubleEqual:       EQUALS,
	token.ItemNotEqual:          EQUALS,
	token.ItemLessThan:          LESSGREATER,
	token.ItemGreaterThan:       LESSGREATER,
	token.ItemLessOrEqual:       LESSGREATER,
	token.ItemGreaterOrEqual:    LESSGREATER,
	token.ItemPlus:              SUM,
	token.ItemMinus:             SUM,
	token.ItemDivide:            PRODUCT,
	token.ItemMultiply:          PRODUCT,
	token.ItemPipe:              SPECIAL,
	token.ItemInfix:             SPECIAL,
	token.ItemModulus:           SPECIAL,
	token.ItemRange:             RANGE,
	token.ItemPower:             POWER,
	token.ItemDollar:            DOLLAR,
	token.ItemNamespace:         NAMESPACE,
	token.ItemNamespaceInternal: NAMESPACE,
	token.ItemLeftParen:         CALL,
	token.ItemLeftSquare:        INDEX,
	token.ItemDoubleLeftSquare:  INDEX,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors diagnostics.Diagnostics

	pos int

	curToken  token.Item
	peekToken token.Item

	prefixParseFns map[token.ItemType]prefixParseFn
	infixParseFns  map[token.ItemType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: diagnostics.Diagnostics{},
	}

	p.prefixParseFns = make(map[token.ItemType]prefixParseFn)
	p.registerPrefix(token.ItemIdent, p.parseIdentifier)
	p.registerPrefix(token.ItemComma, p.parseComma)
	p.registerPrefix(token.ItemAttribute, p.parseAttribute)
	p.registerPrefix(token.ItemInteger, p.parseIntegerLiteral)
	p.registerPrefix(token.ItemFloat, p.parseFloatLiteral)
	p.registerPrefix(token.ItemBang, p.parsePrefixExpression)
	p.registerPrefix(token.ItemMinus, p.parsePrefixExpression)
	p.registerPrefix(token.ItemBool, p.parseBoolean)
	p.registerPrefix(token.ItemLeftParen, p.parseGroupedExpression)
	p.registerPrefix(token.ItemIf, p.parseIfExpression)
	p.registerPrefix(token.ItemFunction, p.parseFunctionLiteral)
	p.registerPrefix(token.ItemDoubleQuote, p.parseStringLiteral)
	p.registerPrefix(token.ItemSingleQuote, p.parseStringLiteral)
	p.registerPrefix(token.ItemNA, p.parseNA)
	p.registerPrefix(token.ItemNan, p.parseNan)
	p.registerPrefix(token.ItemNAComplex, p.parseNaComplex)
	p.registerPrefix(token.ItemNAReal, p.parseNaReal)
	p.registerPrefix(token.ItemNAInteger, p.parseNaInteger)
	p.registerPrefix(token.ItemInf, p.parseInf)
	p.registerPrefix(token.ItemNULL, p.parseNull)
	p.registerPrefix(token.ItemThreeDot, p.parseElipsis)
	p.registerPrefix(token.ItemString, p.parseNaString)
	p.registerPrefix(token.ItemFor, p.parseFor)
	p.registerPrefix(token.ItemWhile, p.parseWhile)
	p.registerPrefix(token.ItemRepeat, p.parseRepeat)
	p.registerPrefix(token.ItemBreak, p.parseLoopControl)
	p.registerPrefix(token.ItemNext, p.parseLoopControl)
	p.registerPrefix(token.ItemDecoratorClass, p.parseDecoratorClass)
	p.registerPrefix(token.ItemDecoratorGeneric, p.parseDecoratorGeneric)
	p.registerPrefix(token.ItemDecoratorDefault, p.parseDecoratorDefault)
	p.registerPrefix(token.ItemDecoratorMatrix, p.parseDecoratorMatrix)
	p.registerPrefix(token.ItemDecoratorFactor, p.parseDecoratorFactor)
	p.registerPrefix(token.ItemDecoratorEnvironment, p.parseDecoratorEnvironment)

	p.infixParseFns = make(map[token.ItemType]infixParseFn)
	p.registerInfix(token.ItemPlus, p.parseInfixExpression)
	p.registerInfix(token.ItemMinus, p.parseInfixExpression)
	p.registerInfix(token.ItemDivide, p.parseInfixExpression)
	p.registerInfix(token.ItemMultiply, p.parseInfixExpression)
	p.registerInfix(token.ItemAssign, p.parseInfixExpression)
	p.registerInfix(token.ItemAssignInc, p.parseInfixExpression)
	p.registerInfix(token.ItemAssignDec, p.parseInfixExpression)
	p.registerInfix(token.ItemAssignParent, p.parseInfixExpression)
	p.registerInfix(token.ItemDoubleEqual, p.parseInfixExpression)
	p.registerInfix(token.ItemNotEqual, p.parseInfixExpression)
	p.registerInfix(token.ItemLessThan, p.parseInfixExpression)
	p.registerInfix(token.ItemGreaterThan, p.parseInfixExpression)
	p.registerInfix(token.ItemPipe, p.parseInfixExpression)
	p.registerInfix(token.ItemComma, p.parseInfixExpression)
	p.registerInfix(token.ItemDollar, p.parseInfixExpression)
	p.registerInfix(token.ItemRange, p.parseInfixExpression)
	p.registerInfix(token.ItemNamespace, p.parseInfixExpression)
	p.registerInfix(token.ItemNamespaceInternal, p.parseInfixExpression)
	p.registerInfix(token.ItemLeftSquare, p.parseIndexExpression)
	p.registerInfix(token.ItemDoubleLeftSquare, p.parseIndexExpression)
	p.registerInfix(token.ItemAnd, p.parseInfixExpression)
	p.registerInfix(token.ItemOr, p.parseInfixExpression)
	p.registerInfix(token.ItemLessOrEqual, p.parseInfixExpression)
	p.registerInfix(token.ItemGreaterOrEqual, p.parseInfixExpression)
	p.registerInfix(token.ItemModulus, p.parseInfixExpression)
	p.registerInfix(token.ItemPower, p.parseInfixExpression)
	p.registerInfix(token.ItemInfix, p.parseInfixExpression)

	p.registerInfix(token.ItemLeftParen, p.parseCallExpression)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	if p.pos >= len(p.l.Items) {
		return
	}
	p.peekToken = p.l.Items[p.pos]
	p.pos++
}

// tokenBefore returns the token preceding the current token,
// or an empty token if the current token is the first one.
func (p *Parser) tokenBefore() token.Item {
	// invariant: curToken is Items[pos-2]
	if p.pos-3 < 0 {
		return token.Item{Class: token.ItemNewLine}
	}
	return p.l.Items[p.pos-3]
}

// missingClosing reports an unclosed delimiter, e.g.: print((1 + 2)
func (p *Parser) missingClosing(closing string) {
	p.errors = append(
		p.errors,
		diagnostics.NewError(p.curToken, fmt.Sprintf("missing closing `%v`", closing)),
	)
}

func (p *Parser) previousToken(n int) {
	for i := 0; i < n; i++ {
		p.pos--
		p.curToken = p.l.Items[p.pos-2]
		p.peekToken = p.l.Items[p.pos-1]
	}
}

func (p *Parser) print() {
	fmt.Println("++++++++++++++++++++ Current ++++++++++++++++++++")
	fmt.Printf("line: %v - character: %v | ", p.curToken.Line+1, p.curToken.Char+1)
	p.curToken.Print()
	fmt.Println("++++++++++++++++++++ Peek")
	fmt.Printf("line: %v - character: %v | ", p.peekToken.Line+1, p.peekToken.Char+1)
	p.peekToken.Print()
}

func (p *Parser) curTokenIs(t token.ItemType) bool {
	return p.curToken.Class == t
}

func (p *Parser) peekTokenIs(t token.ItemType) bool {
	return p.peekToken.Class == t
}

func (p *Parser) expectPeek(t token.ItemType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) expectCurrent(t token.ItemType) bool {
	if p.curTokenIs(t) {
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) HasError() bool {
	return len(p.errors) > 0
}

func (p *Parser) Errors() diagnostics.Diagnostics {
	return p.errors.UniqueLine()
}

func (p *Parser) peekError(t token.ItemType) {
	// we already got an error on the lexer: use it
	if p.peekToken.Class == token.ItemError {
		return
	}

	msg := fmt.Sprintf(
		"expected next token to be `%v`, got `%v` instead",
		t,
		p.peekToken.Class,
	)

	p.errors = append(
		p.errors,
		diagnostics.NewError(p.curToken, msg),
	)
}

func (p *Parser) noPrefixParseFnError(t token.ItemType) {
	msg := fmt.Sprintf(
		"no prefix parse function for `%v` found",
		t,
	)
	p.errors = append(
		p.errors,
		diagnostics.NewError(p.curToken, msg),
	)
}

func (p *Parser) Run() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.ItemEOF) && !p.curTokenIs(token.ItemError) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Class {
	case token.ItemLet:
		return p.parseLetStatement()
	case token.ItemConst:
		return p.parseConstStatement()
	case token.ItemReturn:
		return p.parseReturnStatement()
	case token.ItemDefer:
		return p.parseDeferStatement()
	case token.ItemComment:
		return p.parseCommentStatement()
	case token.ItemNewLine:
		return p.parseNewLine()
	case token.ItemTypesDecl:
		return p.parseTypeDeclarations()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseFor() ast.Expression {
	lit := &ast.For{
		Token: p.curToken,
	}

	if !p.expectPeek(token.ItemLeftParen) {
		return nil
	}

	if !p.expectPeek(token.ItemLet) {
		return nil
	}

	lit.Name = p.parseLetStatement()

	if !p.expectPeek(token.ItemIn) {
		return nil
	}

	p.nextToken()

	lit.Vector = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.ItemRightParen) {
		p.nextToken()
	}

	if !p.curTokenIs(token.ItemRightParen) {
		return nil
	}

	p.skipNewLine()

	if !p.expectPeek(token.ItemLeftCurly) {
		return nil
	}

	lit.Value = p.parseBlockStatement()

	p.nextToken()

	return lit
}

func (p *Parser) parseWhile() ast.Expression {
	lit := &ast.While{
		Token: p.curToken,
	}

	if !p.expectPeek(token.ItemLeftParen) {
		return nil
	}

	p.nextToken()

	lit.Statement = p.parseStatement()

	if p.peekTokenIs(token.ItemRightParen) {
		p.nextToken()
	}

	if p.peekTokenIs(token.ItemLeftCurly) {
		p.nextToken()
	}

	lit.Value = p.parseBlockStatement()

	if p.peekTokenIs(token.ItemRightCurly) {
		p.nextToken()
	}

	return lit
}

func (p *Parser) parseRepeat() ast.Expression {
	lit := &ast.Repeat{
		Token: p.curToken,
	}

	if !p.expectPeek(token.ItemLeftCurly) {
		return nil
	}

	lit.Value = p.parseBlockStatement()

	if p.peekTokenIs(token.ItemRightCurly) {
		p.nextToken()
	}

	return lit
}

// parseLoopControl parses break and next
func (p *Parser) parseLoopControl() ast.Expression {
	return &ast.LoopControl{
		Token: p.curToken,
		Value: p.curToken.Value,
	}
}

func (p *Parser) parseComma() ast.Expression {
	return &ast.Comma{Token: p.curToken}
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Value}
}

func (p *Parser) parseAttribute() ast.Expression {
	return &ast.Attribute{Token: p.curToken, Value: p.curToken.Value}
}

func (p *Parser) parseNull() ast.Expression {
	return &ast.Null{
		Token: p.curToken,
		Value: "NULL",
		Type:  &ast.Type{Name: "null"},
	}
}

func (p *Parser) parseElipsis() ast.Expression {
	return &ast.Keyword{Token: p.curToken, Value: "..."}
}

func (p *Parser) parseNA() ast.Expression {
	return &ast.Keyword{
		Token: p.curToken,
		Value: "NA",
		Type:  &ast.Type{Name: ""},
	}
}

func (p *Parser) parseNan() ast.Expression {
	return &ast.Keyword{
		Token: p.curToken,
		Value: "NaN",
		Type:  &ast.Type{Name: ""},
	}
}

func (p *Parser) parseNaString() ast.Expression {
	return &ast.Keyword{
		Token: p.curToken,
		Value: "NA_character_",
		Type:  &ast.Type{Name: "na_char"},
	}
}

func (p *Parser) parseNaReal() ast.Expression {
	return &ast.Keyword{
		Token: p.curToken,
		Value: "NA_real_",
		Type:  &ast.Type{Name: "na_real"},
	}
}

func (p *Parser) parseNaComplex() ast.Expression {
	return &ast.Keyword{
		Token: p.curToken,
		Value: "NA_complex_",
		Type:  &ast.Type{Name: "na_complex"},
	}
}

func (p *Parser) parseNaInteger() ast.Expression {
	return &ast.Keyword{
		Token: p.curToken,
		Value: "NA_integer_",
		Type:  &ast.Type{Name: "na_int"},
	}
}

func (p *Parser) parseInf() ast.Expression {
	return &ast.Keyword{
		Token: p.curToken,
		Value: "Inf",
		Type:  &ast.Type{Name: "num"},
	}
}

func (p *Parser) parseTypeDeclarations() ast.Statement {
	p.nextToken()
	p.nextToken()
	if p.peekToken.Value == "func" {
		p.previousToken(2)
		return p.parseTypeDeclarationFunc()
	}
	p.previousToken(2)
	return p.parseTypeDeclaration()
}

func (p *Parser) parseTypeDeclarationFunc() *ast.TypeFunction {
	fn := &ast.TypeFunction{
		Token: p.curToken,
	}

	// expect the custom type
	if !p.expectPeek(token.ItemTypes) {
		return nil
	}

	fn.Name = p.curToken.Value

	if !p.expectPeek(token.ItemColon) {
		return nil
	}

	p.nextToken()

	if !p.expectPeek(token.ItemLeftParen) {
		return nil
	}

	for !p.peekTokenIs(token.ItemRightParen) {
		// signatures are on a single line
		if p.peekTokenIs(token.ItemEOF) || p.peekTokenIs(token.ItemNewLine) {
			p.missingClosing(")")
			return nil
		}
		fn.Arguments = append(fn.Arguments, p.parseTypes())
		if p.peekTokenIs(token.ItemComma) {
			p.nextToken()
		}
	}

	if !p.expectPeek(token.ItemRightParen) {
		return nil
	}

	// parseTypes expects the return type to be the next token
	fn.Return = p.parseTypes()

	return fn
}

func (p *Parser) parseTypeDeclaration() *ast.TypeStatement {
	typ := &ast.TypeStatement{
		Token:      p.curToken,
		Type:       ast.Types{},
		Attributes: []*ast.TypeAttributesStatement{},
	}

	// expect the custom type
	if !p.expectPeek(token.ItemTypes) {
		return nil
	}

	typ.Name = p.curToken.Value

	// expect colon
	if !p.expectPeek(token.ItemColon) {
		return nil
	}

	if p.peekTokenIs(token.ItemTypes) {
		typ.Object = "vector"
		typ.Type = p.parseTypes()
		p.nextToken()
		return typ
	}

	if p.peekTokenIs(token.ItemTypesList) {
		typ.Object = "impliedList"
		typ.Type = p.parseTypes()
		p.nextToken()
		return typ
	}

	if p.peekTokenIs(token.ItemObjList) {
		typ.Object = "list"
		p.nextToken()
		if !p.expectPeek(token.ItemLeftCurly) {
			return nil
		}
		typ.Type = p.parseTypes()
		p.skipNewLine()
		p.nextToken()
		p.nextToken()
		return typ
	}

	if p.peekTokenIs(token.ItemObjFactor) {
		typ.Object = "factor"
		p.nextToken()
		if !p.expectPeek(token.ItemLeftCurly) {
			return nil
		}
		typ.Type = p.parseTypes()
		p.skipNewLine()
		p.nextToken()
		p.nextToken()
		return typ
	}

	if p.peekTokenIs(token.ItemObjMatrix) {
		typ.Object = "matrix"
		p.nextToken()
		if !p.expectPeek(token.ItemLeftCurly) {
			return nil
		}
		typ.Type = p.parseTypes()
		p.skipNewLine()
		p.nextToken()
		p.nextToken()
		return typ
	}

	if p.peekTokenIs(token.ItemObjStruct) {
		typ.Object = "struct"
		p.nextToken()
		p.nextToken()
		typ.Type = p.parseTypes()
		p.skipNewLine()

		// struct with no attributes
		if p.peekTokenIs(token.ItemRightCurly) {
			p.nextToken()
			p.nextToken()
			return typ
		}

		if !p.expectPeek(token.ItemComma) {
			return nil
		}
		p.skipNewLine()
	}

	if p.peekTokenIs(token.ItemObjEnvironment) {
		p.nextToken()
		p.nextToken()
		typ.Object = "environment"
		p.skipNewLine()
	}

	if p.peekTokenIs(token.ItemObjObject) {
		p.nextToken()
		p.nextToken()
		typ.Object = "object"
		p.skipNewLine()
	}

	if p.peekTokenIs(token.ItemObjDataframe) {
		p.nextToken()
		p.nextToken()
		typ.Object = "dataframe"
		p.skipNewLine()
	}

	typ.Attributes = p.parseTypeAttributes()

	return typ
}

func (p *Parser) parseTypeAttributes() []*ast.TypeAttributesStatement {
	var attrs []*ast.TypeAttributesStatement

	for !p.peekTokenIs(token.ItemRightCurly) && !p.peekTokenIs(token.ItemEOF) {
		p.nextToken()
		attrs = append(attrs, p.parseTypeAttribute())
		if p.curTokenIs(token.ItemRightCurly) {
			return attrs
		}
	}

	p.nextToken()

	return attrs
}

func (p *Parser) parseTypeAttribute() *ast.TypeAttributesStatement {
	if p.curTokenIs(token.ItemNewLine) {
		p.nextToken()
	}

	attr := &ast.TypeAttributesStatement{}

	attr.Name = p.curToken.Value

	if !p.expectPeek(token.ItemColon) {
		return nil
	}

	attr.Type = p.parseTypes()

	p.nextToken()

	return attr
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.ItemIdent) {
		return nil
	}

	stmt.Name = p.curToken.Value

	if !p.expectPeek(token.ItemColon) {
		return nil
	}

	stmt.Type = p.parseTypes()

	if !p.peekTokenIs(token.ItemAssign) {
		return stmt
	}

	p.nextToken()
	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.ItemNewLine) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseConstStatement() *ast.ConstStatement {
	stmt := &ast.ConstStatement{Token: p.curToken}

	if !p.expectPeek(token.ItemIdent) {
		return nil
	}

	stmt.Name = p.curToken.Value

	if !p.expectPeek(token.ItemColon) {
		return nil
	}

	stmt.Type = p.parseTypes()

	if !p.expectPeek(token.ItemAssign) {
		return stmt
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.ItemNewLine) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseDeferStatement() *ast.DeferStatement {
	stmt := &ast.DeferStatement{Token: p.curToken}

	p.nextToken()

	stmt.Func = p.parseExpression(LOWEST)

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.ItemNewLine) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	if p.curTokenIs(token.ItemLeftParen) {
		tk := p.tokenBefore()

		// skip paren left (
		p.nextToken()

		// no prefix, it's not a call
		// it's either a vector (1, 2, 3)
		// or an anonymous function
		// (x: char): char => { print(x) }
		if tk.Class != token.ItemIdent {
			i := 0
			for !p.curTokenIs(token.ItemRightParen) {
				if p.curTokenIs(token.ItemEOF) {
					p.missingClosing(")")
					return nil
				}
				i++
				p.nextToken()
			}

			// if the closing paren ) is followed by
			// a type or => it's an anonymous function
			if p.peekTokenIs(token.ItemColon) || p.peekTokenIs(token.ItemArrow) {
				p.previousToken(i)
				return p.parseAnonymousFunction()
			}

			// otherwise it's a vector
			p.previousToken(i + 1)
			return p.parseVector()
		}
	}

	// skip newline as pipe is often followed by a new line
	if p.curTokenIs(token.ItemNewLine) {
		p.nextToken()
	}

	prefix := p.prefixParseFns[p.curToken.Class]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Class)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.ItemEOF) &&
		precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Class]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Class]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Class]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	return &ast.IntegerLiteral{
		Token: p.curToken,
		Value: p.curToken.Value,
		Type:  &ast.Type{Name: "int", List: false},
	}
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	typ := "num"

	// complex, e.g.: 1i
	if strings.HasSuffix(p.curToken.Value, "i") {
		typ = "complex"
	}

	return &ast.FloatLiteral{
		Token: p.curToken,
		Value: p.curToken.Value,
		Type:  &ast.Type{Name: typ, List: false},
	}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.curToken,
		Value: p.curToken.Value == "true" || p.curToken.Value == "TRUE",
		Type:  &ast.Type{Name: "bool"},
	}
}

func (p *Parser) parseCommentStatement() ast.Statement {
	return &ast.CommentStatement{Token: p.curToken, Value: p.curToken.Value}
}

func (p *Parser) parseNewLine() ast.Statement {
	return &ast.NewLine{Token: p.curToken}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	str := &ast.StringLiteral{
		Token: p.curToken,
		Type:  &ast.Type{Name: "char"},
	}

	// it's an empty string ""
	if p.peekTokenIs(p.curToken.Class) {
		p.nextToken()
		return str
	}

	p.expectPeek(token.ItemString)

	str.Str = p.curToken.Value

	p.nextToken()

	return str
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Value,
	}

	p.nextToken()

	// in R ! binds looser than comparisons: !x == y is !(x == y)
	precedence := PREFIX
	if expression.Operator == "!" {
		precedence = NOT
	}

	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	operator := p.curToken.Value

	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: operator,
		Left:     left,
	}

	precedence := p.curPrecedence()

	// right associative: 2^3^2 is 2^(3^2), x = y = 1 is x = (y = 1)
	if precedence == POWER || precedence == ASSIGN {
		precedence--
	}

	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseVector() ast.Expression {
	vec := &ast.VectorLiteral{
		Token: p.curToken,
	}

	for !p.peekTokenIs(token.ItemRightParen) {
		if p.peekTokenIs(token.ItemEOF) {
			p.missingClosing(")")
			return nil
		}
		p.nextToken()
		if p.curTokenIs(token.ItemComma) || p.peekTokenIs(token.ItemNewLine) {
			continue
		}
		vec.Value = append(vec.Value, p.parseExpression(LOWEST))
	}

	p.nextToken()
	p.nextToken()

	return vec
}

func (p *Parser) parseAnonymousFunction() ast.Expression {
	lit := &ast.FunctionLiteral{}

	p.previousToken(1)
	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.ItemColon) {
		return nil
	}

	// parse types
	lit.ReturnType = p.parseTypes()

	lit.Name = ""

	if !p.expectPeek(token.ItemArrow) {
		return nil
	}

	// set token as => so we know it's anonymous downstream
	lit.Token = p.curToken

	if !p.expectPeek(token.ItemLeftCurly) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	tk := p.tokenBefore()

	// skip paren left (
	p.nextToken()

	// no prefix, it's not a call
	// it's either a vector (1, 2, 3)
	// or an anonymous function
	// (x: char): char => { print(x) }
	if tk.Class != token.ItemIdent {
		i := 0
		for !p.curTokenIs(token.ItemRightParen) {
			if p.curTokenIs(token.ItemEOF) {
				p.missingClosing(")")
				return nil
			}
			i++
			p.nextToken()
		}

		// if the closing paren ) is followed by
		// a type or => it's an anonymous function
		if p.peekTokenIs(token.ItemColon) || p.peekTokenIs(token.ItemArrow) {
			p.previousToken(i)
			return p.parseAnonymousFunction()
		}

		// otherwise it's a vector
		p.previousToken(i + 1)
		return p.parseVector()
	}

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.ItemRightParen) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.ItemLeftParen) {
		return nil
	}

	p.nextToken()

	expression.Condition = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.ItemRightParen) {
		p.nextToken()
	}

	if !p.expectPeek(token.ItemLeftCurly) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ItemElse) {
		p.nextToken()

		// else if: the nested if becomes the sole statement
		// of the else block, i.e.: else { if (...) {...} }
		if p.peekTokenIs(token.ItemIf) {
			p.nextToken()
			tok := p.curToken
			nested := p.parseIfExpression()
			if nested == nil {
				return nil
			}

			expression.Alternative = &ast.BlockStatement{
				Token: tok,
				Statements: []ast.Statement{
					&ast.ExpressionStatement{Token: tok, Expression: nested},
				},
			}

			return expression
		}

		if !p.expectPeek(token.ItemLeftCurly) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.ItemRightCurly) && !p.curTokenIs(token.ItemEOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	// reached the end of the file without a closing }
	// we report it on the opening { so it points at the unclosed block
	if !p.curTokenIs(token.ItemRightCurly) {
		p.errors = append(
			p.errors,
			diagnostics.NewError(block.Token, "missing closing `}`"),
		)
	}

	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	// it's a method
	if p.peekTokenIs(token.ItemLeftParen) {
		// skip func & left paren
		p.nextToken()
		p.nextToken()

		lit.MethodVariable = p.curToken.Value

		if !p.expectPeek(token.ItemColon) {
			return nil
		}

		// get type
		p.nextToken()

		lit.Method = &ast.Type{Name: p.curToken.Value, List: false}

		p.nextToken()
	}

	if !p.expectPeek(token.ItemIdent) {
		return nil
	}

	lit.Name = p.curToken.Value
	lit.NameToken = p.curToken

	lit.Operator = "="

	if !p.expectPeek(token.ItemLeftParen) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.ItemColon) {
		return nil
	}

	// parse types
	lit.ReturnType = p.parseTypes()

	// we could be in @generic which does not expect a body
	if !p.peekTokenIs(token.ItemLeftCurly) {
		return lit
	}

	p.nextToken()

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Parameter {
	parameters := []*ast.Parameter{}

	// function has no parameters
	if p.peekTokenIs(token.ItemRightParen) {
		p.nextToken()
		return parameters
	}

	p.skipNewLine()

	for p.peekTokenIs(token.ItemIdent) || p.peekTokenIs(token.ItemThreeDot) {
		p.nextToken()

		if p.curTokenIs(token.ItemComma) {
			p.nextToken()
		}

		parameter := &ast.Parameter{Token: p.curToken, Name: p.curToken.Value}

		if !p.expectPeek(token.ItemColon) {
			continue
		}

		// parse types
		parameter.Type = p.parseTypes()

		// if we have an assign we parse a statement, the function default
		if p.peekTokenIs(token.ItemAssign) {
			p.nextToken()
			p.nextToken()
			parameter.Operator = "="
			parameter.Default = p.parseExpressionStatement()
		}

		parameters = append(parameters, parameter)

		p.skipNewLine()

		if p.peekTokenIs(token.ItemComma) {
			p.nextToken()
		}

		p.skipNewLine()
	}

	p.skipNewLine()

	if !p.expectPeek(token.ItemRightParen) {
		return nil
	}

	return parameters
}

// parseIndexExpression parses x[i], x[i, j], x[, j], x[[i]]
// on entry the current token is the opening [ or [[
// on exit the current token is the closing ] or ]]
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{
		Token:  p.curToken,
		Left:   left,
		Double: p.curTokenIs(token.ItemDoubleLeftSquare),
	}

	closing := token.ItemRightSquare
	if exp.Double {
		closing = token.ItemDoubleRightSquare
	}

	p.nextToken()
	p.skipCurrentNewLine()

	for !p.curTokenIs(closing) {
		if p.curTokenIs(token.ItemEOF) || p.curTokenIs(token.ItemError) {
			p.errors = append(
				p.errors,
				diagnostics.NewError(exp.Token, fmt.Sprintf("missing closing `%v`", exp.Close())),
			)
			return nil
		}

		// empty slot, e.g.: x[, 1]
		var idx ast.Expression
		if !p.curTokenIs(token.ItemComma) {
			idx = p.parseExpression(LOWEST)
			p.nextToken()
			p.skipCurrentNewLine()
		}

		exp.Indices = append(exp.Indices, idx)

		if p.curTokenIs(token.ItemComma) {
			p.nextToken()
			p.skipCurrentNewLine()

			// trailing empty slot, e.g.: x[1, ]
			if p.curTokenIs(closing) {
				exp.Indices = append(exp.Indices, nil)
			}
			continue
		}

		if !p.curTokenIs(closing) {
			p.errors = append(
				p.errors,
				diagnostics.NewError(
					p.curToken,
					fmt.Sprintf("expected `%v` or `,`, got `%v` instead", exp.Close(), p.curToken.Class),
				),
			)
			return nil
		}
	}

	return exp
}

func (p *Parser) skipCurrentNewLine() {
	for p.curTokenIs(token.ItemNewLine) {
		p.nextToken()
	}
}

func (p *Parser) parseDecoratorGeneric() ast.Expression {
	dec := &ast.DecoratorGeneric{
		Token: p.curToken,
	}

	p.skipNewLine()
	p.nextToken()

	dec.Func = p.parseFunctionLiteral()

	return dec
}

func (p *Parser) parseDecoratorDefault() ast.Expression {
	dec := &ast.DecoratorDefault{
		Token: p.curToken,
	}

	p.skipNewLine()
	p.nextToken()

	dec.Func = p.parseFunctionLiteral()

	return dec
}

func (p *Parser) parseDecoratorEnvironment() ast.Expression {
	dec := &ast.DecoratorEnvironment{
		Token: p.curToken,
	}

	if !p.expectPeek(token.ItemLeftParen) {
		return nil
	}

	dec.Arguments = p.parseCallArguments()

	p.endDecoratorArguments()

	if !p.expectPeek(token.ItemNewLine) {
		return nil
	}

	if !p.expectPeek(token.ItemTypesDecl) {
		return nil
	}

	dec.Type = p.parseTypeDeclaration()

	return dec
}

func (p *Parser) parseDecoratorFactor() ast.Expression {
	dec := &ast.DecoratorFactor{
		Token: p.curToken,
	}

	if !p.expectPeek(token.ItemLeftParen) {
		return nil
	}

	dec.Arguments = p.parseCallArguments()

	p.endDecoratorArguments()

	if !p.expectPeek(token.ItemNewLine) {
		return nil
	}

	if !p.expectPeek(token.ItemTypesDecl) {
		return nil
	}

	dec.Type = p.parseTypeDeclaration()

	return dec
}

func (p *Parser) parseDecoratorMatrix() ast.Expression {
	dec := &ast.DecoratorMatrix{
		Token: p.curToken,
	}

	if !p.expectPeek(token.ItemLeftParen) {
		return nil
	}

	dec.Arguments = p.parseCallArguments()

	p.endDecoratorArguments()

	if !p.expectPeek(token.ItemNewLine) {
		return nil
	}

	if !p.expectPeek(token.ItemTypesDecl) {
		return nil
	}

	dec.Type = p.parseTypeDeclaration()

	return dec
}

// endDecoratorArguments moves to the closing parenthesis of decorator
// arguments, parseCallArguments may stop before it, on it,
// or on the closing parenthesis of the last argument, e.g.: c("a")
func (p *Parser) endDecoratorArguments() {
	for p.peekTokenIs(token.ItemRightParen) {
		p.nextToken()
	}
}

func (p *Parser) parseDecoratorClass() ast.Expression {
	dec := &ast.DecoratorClass{
		Token: p.curToken,
	}

	if !p.expectPeek(token.ItemLeftParen) {
		return nil
	}

	for p.peekTokenIs(token.ItemIdent) {
		p.nextToken()
		dec.Classes = append(dec.Classes, p.curToken.Value)
		p.nextToken()
	}

	if !p.expectPeek(token.ItemNewLine) {
		return nil
	}

	if !p.expectPeek(token.ItemTypesDecl) {
		return nil
	}

	dec.Type = p.parseTypeDeclaration()

	return dec
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function.Item().Value}

	switch f := function.(type) {
	case *ast.Identifier:
		exp.Name = f.Value
	}

	p.skipNewLine()
	exp.Arguments = p.parseCallArguments()

	// skip last token
	p.nextToken()

	if p.peekTokenIs(token.ItemRightParen) {
		p.nextToken()
	}

	// if it's a nested call we may have a trailing
	if p.peekTokenIs(token.ItemComma) {
		p.nextToken()
	}

	return exp
}

func (p *Parser) parseCallArguments() []ast.Argument {
	args := []ast.Argument{}

	if p.peekTokenIs(token.ItemRightParen) {
		return args
	}

	for !p.peekTokenIs(token.ItemRightParen) && !p.peekTokenIs(token.ItemComma) {
		if p.peekTokenIs(token.ItemEOF) {
			p.missingClosing(")")
			return args
		}
		p.nextToken()

		var arg ast.Argument
		if p.peekTokenIs(token.ItemAssign) {
			arg.Name = p.curToken.Value
		}
		arg.Token = p.curToken

		arg.Value = p.parseExpression(LOWEST)

		args = append(args, arg)

		if p.curTokenIs(token.ItemRightParen) {
			return args
		}

		if p.peekTokenIs(token.ItemComma) || p.peekTokenIs(token.ItemNewLine) {
			p.nextToken()
		}
		p.skipNewLine()
	}

	return args
}

func (p *Parser) registerPrefix(tokenType token.ItemType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.ItemType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) skipNewLine() {
	for p.peekTokenIs(token.ItemNewLine) {
		p.nextToken()
	}
}

func (p *Parser) parseTypes() ast.Types {
	var t ast.Types

	for p.peekTokenIs(token.ItemTypes) || p.peekTokenIs(token.ItemTypesList) ||
		p.peekTokenIs(token.ItemOr) || p.peekTokenIs(token.ItemTypesPkg) {

		p.nextToken()

		if p.curTokenIs(token.ItemOr) {
			continue
		}

		if p.curTokenIs(token.ItemTypesList) {
			continue
		}

		pkg := ""
		if p.curTokenIs(token.ItemTypesPkg) {
			pkg = p.curToken.Value
			// skip namespace
			p.nextToken()
			p.nextToken()
		}

		// is list
		p.previousToken(1)
		tok := p.curToken
		p.nextToken()

		list := false
		if tok.Class == token.ItemTypesList {
			list = true
		}

		t = append(t, &ast.Type{Name: p.curToken.Value, List: list, Package: pkg})
	}
	return t
}
