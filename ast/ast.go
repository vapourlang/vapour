package ast

import (
	"bytes"
	"strings"

	"github.com/devOpifex/vapour/token"
)

type Node interface {
	TokenLiteral() string
	String() string
	Item() token.Item
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

func (p *Program) Item() token.Item { return token.Item{} }

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
}

func (ls *LetStatement) Item() token.Item     { return ls.Token }
func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Value }
func (ls *LetStatement) String() string {
	if ls.Value == nil {
		return ""
	}

	var out bytes.Buffer

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
}

func (cs *ConstStatement) Item() token.Item     { return cs.Token }
func (cs *ConstStatement) statementNode()       {}
func (cs *ConstStatement) TokenLiteral() string { return cs.Token.Value }
func (cs *ConstStatement) String() string {
	if cs.Value == nil {
		return ""
	}

	var out bytes.Buffer

	out.WriteString(cs.Name.String())
	out.WriteString(" = ")

	if cs.Value != nil {
		out.WriteString(cs.Value.String())
	}

	return out.String()
}

type RightSquare struct {
	Token token.Item
	Value string
}

func (rs *RightSquare) Item() token.Item     { return rs.Token }
func (rs *RightSquare) statementNode()       {}
func (rs *RightSquare) TokenLiteral() string { return rs.Token.Value }
func (rs *RightSquare) String() string {
	var out bytes.Buffer
	out.WriteString("]")
	return out.String()
}

type DoubleRightSquare struct {
	Token token.Item
	Value string
}

func (ds *DoubleRightSquare) Item() token.Item     { return ds.Token }
func (ds *DoubleRightSquare) statementNode()       {}
func (ds *DoubleRightSquare) TokenLiteral() string { return ds.Token.Value }
func (ds *DoubleRightSquare) String() string {
	var out bytes.Buffer
	out.WriteString("]]")
	return out.String()
}

type Type struct {
	Name string
	List bool
}

type TypeStatement struct {
	Token      token.Item  // type token
	Name       *Identifier // custom type
	Type       []*Type
	Attributes []*TypeAttributesStatement
	List       bool
}

func (ts *TypeStatement) Item() token.Item     { return ts.Token }
func (ts *TypeStatement) statementNode()       {}
func (ts *TypeStatement) TokenLiteral() string { return ts.Token.Value }
func (ts *TypeStatement) String() string {
	var out bytes.Buffer

	out.WriteString("# type ")
	out.WriteString(ts.Name.String() + " ")
	for _, v := range ts.Type {
		out.WriteString(v.Name + " ")
	}
	out.WriteString("\n")
	for _, v := range ts.Attributes {
		out.WriteString(v.String())
	}
	out.WriteString("\n")

	return out.String()
}

type TypeAttributesStatement struct {
	Token token.Item // type token
	Name  *Identifier
	Type  []*Type
}

func (ta *TypeAttributesStatement) Item() token.Item     { return ta.Token }
func (ta *TypeAttributesStatement) statementNode()       {}
func (ta *TypeAttributesStatement) TokenLiteral() string { return ta.Token.Value }
func (ta *TypeAttributesStatement) String() string {
	var out bytes.Buffer

	out.WriteString("# attribute ")
	out.WriteString(ta.Name.String() + ": ")
	for _, v := range ta.Type {
		out.WriteString(v.Name + " ")
	}
	out.WriteString("\n")

	return out.String()
}

type CommentStatement struct {
	Token token.Item
	Name  *Identifier
	Value string
}

func (c *CommentStatement) Item() token.Item     { return c.Token }
func (c *CommentStatement) statementNode()       {}
func (c *CommentStatement) TokenLiteral() string { return c.Token.Value }
func (c *CommentStatement) String() string {
	var out bytes.Buffer

	out.WriteString(c.TokenLiteral() + "\n")

	return out.String()
}

type NewLine struct {
	Token token.Item
}

func (nl *NewLine) Item() token.Item     { return nl.Token }
func (nl *NewLine) statementNode()       {}
func (nl *NewLine) TokenLiteral() string { return nl.Token.Value }
func (nl *NewLine) String() string {
	var out bytes.Buffer

	out.WriteString("\n")

	return out.String()
}

type DeferStatement struct {
	Token token.Item
	Func  Expression
}

func (ds *DeferStatement) Item() token.Item     { return ds.Token }
func (ds *DeferStatement) statementNode()       {}
func (ds *DeferStatement) TokenLiteral() string { return ds.Token.Value }
func (ds *DeferStatement) String() string {
	var out bytes.Buffer

	out.WriteString("\non.exit((")

	if ds.Func != nil {
		out.WriteString(ds.Func.String())
	}

	out.WriteString(")())")

	return out.String()
}

type ReturnStatement struct {
	Token       token.Item
	ReturnValue Expression
}

func (rs *ReturnStatement) Item() token.Item     { return rs.Token }
func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Value }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString("\n" + rs.TokenLiteral() + "(")

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

func (es *ExpressionStatement) Item() token.Item     { return es.Token }
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

func (bs *BlockStatement) Item() token.Item     { return bs.Token }
func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Value }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// Expressions
type Identifier struct {
	Token token.Item // the token.IDENT token
	Type  []*Type
	Value string
}

func (i *Identifier) Item() token.Item     { return i.Token }
func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Value }
func (i *Identifier) String() string {
	return i.Value
}

type Attrbute struct {
	Token token.Item // the token.IDENT token
	Value string
}

func (a *Attrbute) Item() token.Item     { return a.Token }
func (a *Attrbute) expressionNode()      {}
func (a *Attrbute) TokenLiteral() string { return a.Token.Value }
func (a *Attrbute) String() string {
	return a.Value
}

type Square struct {
	Token      token.Item
	Statements []Statement
	Type       []*Type
}

func (s *Square) Item() token.Item     { return s.Token }
func (s *Square) expressionNode()      {}
func (s *Square) TokenLiteral() string { return s.Token.Value }
func (s *Square) String() string {
	var out bytes.Buffer

	out.WriteString(s.Token.Value)

	if s.Statements != nil {
		for i, ss := range s.Statements {
			out.WriteString(ss.String())
			if i < len(s.Statements)-1 {
				out.WriteString(", ")
			}
		}
	}

	if s.Token.Value == "[" {
		out.WriteString("]")
	}

	if s.Token.Value == "[[" {
		out.WriteString("]]")
	}

	return out.String()
}

type Boolean struct {
	Token token.Item
	Value bool
	Type  []*Type
}

func (b *Boolean) Item() token.Item     { return b.Token }
func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Value }
func (b *Boolean) String() string {
	if b.Value {
		return "TRUE"
	}

	return "FALSE"
}

type IntegerLiteral struct {
	Token token.Item
	Value string
	Type  []*Type
}

func (il *IntegerLiteral) Item() token.Item     { return il.Token }
func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Value }
func (il *IntegerLiteral) String() string       { return il.Token.Value }

type FloatLiteral struct {
	Token token.Item
	Value string
	Type  []*Type
}

func (fl *FloatLiteral) Item() token.Item     { return fl.Token }
func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Value }
func (fl *FloatLiteral) String() string       { return fl.Token.Value }

type VectorLiteral struct {
	Token token.Item
	Value []Expression
}

func (v *VectorLiteral) Item() token.Item     { return v.Token }
func (v *VectorLiteral) expressionNode()      {}
func (v *VectorLiteral) TokenLiteral() string { return v.Token.Value }
func (v *VectorLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("c(")
	for i, e := range v.Value {
		if e == nil {
			continue
		}
		out.WriteString(e.String())
		if i < len(v.Value)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString(")\n")

	return out.String()
}

type SquareRightLiteral struct {
	Token token.Item
	Value string
}

func (s *SquareRightLiteral) Item() token.Item     { return s.Token }
func (s *SquareRightLiteral) expressionNode()      {}
func (s *SquareRightLiteral) TokenLiteral() string { return s.Token.Value }
func (s *SquareRightLiteral) String() string {
	var out bytes.Buffer

	out.WriteString(s.Value)
	out.WriteString("\n")

	return out.String()
}

type For struct {
	Token  token.Item
	Name   *LetStatement
	Vector Expression
	Value  *BlockStatement
}

func (f *For) Item() token.Item     { return f.Token }
func (f *For) expressionNode()      {}
func (f *For) TokenLiteral() string { return f.Token.Value }
func (f *For) String() string {
	var out bytes.Buffer

	out.WriteString("for(")
	out.WriteString(f.Name.String())
	out.WriteString(" in ")
	out.WriteString(f.Vector.String())
	out.WriteString(")\n {")
	out.WriteString(f.Value.String())
	out.WriteString("}\n")

	return out.String()
}

type While struct {
	Token     token.Item
	Statement Statement
	Value     *BlockStatement
}

func (w *While) Item() token.Item     { return w.Token }
func (w *While) expressionNode()      {}
func (w *While) TokenLiteral() string { return w.Token.Value }
func (w *While) String() string {
	var out bytes.Buffer

	out.WriteString("while(")
	out.WriteString(w.Statement.String())
	out.WriteString(")\n {")
	out.WriteString(w.Value.String())
	out.WriteString("}\n")

	return out.String()
}

type Null struct {
	Token token.Item
	Value string
	Type  []*Type
}

func (n *Null) Item() token.Item     { return n.Token }
func (n *Null) expressionNode()      {}
func (n *Null) TokenLiteral() string { return n.Token.Value }
func (n *Null) String() string {
	var out bytes.Buffer

	out.WriteString(n.Value)

	return out.String()
}

type Keyword struct {
	Token token.Item
	Value string
	Type  []*Type
}

func (kw *Keyword) Item() token.Item     { return kw.Token }
func (kw *Keyword) expressionNode()      {}
func (kw *Keyword) TokenLiteral() string { return kw.Token.Value }
func (kw *Keyword) String() string {
	var out bytes.Buffer

	out.WriteString(kw.Value)

	return out.String()
}

type StringLiteral struct {
	Token token.Item
	Str   string
	Type  []*Type
}

func (sl *StringLiteral) Item() token.Item     { return sl.Token }
func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Value }
func (sl *StringLiteral) String() string {
	var out bytes.Buffer

	out.WriteString(sl.Token.Value + sl.Str + sl.Token.Value)

	return out.String()
}

type PrefixExpression struct {
	Token    token.Item // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) Item() token.Item     { return pe.Token }
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

func (ie *InfixExpression) Item() token.Item     { return ie.Token }
func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Value }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString(ie.Left.String())
	if ie.Operator != "::" && ie.Operator != "$" && ie.Operator != ".." {
		out.WriteString(" ")
	}

	if ie.Operator != ".." {
		out.WriteString(ie.Operator)
	}

	if ie.Operator == ".." {
		out.WriteString(":")
	}

	if ie.Operator != "::" && ie.Operator != "$" && ie.Operator != ".." {
		out.WriteString(" ")
	}

	if ie.Right != nil {
		out.WriteString(ie.Right.String())
	}

	return out.String()
}

type IfExpression struct {
	Token       token.Item // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) Item() token.Item     { return ie.Token }
func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Value }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if(")
	out.WriteString(ie.Condition.String())
	out.WriteString("){\n")
	out.WriteString(ie.Consequence.String())
	out.WriteString("}")

	if ie.Alternative != nil {
		out.WriteString(" else {\n")
		out.WriteString(ie.Alternative.String())
		out.WriteString("\n}\n")
	}

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Item // The 'func' token
	Method     string
	Name       *Identifier
	Operator   string
	Type       []*Type
	Parameters []*Parameter
	Body       *BlockStatement
}

func (fl *FunctionLiteral) Item() token.Item     { return fl.Token }
func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Value }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	if fl.Name.String() != "" {
		out.WriteString("#' @yield ")
		for i, v := range fl.Type {
			out.WriteString(v.Name)
			if i < len(fl.Type)-1 {
				out.WriteString(" | ")
			}
			out.WriteString("\n")
		}
		out.WriteString(fl.Name.String())
	}

	if fl.Method != "" {
		out.WriteString("." + fl.Method)
	}

	if fl.Operator != "" {
		out.WriteString(" " + fl.Operator + " ")
	}

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
	Type     []*Type
	Default  *ExpressionStatement
	Method   bool
}

func (p *Parameter) Item() token.Item     { return p.Token }
func (p *Parameter) expressionNode()      {}
func (p *Parameter) TokenLiteral() string { return p.Token.Value }
func (p *Parameter) String() string {
	var out bytes.Buffer

	out.WriteString(p.Name)
	if p.Operator != "" {
		out.WriteString(" " + p.Operator + " " + p.Default.String())
	}
	return out.String()
}

type Argument struct {
	Token token.Item
	Name  string
	Value Expression
}

type CallExpression struct {
	Token     token.Item // The '(' token
	Function  Expression // Identifier or FunctionLiteral
	Name      string
	Arguments []Argument
}

func (ce *CallExpression) Item() token.Item     { return ce.Token }
func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Value }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		if a.Value != nil {
			args = append(args, a.Value.String())
		}
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type Decorator struct {
	Token   token.Item // The '@' token
	Name    string
	Classes []string
	Type    *TypeStatement
}

func (d *Decorator) Item() token.Item     { return d.Token }
func (d *Decorator) expressionNode()      {}
func (d *Decorator) TokenLiteral() string { return d.Token.Value }
func (d *Decorator) String() string {
	var out bytes.Buffer

	out.WriteString("# classes: ")
	out.WriteString(strings.Join(d.Classes, ","))
	out.WriteString("\n")
	out.WriteString(d.Type.String())

	return out.String()
}
