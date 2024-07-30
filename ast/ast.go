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
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Value }
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString("#' @type " + ls.Name.String() + " ")
	for i, v := range ls.Name.Type {
		if v.List {
			out.WriteString("[]")
		}
		out.WriteString(v.Name)
		if i < len(ls.Name.Type)-1 {
			out.WriteString(" | ")
		}
	}
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
}

func (cs *ConstStatement) statementNode()       {}
func (cs *ConstStatement) TokenLiteral() string { return cs.Token.Value }
func (cs *ConstStatement) String() string {
	var out bytes.Buffer

	out.WriteString("#' @type " + cs.Name.String() + " ")
	for i, v := range cs.Name.Type {
		if v.List {
			out.WriteString("[]")
		}
		out.WriteString(v.Name)
		if i < len(cs.Name.Type)-1 {
			out.WriteString(" | ")
		}
	}
	out.WriteString("\n")
	out.WriteString(cs.Name.String())
	out.WriteString(" = ")

	if cs.Value != nil {
		out.WriteString(cs.Value.String())
	}

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
	Object     []*Type
	Attributes []*TypeAttributesStatement
	List       bool
}

func (ts *TypeStatement) statementNode()       {}
func (ts *TypeStatement) TokenLiteral() string { return ts.Token.Value }
func (ts *TypeStatement) String() string {
	var out bytes.Buffer

	out.WriteString("# type ")
	out.WriteString(ts.Name.String() + " ")
	if len(ts.Object) > 0 {
		for _, v := range ts.Type {
			out.WriteString(v.Name + " ")
		}

		for _, v := range ts.Object {
			out.WriteString(v.Name + " ")
		}
	} else {
		for _, v := range ts.Type {
			out.WriteString(v.Name + " ")
		}
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

func (c *CommentStatement) statementNode()       {}
func (c *CommentStatement) TokenLiteral() string { return c.Token.Value }
func (c *CommentStatement) String() string {
	var out bytes.Buffer

	out.WriteString(c.TokenLiteral() + "\n")

	return out.String()
}

type SemiColon struct {
	Token token.Item
}

func (s *SemiColon) statementNode()       {}
func (s *SemiColon) TokenLiteral() string { return s.Token.Value }
func (s *SemiColon) String() string {
	var out bytes.Buffer

	out.WriteString(";")

	return out.String()
}

type NewLine struct {
	Token token.Item
}

func (nl *NewLine) statementNode()       {}
func (nl *NewLine) TokenLiteral() string { return nl.Token.Value }
func (nl *NewLine) String() string {
	var out bytes.Buffer

	out.WriteString("\n")

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
		return es.Expression.String() + "\n"
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
	Type  []*Type
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Value }
func (i *Identifier) String() string {
	var out bytes.Buffer
	out.WriteString(i.Value)
	return out.String()
}

type Boolean struct {
	Token token.Item
	Value bool
	Type  []*Type
}

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

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Value }
func (il *IntegerLiteral) String() string       { return il.Token.Value }

type VectorLiteral struct {
	Token token.Item
	Value []Expression
}

func (v *VectorLiteral) expressionNode()      {}
func (v *VectorLiteral) TokenLiteral() string { return v.Token.Value }
func (v *VectorLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("c(")
	for i, e := range v.Value {
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

func (s *SquareRightLiteral) expressionNode()      {}
func (s *SquareRightLiteral) TokenLiteral() string { return s.Token.Value }
func (s *SquareRightLiteral) String() string {
	var out bytes.Buffer

	out.WriteString(s.Value)
	out.WriteString("\n")

	return out.String()
}

type For struct {
	Token     token.Item
	Statement Statement
	Value     *BlockStatement
}

func (f *For) expressionNode()      {}
func (f *For) TokenLiteral() string { return f.Token.Value }
func (f *For) String() string {
	var out bytes.Buffer

	out.WriteString("for(")
	out.WriteString(f.Statement.String())
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
	out.WriteString(ie.Operator)

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
		}
		out.WriteString("\n")
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
	Default  *Identifier
}

func (p *Parameter) statementNode()       {}
func (p *Parameter) TokenLiteral() string { return p.Token.Value }
func (p *Parameter) String() string {
	var out bytes.Buffer

	out.WriteString(p.Name)
	if p.Operator != "" {
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
