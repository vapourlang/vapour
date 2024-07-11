package parser

import (
	"fmt"
	"strconv"

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
	CALL        // myFunction(X)
)

var precedences = map[token.ItemType]int{
	token.ItemAssign:      EQUALS,
	token.ItemDoubleEqual: EQUALS,
	token.ItemNotEqual:    EQUALS,
	token.ItemLessThan:    LESSGREATER,
	token.ItemGreaterThan: LESSGREATER,
	token.ItemPlus:        SUM,
	token.ItemMinus:       SUM,
	token.ItemDivide:      PRODUCT,
	token.ItemMultiply:    PRODUCT,
	token.ItemLeftParen:   CALL,
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
	p.registerPrefix(token.ItemBool, p.parseBoolean)
	p.registerPrefix(token.ItemLeftParen, p.parseGroupedExpression)
	p.registerPrefix(token.ItemIf, p.parseIfExpression)
	p.registerPrefix(token.ItemFunction, p.parseFunctionLiteral)
	p.registerPrefix(token.ItemDoubleQuote, p.parseStringLiteral)
	p.registerPrefix(token.ItemSingleQuote, p.parseStringLiteral)

	p.infixParseFns = make(map[token.ItemType]infixParseFn)
	p.registerInfix(token.ItemPlus, p.parseInfixExpression)
	p.registerInfix(token.ItemMinus, p.parseInfixExpression)
	p.registerInfix(token.ItemDivide, p.parseInfixExpression)
	p.registerInfix(token.ItemMultiply, p.parseInfixExpression)
	p.registerInfix(token.ItemAssign, p.parseInfixExpression)
	p.registerInfix(token.ItemAssign, p.parseInfixExpression)
	p.registerInfix(token.ItemDoubleEqual, p.parseInfixExpression)
	p.registerInfix(token.ItemLessThan, p.parseInfixExpression)
	p.registerInfix(token.ItemGreaterThan, p.parseInfixExpression)

	p.registerInfix(token.ItemLeftParen, p.parseCallExpression)

	p.nextToken()
	p.nextToken()

	return p
}

// TODO simplify: it should not be this complex
func (p *Parser) nextToken() {
	if p.pos >= len(p.l.Items) {
		p.curToken = token.Item{Class: token.ItemEOF}
		return
	}
	p.curToken = p.peekToken
	p.peekToken = p.l.Items[p.pos]
	p.pos++
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
	default:
		return p.parseExpressionStatement()
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

	if p.peekTokenIs(token.ItemTypesList) {
		p.nextToken()
	}

	if !p.expectPeek(token.ItemTypes) {
		return nil
	}

	stmt.Type = append(stmt.Type, p.curToken.Value)

	if !p.expectPeek(token.ItemAssign) {
		return nil
	}

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.ItemEOL) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseConstStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.ItemIdent) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Value}

	if !p.peekTokenIs(token.ItemAssign) {
		return nil
	}

	// skip identifier
	p.nextToken()
	// skip assignment
	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.ItemEOL) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.ItemEOL) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.ItemEOL) {
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

	for !p.peekTokenIs(token.ItemEOL) && precedence < p.peekPrecedence() {
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

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Value}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Value, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Value)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.ItemBool)}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	str := &ast.StringLiteral{Token: p.curToken}
	str.Str = []string{}

	p.nextToken()

	for !p.curTokenIs(token.ItemString) && !p.curTokenIs(token.ItemEOF) {
		str.Str = append(str.Str, p.curToken.Value)
		p.nextToken()
	}

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
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Value,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

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

	if !p.expectPeek(token.ItemLeftParen) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.ItemLeftCurly) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.ItemRightCurly) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Value}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.ItemEOL) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Value}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.ItemRightParen) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments()
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

	for p.peekTokenIs(token.ItemEOL) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.ItemEOL) {
		return nil
	}

	return args
}

func (p *Parser) registerPrefix(tokenType token.ItemType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.ItemType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
