package ast

import (
	"bytes"
	"strings"

	"github.com/devOpifex/vapour/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// Statements
type LetStatement struct {
	Token token.Item
	Name  *Identifier
	Value Expression
	Type  []string
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Value }
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString("# type: " + strings.Join(ls.Type, ", or "))
	out.WriteString("\n")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	return out.String()
}

type ConstStatement struct {
	Token token.Item
	Name  *Identifier
	Value Expression
	Type  []string
}

func (cs *ConstStatement) statementNode()       {}
func (cs *ConstStatement) TokenLiteral() string { return cs.Token.Value }
func (cs *ConstStatement) String() string {
	var out bytes.Buffer

	out.WriteString(cs.Name.String())
	out.WriteString(" = ")

	if cs.Value != nil {
		out.WriteString(cs.Value.String())
	}

	return out.String()
}

type ReturnStatement struct {
	Token       token.Item
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Value }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + "(")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(")")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Item // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Value }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type BlockStatement struct {
	Token      token.Item // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Value }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
		out.WriteString("\n")
	}

	return out.String()
}

// Expressions
type Identifier struct {
	Token token.Item // the token.IDENT token
	Type  []string
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Value }
func (i *Identifier) String() string       { return i.Value }

type Boolean struct {
	Token token.Item
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Value }
func (b *Boolean) String() string       { return b.Token.Value }

type IntegerLiteral struct {
	Token token.Item
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Value }
func (il *IntegerLiteral) String() string       { return il.Token.Value }

type StringLiteral struct {
	Token token.Item
	Str   []string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Value }
func (sl *StringLiteral) String() string {
	var out bytes.Buffer

	out.WriteString(sl.TokenLiteral())
	for _, s := range sl.Str {
		out.WriteString(s)
	}
	out.WriteString(sl.TokenLiteral())

	return out.String()
}

type PrefixExpression struct {
	Token    token.Item // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Value }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Item // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Value }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())

	return out.String()
}

type IfExpression struct {
	Token       token.Item // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Value }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Item // The 'func' token
	Name       *Identifier
	Operator   string
	Type       []string
	Parameters []*Parameter
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Value }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("# type: " + strings.Join(fl.Type, ", or") + "\n")
	out.WriteString(fl.Name.String() + " " + fl.Operator + " ")
	out.WriteString("function")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(fl.Body.String())
	out.WriteString("}\n")

	return out.String()
}

type Parameter struct {
	Token    token.Item // The 'func' token
	Name     string
	Operator string
	Type     []string
	Default  *Identifier
}

func (p *Parameter) statementNode()       {}
func (p *Parameter) TokenLiteral() string { return p.Token.Value }
func (p *Parameter) String() string {
	var out bytes.Buffer

	out.WriteString(p.Name)
	if p.Default.Value != "" {
		out.WriteString(" " + p.Operator + " " + p.Default.Value)
	}
	return out.String()
}

type CallExpression struct {
	Token     token.Item // The '(' token
	Function  Expression // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Value }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}
