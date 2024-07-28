package parser

import (
	"fmt"

	"github.com/devOpifex/vapour/ast"
	"github.com/devOpifex/vapour/lexer"
	"github.com/devOpifex/vapour/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // call(X)
	INDEX
)

var precedences = map[token.ItemType]int{
	token.ItemAssign:           EQUALS,
	token.ItemDoubleEqual:      EQUALS,
	token.ItemNotEqual:         EQUALS,
	token.ItemLessThan:         LESSGREATER,
	token.ItemGreaterThan:      LESSGREATER,
	token.ItemPlus:             SUM,
	token.ItemMinus:            SUM,
	token.ItemDivide:           PRODUCT,
	token.ItemMultiply:         PRODUCT,
	token.ItemPipe:             INDEX,
	token.ItemInfix:            PRODUCT,
	token.ItemLeftParen:        CALL,
	token.ItemLeftCurly:        INDEX,
	token.ItemSemiColon:        INDEX,
	token.ItemNewLine:          INDEX,
	token.ItemDot:              SUM,
	token.ItemLeftSquare:       SUM,
	token.ItemDoubleLeftSquare: SUM,
	token.ItemIn:               EQUALS,
	token.ItemRange:            EQUALS,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	pos int

	curToken  token.Item
	peekToken token.Item

	prefixParseFns map[token.ItemType]prefixParseFn
	infixParseFns  map[token.ItemType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.ItemType]prefixParseFn)
	p.registerPrefix(token.ItemIdent, p.parseIdentifier)
	p.registerPrefix(token.ItemInteger, p.parseIntegerLiteral)
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
	p.registerPrefix(token.ItemSemiColon, p.parseIdentifier)
	p.registerPrefix(token.ItemNewLine, p.parseIdentifier)
	p.registerPrefix(token.ItemRightSquare, p.parseIdentifier)
	p.registerPrefix(token.ItemDoubleRightSquare, p.parseIdentifier)
	p.registerPrefix(token.ItemFor, p.parseFor)

	p.infixParseFns = make(map[token.ItemType]infixParseFn)
	p.registerInfix(token.ItemPlus, p.parseInfixExpression)
	p.registerInfix(token.ItemMinus, p.parseInfixExpression)
	p.registerInfix(token.ItemDivide, p.parseInfixExpression)
	p.registerInfix(token.ItemMultiply, p.parseInfixExpression)
	p.registerInfix(token.ItemAssign, p.parseInfixExpression)
	p.registerInfix(token.ItemDoubleEqual, p.parseInfixExpression)
	p.registerInfix(token.ItemLessThan, p.parseInfixExpression)
	p.registerInfix(token.ItemGreaterThan, p.parseInfixExpression)
	p.registerInfix(token.ItemPipe, p.parseInfixExpression)
	p.registerInfix(token.ItemLeftSquare, p.parseInfixExpression)
	p.registerInfix(token.ItemDoubleLeftSquare, p.parseInfixExpression)
	p.registerInfix(token.ItemComma, p.parseInfixExpression)
	p.registerInfix(token.ItemDot, p.parseInfixExpression)
	p.registerInfix(token.ItemIn, p.parseInfixExpression)
	p.registerInfix(token.ItemRange, p.parseInfixExpression)

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

func (p *Parser) previousToken(n int) {
	for i := 0; i < n; i++ {
		p.pos--
		p.curToken = p.l.Items[p.pos-2]
		p.peekToken = p.l.Items[p.pos-1]
	}
}

func (p *Parser) print() {
	fmt.Println("++++++++++++++++++++++++++\nCurrent")
	p.curToken.Print()
	fmt.Println("Peek")
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

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.ItemType) {
	msg := fmt.Sprintf("expected next token to be %c, got %c instead",
		t, p.peekToken.Class)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.ItemType) {
	msg := fmt.Sprintf("no prefix parse function for %c found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Run() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.ItemEOF) {
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
	case token.ItemComment:
		return p.parseCommentStatement()
	case token.ItemNewLine:
		return p.parseNewLine()
	case token.ItemSemiColon:
		return p.parseSemiColon()
	case token.ItemTypesDecl:
		return p.parseTypeDeclaration()
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

	p.nextToken()

	lit.Statement = p.parseExpression(LOWEST)

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

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Value}
}

func (p *Parser) parseNull() ast.Expression {
	return &ast.Null{
		Token: p.curToken,
		Value: "NULL",
		Type:  []*ast.Type{{Name: "null"}},
	}
}

func (p *Parser) parseElipsis() ast.Expression {
	return &ast.Keyword{Token: p.curToken, Value: "..."}
}

func (p *Parser) parseNA() ast.Expression {
	return &ast.Keyword{
		Token: p.curToken,
		Value: "NA",
		Type:  []*ast.Type{{Name: "na"}},
	}
}

func (p *Parser) parseDot() ast.Expression {
	return &ast.Keyword{Token: p.curToken, Value: "$"}
}

func (p *Parser) parseNan() ast.Expression {
	return &ast.Keyword{
		Token: p.curToken,
		Value: "NaN",
		Type:  []*ast.Type{{Name: "null"}},
	}
}

func (p *Parser) parseNaString() ast.Expression {
	return &ast.Keyword{
		Token: p.curToken,
		Value: "NA_character_",
		Type:  []*ast.Type{{Name: "na_char"}},
	}
}

func (p *Parser) parseNaReal() ast.Expression {
	return &ast.Keyword{
		Token: p.curToken,
		Value: "NA_real_",
		Type:  []*ast.Type{{Name: "na_real"}},
	}
}

func (p *Parser) parseNaComplex() ast.Expression {
	return &ast.Keyword{
		Token: p.curToken,
		Value: "NA_complex_",
		Type:  []*ast.Type{{Name: "na_complex"}},
	}
}

func (p *Parser) parseNaInteger() ast.Expression {
	return &ast.Keyword{
		Token: p.curToken,
		Value: "NA_integer_",
		Type:  []*ast.Type{{Name: "na_int"}},
	}
}

func (p *Parser) parseInf() ast.Expression {
	return &ast.Keyword{
		Token: p.curToken,
		Value: "Inf",
		Type:  []*ast.Type{{Name: "inf"}},
	}
}

func (p *Parser) parseTypeDeclaration() *ast.TypeStatement {
	typ := &ast.TypeStatement{Token: p.curToken}

	// expect the custom type
	if !p.peekTokenIs(token.ItemTypes) {
		return nil
	}

	p.nextToken()

	typ.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Value}

	// expect colon
	if !p.peekTokenIs(token.ItemColon) {
		return nil
	}

	p.nextToken()

	// may also expect list []
	if p.peekTokenIs(token.ItemTypesList) {
		typ.List = true
		p.nextToken()
	}

	// expect native type
	if !p.peekTokenIs(token.ItemTypes) {
		return nil
	}

	for p.expectPeek(token.ItemOr) || p.expectPeek(token.ItemTypes) {
		if p.peekTokenIs(token.ItemOr) {
			continue
		}

		p.previousToken(1)
		tok := p.curToken
		p.nextToken()

		list := false
		if tok.Class == token.ItemTypesList {
			list = true
		}

		typ.Type = append(typ.Type, &ast.Type{Name: p.curToken.Value, List: list})
	}

	// no struct
	if !p.peekTokenIs(token.ItemLeftCurly) {
		return typ
	}

	// skip left curly {
	p.nextToken()

	// it's an object or dataframe with named attributes only
	if !p.peekTokenIs(token.ItemIdent) {
		typ.Attributes = p.parseTypeAttributes()
		return typ
	}

	p.nextToken()

	// it's a struct or list
	typ.Object = []*ast.Type{{Name: p.curToken.Value}}

	for p.expectPeek(token.ItemOr) || p.expectPeek(token.ItemIdent) {
		if p.curTokenIs(token.ItemOr) {
			continue
		}

		p.previousToken(1)
		tok := p.curToken
		p.nextToken()

		list := false
		if tok.Class == token.ItemOr {
			list = true
		}

		typ.Object = append(typ.Object, &ast.Type{Name: p.curToken.Value, List: list})
	}

	p.nextToken()

	typ.Attributes = p.parseTypeAttributes()

	return typ
}

func (p *Parser) parseTypeAttributes() []*ast.TypeAttributesStatement {
	var attrs []*ast.TypeAttributesStatement

	for !p.peekTokenIs(token.ItemRightCurly) {
		p.nextToken()
		attrs = append(attrs, p.parseTypeAttribute())
	}

	p.nextToken()

	return attrs
}

func (p *Parser) parseTypeAttribute() *ast.TypeAttributesStatement {
	tok := p.curToken

	ident := &ast.Identifier{
		Token: tok,
		Value: p.curToken.Value,
	}

	// skip colon
	p.nextToken()

	var types []*ast.Type

	for p.expectPeek(token.ItemTypes) || p.expectPeek(token.ItemTypesOr) {
		if p.curTokenIs(token.ItemOr) {
			continue
		}

		p.previousToken(1)
		tok := p.curToken
		p.nextToken()

		list := false
		if tok.Class == token.ItemOr {
			list = true
		}

		types = append(types, &ast.Type{Name: p.curToken.Value, List: list})
	}

	if p.peekTokenIs(token.ItemComma) {
		p.nextToken()
	}

	return &ast.TypeAttributesStatement{
		Token: tok,
		Name:  ident,
		Type:  types,
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.ItemIdent) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Value}

	if !p.expectPeek(token.ItemColon) {
		return nil
	}

	for p.expectPeek(token.ItemTypes) ||
		p.expectPeek(token.ItemTypesList) ||
		p.expectPeek(token.ItemTypesOr) {
		if p.curTokenIs(token.ItemTypesOr) {
			continue
		}

		if p.curTokenIs(token.ItemTypesList) {
			continue
		}

		p.previousToken(1)
		tok := p.curToken
		p.nextToken()

		list := false
		if tok.Class == token.ItemTypesList {
			list = true
		}

		stmt.Name.Type = append(stmt.Name.Type, &ast.Type{Name: p.curToken.Value, List: list})
	}

	if !p.expectPeek(token.ItemAssign) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.ItemNewLine) || p.peekTokenIs(token.ItemSemiColon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseConstStatement() *ast.ConstStatement {
	stmt := &ast.ConstStatement{Token: p.curToken}

	if !p.expectPeek(token.ItemIdent) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Value}

	if !p.expectPeek(token.ItemColon) {
		return nil
	}

	if p.peekTokenIs(token.ItemTypesList) {
		p.nextToken()
	}

	if !p.peekTokenIs(token.ItemTypes) {
		return nil
	}

	for p.expectPeek(token.ItemTypes) ||
		p.expectPeek(token.ItemTypesList) ||
		p.expectPeek(token.ItemTypesOr) {
		if p.curTokenIs(token.ItemTypesOr) {
			continue
		}

		if p.curTokenIs(token.ItemTypesList) {
			continue
		}

		p.previousToken(1)
		tok := p.curToken
		p.nextToken()

		list := false
		if tok.Class == token.ItemTypesList {
			list = true
		}

		stmt.Name.Type = append(stmt.Name.Type, &ast.Type{Name: p.curToken.Value, List: list})
	}

	if !p.expectPeek(token.ItemAssign) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.ItemNewLine) || p.peekTokenIs(token.ItemSemiColon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.ItemNewLine) || p.peekTokenIs(token.ItemSemiColon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.ItemNewLine) ||
		p.peekTokenIs(token.ItemSemiColon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
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
		Type:  []*ast.Type{{Name: "int", List: false}},
	}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.curToken,
		Value: p.curToken.Value == "true",
		Type:  []*ast.Type{{Name: "bool", List: false}},
	}
}

func (p *Parser) parseCommentStatement() ast.Statement {
	return &ast.CommentStatement{Token: p.curToken, Value: p.curToken.Value}
}

func (p *Parser) parseSemiColon() ast.Statement {
	return &ast.SemiColon{Token: p.curToken}
}

func (p *Parser) parseNewLine() ast.Statement {
	return &ast.NewLine{Token: p.curToken}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	str := &ast.StringLiteral{
		Token: p.curToken,
		Type:  []*ast.Type{{Name: "char", List: false}},
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

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	operator := p.curToken.Value

	if operator == "." {
		operator = "$"
	}

	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: operator,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseVector() ast.Expression {
	vec := &ast.VectorLiteral{}

	for !p.peekTokenIs(token.ItemRightParen) {
		p.nextToken()
		if p.curTokenIs(token.ItemComma) {
			continue
		}
		vec.Value = append(vec.Value, p.parseExpression(LOWEST))
	}

	p.nextToken()

	return vec
}

func (p *Parser) parseAnonymousFunction() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	lit.Parameters = p.parseFunctionParameters()

	// parse types
	for p.expectPeek(token.ItemTypes) ||
		p.expectPeek(token.ItemTypesList) || p.expectPeek(token.ItemTypesOr) {
		if p.curTokenIs(token.ItemTypesOr) {
			continue
		}

		if p.curTokenIs(token.ItemTypesList) {
			continue
		}

		p.previousToken(1)
		tok := p.curToken
		p.nextToken()

		list := false
		if tok.Class == token.ItemTypesList {
			list = true
		}

		lit.Type = append(lit.Type, &ast.Type{Name: p.curToken.Value, List: list})
	}

	lit.Name = &ast.Identifier{Token: p.curToken, Value: ""}

	if !p.expectPeek(token.ItemArrow) {
		return nil
	}

	if !p.expectPeek(token.ItemLeftCurly) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.previousToken(1)
	tk := p.curToken
	p.nextToken()

	// skip paren left (
	p.nextToken()

	// no prefix
	// it's either a vector (1, 2, 3)
	// or an anonymous function
	// (x: string) string => print(x)
	// or
	// (x: string) string => { print(x) }
	if tk.Class != token.ItemIdent {
		i := 0
		for !p.curTokenIs(token.ItemRightParen) {
			i++
			p.nextToken()
		}

		// if the closing paren ) is followed by
		// a type or => it's an anonymous function
		if p.peekTokenIs(token.ItemTypes) || p.peekTokenIs(token.ItemArrow) {
			p.previousToken(i + 1)
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

	if !p.expectPeek(token.ItemRightParen) {
		return nil
	}

	if !p.expectPeek(token.ItemLeftCurly) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ItemElse) {
		p.nextToken()

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

	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	param := &ast.Parameter{}

	if p.expectPeek(token.ItemLeftParen) {
		// skip left paren
		p.nextToken()

		param = &ast.Parameter{
			Token: p.curToken,
			Name:  p.curToken.Value,
		}

		list := false
		if p.peekTokenIs(token.ItemTypesOr) {
			list = true
			p.nextToken()
		}

		// get type
		p.nextToken()

		lit.Method = p.curToken.Value
		param.Type = []*ast.Type{{Name: p.curToken.Value, List: list}}

		p.nextToken()
	}

	if !p.expectPeek(token.ItemIdent) {
		return nil
	}

	lit.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Value}

	lit.Operator = "="

	if !p.expectPeek(token.ItemLeftParen) {
		return nil
	}

	var params []*ast.Parameter
	otherParams := p.parseFunctionParameters()
	if param.Name != "" {
		params = append(params, param)
	}
	params = append(params, otherParams...)
	lit.Parameters = params

	// parse types
	for p.expectPeek(token.ItemTypes) ||
		p.expectPeek(token.ItemTypesList) || p.expectPeek(token.ItemTypesOr) {
		if p.curTokenIs(token.ItemTypesOr) {
			continue
		}

		if p.curTokenIs(token.ItemTypesList) {
			continue
		}

		p.previousToken(1)
		tok := p.curToken
		p.nextToken()

		list := false
		if tok.Class == token.ItemTypesList {
			list = true
		}

		lit.Type = append(lit.Type, &ast.Type{Name: p.curToken.Value, List: list})
	}

	if !p.expectPeek(token.ItemLeftCurly) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Parameter {
	parameters := []*ast.Parameter{}

	if p.peekTokenIs(token.ItemRightParen) {
		p.nextToken()
		return parameters
	}

	p.nextToken()

	for !p.peekTokenIs(token.ItemRightParen) {
		parameter := &ast.Parameter{Token: p.curToken, Name: p.curToken.Value}

		if !p.peekTokenIs(token.ItemColon) {
			continue
		}

		// skip colon
		p.nextToken()

		// parse types
		for p.expectPeek(token.ItemTypes) || p.expectPeek(token.ItemTypesList) || p.expectPeek(token.ItemTypesOr) {
			if p.curTokenIs(token.ItemTypesOr) {
				continue
			}

			if p.curTokenIs(token.ItemTypesList) {
				continue
			}

			p.previousToken(1)
			tok := p.curToken
			p.nextToken()

			list := false
			if tok.Class == token.ItemTypesList {
				list = true
			}

			parameter.Type = append(parameter.Type, &ast.Type{Name: p.curToken.Value, List: list})
		}

		if p.expectPeek(token.ItemAssign) {
			p.nextToken()
			parameter.Operator = "="
			parameter.Default = &ast.Identifier{Token: p.curToken, Value: p.curToken.Value}
		}

		parameters = append(parameters, parameter)

		if p.peekTokenIs(token.ItemComma) {
			p.nextToken()
			p.nextToken()
		}
	}

	if !p.expectPeek(token.ItemRightParen) {
		return nil
	}

	return parameters
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	p.nextToken()
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.peekTokenIs(token.ItemRightParen) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for !p.peekTokenIs(token.ItemRightParen) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	return args
}

func (p *Parser) registerPrefix(tokenType token.ItemType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.ItemType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
