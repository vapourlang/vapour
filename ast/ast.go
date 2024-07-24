package ast

import (
	"bytes"
	"strings"

	"github.com/devOpifex/vapour/token"
)

type Node interface {
	TokenLiteral() string
	Transpile() string
	check(*Environment) astErrors
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

func (p *Program) Transpile() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.Transpile())
	}

	return out.String()
}

func (p *Program) Check(env *Environment) astErrors {
	var errs astErrors

	for _, s := range p.Statements {
		err := s.check(env)
		errs = append(errs, err...)
	}

	return errs
}

// Statements
type LetStatement struct {
	Token token.Item
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) check(env *Environment) astErrors {
	var errs astErrors
	_, exists := env.GetVariable(ls.Name.Value)

	if exists {
		err := astError{Token: ls.Token, Message: ls.Name.Value + " is already declared"}
		errs = append(errs, err)
	}

	env.SetVariable(ls.Name.Value, ls)

	if ls.Name != nil {
		nameErr := ls.Name.check(env)
		errs = append(errs, nameErr...)
	}

	if ls.Value != nil {
		valueErr := ls.Value.check(env)
		errs = append(errs, valueErr...)
	}

	return errs
}
func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Value }
func (ls *LetStatement) Transpile() string {
	var out bytes.Buffer

	out.WriteString("#' @type " + ls.Name.Transpile() + " ")
	for i, v := range ls.Name.Type {
		out.WriteString(v.Name)
		if i < len(ls.Name.Type)-1 {
			out.WriteString(" | ")
		}
	}
	out.WriteString("\n")
	out.WriteString(ls.Name.Transpile())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.Transpile())
	}

	return out.String()
}

type ConstStatement struct {
	Token token.Item
	Name  *Identifier
	Value Expression
}

func (cs *ConstStatement) check(env *Environment) astErrors {
	var errs astErrors
	_, exists := env.GetVariable(cs.Name.Value)

	if exists {
		err := astError{Token: cs.Token, Message: cs.Name.Value + " is already declared"}
		errs = append(errs, err)
	}

	env.SetVariable(cs.Name.Value, cs)

	if cs.Name != nil {
		nameErr := cs.Name.check(env)
		errs = append(errs, nameErr...)
	}

	if cs.Value != nil {
		valueErr := cs.Value.check(env)
		errs = append(errs, valueErr...)
	}

	return errs
}
func (cs *ConstStatement) statementNode()       {}
func (cs *ConstStatement) TokenLiteral() string { return cs.Token.Value }
func (cs *ConstStatement) Transpile() string {
	var out bytes.Buffer

	out.WriteString("#' @type " + cs.Name.Transpile() + " ")
	for i, v := range cs.Name.Type {
		out.WriteString(v.Name)
		if i < len(cs.Name.Type)-1 {
			out.WriteString(" | ")
		}
	}
	out.WriteString("\n")
	out.WriteString(cs.Name.Transpile())
	out.WriteString(" = ")

	if cs.Value != nil {
		out.WriteString(cs.Value.Transpile())
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

func (ts *TypeStatement) check(env *Environment) astErrors {
	var errs astErrors

	if ts.Name == nil {
		return errs
	}

	_, exists := env.GetType(ts.Name.Value)

	if exists {
		err := astError{Token: ts.Token, Message: ts.Name.Value + " is already declared"}
		errs = append(errs, err)
	}

	env.SetType(ts.Name.Value, ts)
	return astErrors{}
}
func (ts *TypeStatement) statementNode()       {}
func (ts *TypeStatement) TokenLiteral() string { return ts.Token.Value }
func (ts *TypeStatement) Transpile() string {
	var out bytes.Buffer

	out.WriteString("# type ")
	out.WriteString(ts.Name.Transpile() + " ")
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
		out.WriteString(v.Transpile())
	}
	out.WriteString("\n")

	return out.String()
}

type TypeAttributesStatement struct {
	Token token.Item // type token
	Name  *Identifier
	Type  []*Type
}

func (ta *TypeAttributesStatement) check(env *Environment) astErrors {
	return astErrors{}
}
func (ta *TypeAttributesStatement) statementNode()       {}
func (ta *TypeAttributesStatement) TokenLiteral() string { return ta.Token.Value }
func (ta *TypeAttributesStatement) Transpile() string {
	var out bytes.Buffer

	out.WriteString("# attribute ")
	out.WriteString(ta.Name.Transpile() + ": ")
	for _, v := range ta.Type {
		out.WriteString(v.Name + " ")
	}
	out.WriteString("\n")

	return out.String()
}

type Keyword struct {
	Token token.Item
	Value Expression
}

func (kw *Keyword) check(env *Environment) astErrors {
	return astErrors{}
}
func (kw *Keyword) expressionNode()      {}
func (kw *Keyword) TokenLiteral() string { return kw.Token.Value }
func (kw *Keyword) Transpile() string {
	var out bytes.Buffer

	out.WriteString(kw.Value.Transpile())

	return out.String()
}

type CommentStatement struct {
	Token token.Item
	Name  *Identifier
	Value string
}

func (c *CommentStatement) check(env *Environment) astErrors {
	return astErrors{}
}
func (c *CommentStatement) statementNode()       {}
func (c *CommentStatement) TokenLiteral() string { return c.Token.Value }
func (c *CommentStatement) Transpile() string {
	var out bytes.Buffer

	out.WriteString(c.TokenLiteral() + "\n")

	return out.String()
}

type SemiColon struct {
	Token token.Item
}

func (s *SemiColon) check(env *Environment) astErrors {
	return astErrors{}
}
func (s *SemiColon) statementNode()       {}
func (s *SemiColon) TokenLiteral() string { return s.Token.Value }
func (s *SemiColon) Transpile() string {
	var out bytes.Buffer

	out.WriteString(";")

	return out.String()
}

type NewLine struct {
	Token token.Item
}

func (nl *NewLine) check(env *Environment) astErrors {
	return astErrors{}
}
func (nl *NewLine) statementNode()       {}
func (nl *NewLine) TokenLiteral() string { return nl.Token.Value }
func (nl *NewLine) Transpile() string {
	var out bytes.Buffer

	out.WriteString("\n")

	return out.String()
}

type ReturnStatement struct {
	Token       token.Item
	ReturnValue Expression
}

func (rs *ReturnStatement) check(env *Environment) astErrors {
	if rs.ReturnValue != nil {
		return astErrors{}
	}

	return rs.ReturnValue.check(env)
}
func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Value }
func (rs *ReturnStatement) Transpile() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + "(")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.Transpile())
	}

	out.WriteString(")")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Item // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) check(env *Environment) astErrors {
	if es.Expression != nil {
		return astErrors{}
	}
	return es.Expression.check(env)
}
func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Value }
func (es *ExpressionStatement) Transpile() string {
	if es.Expression != nil {
		return es.Expression.Transpile() + "\n"
	}
	return ""
}

type BlockStatement struct {
	Token      token.Item // the { token
	Statements []Statement
}

func (bs *BlockStatement) check(env *Environment) astErrors {
	env.addEnclosedEnvironment()

	var errs astErrors

	for _, s := range bs.Statements {
		err := s.check(env)
		errs = append(errs, err...)
	}

	return errs
}
func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Value }
func (bs *BlockStatement) Transpile() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.Transpile())
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

func (i *Identifier) check(env *Environment) astErrors {
	return astErrors{}
}
func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Value }
func (i *Identifier) Transpile() string {
	var out bytes.Buffer
	out.WriteString(i.Value)
	return out.String()
}

type Boolean struct {
	Token token.Item
	Value bool
}

func (b *Boolean) check(env *Environment) astErrors {
	return astErrors{}
}
func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Value }
func (b *Boolean) Transpile() string {
	if b.Value {
		return "TRUE"
	}

	return "FALSE"
}

type IntegerLiteral struct {
	Token token.Item
	Value string
}

func (il *IntegerLiteral) check(env *Environment) astErrors {
	return astErrors{}
}
func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Value }
func (il *IntegerLiteral) Transpile() string    { return il.Token.Value }

type VectorLiteral struct {
	Token token.Item
	Value []Expression
}

func (v *VectorLiteral) check(env *Environment) astErrors {
	var errs astErrors

	for _, e := range v.Value {
		err := e.check(env)
		errs = append(errs, err...)
	}

	return errs
}
func (v *VectorLiteral) expressionNode()      {}
func (v *VectorLiteral) TokenLiteral() string { return v.Token.Value }
func (v *VectorLiteral) Transpile() string {
	var out bytes.Buffer

	out.WriteString("c(")
	for i, e := range v.Value {
		out.WriteString(e.Transpile())
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

func (s *SquareRightLiteral) check(env *Environment) astErrors {
	return astErrors{}
}
func (s *SquareRightLiteral) expressionNode()      {}
func (s *SquareRightLiteral) TokenLiteral() string { return s.Token.Value }
func (s *SquareRightLiteral) Transpile() string {
	var out bytes.Buffer

	out.WriteString(s.Value)
	out.WriteString("\n")

	return out.String()
}

type StringLiteral struct {
	Token token.Item
	Str   string
}

func (sl *StringLiteral) check(env *Environment) astErrors {
	return astErrors{}
}
func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Value }
func (sl *StringLiteral) Transpile() string {
	var out bytes.Buffer

	out.WriteString(sl.Token.Value + sl.Str + sl.Token.Value)

	return out.String()
}

type PrefixExpression struct {
	Token    token.Item // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pw *PrefixExpression) check(env *Environment) astErrors {
	return pw.Right.check(env)
}
func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Value }
func (pe *PrefixExpression) Transpile() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.Transpile())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Item // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) check(env *Environment) astErrors {
	var errs astErrors

	lErr := ie.Left.check(env)
	errs = append(errs, lErr...)

	if ie.Right != nil {
		rErr := ie.Right.check(env)
		errs = append(errs, rErr...)
	}

	return errs
}
func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Value }
func (ie *InfixExpression) Transpile() string {
	var out bytes.Buffer

	out.WriteString(ie.Left.Transpile())
	out.WriteString(" " + ie.Operator + " ")

	if ie.Right != nil {
		out.WriteString(ie.Right.Transpile())
	}

	return out.String()
}

type IfExpression struct {
	Token       token.Item // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) check(env *Environment) astErrors {
	var errs astErrors

	condErr := ie.Condition.check(env)
	errs = append(errs, condErr...)

	consErr := ie.Consequence.check(env)
	errs = append(errs, consErr...)

	if ie.Alternative != nil {
		altErr := ie.Alternative.check(env)
		errs = append(errs, altErr...)
	}

	return errs
}
func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Value }
func (ie *IfExpression) Transpile() string {
	var out bytes.Buffer

	out.WriteString("if(")
	out.WriteString(ie.Condition.Transpile())
	out.WriteString("){\n")
	out.WriteString(ie.Consequence.Transpile())
	out.WriteString("}")

	if ie.Alternative != nil {
		out.WriteString(" else {\n")
		out.WriteString(ie.Alternative.Transpile())
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

func (fl *FunctionLiteral) check(env *Environment) astErrors {
	env.SetVariable(fl.Name.Value, fl)
	env.addEnclosedEnvironment()

	var errs astErrors

	bodyErr := fl.Body.check(env)
	errs = append(errs, bodyErr...)

	for _, p := range fl.Parameters {
		err := p.check(env)
		errs = append(errs, err...)
	}

	return errs
}
func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Value }
func (fl *FunctionLiteral) Transpile() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.Transpile())
	}

	if fl.Name.Transpile() != "" {
		out.WriteString("#' @yield ")
		for i, v := range fl.Type {
			out.WriteString(v.Name)
			if i < len(fl.Type)-1 {
				out.WriteString(" | ")
			}
		}
		out.WriteString("\n")
		out.WriteString(fl.Name.Transpile())
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
	out.WriteString(fl.Body.Transpile())
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

func (p *Parameter) check(env *Environment) astErrors {
	env.SetVariable(p.Name, p)
	return astErrors{}
}
func (p *Parameter) statementNode()       {}
func (p *Parameter) TokenLiteral() string { return p.Token.Value }
func (p *Parameter) Transpile() string {
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

func (ce *CallExpression) check(env *Environment) astErrors {
	var errs astErrors

	for _, a := range ce.Arguments {
		err := a.check(env)
		errs = append(errs, err...)
	}

	return errs
}
func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Value }
func (ce *CallExpression) Transpile() string {
	var out bytes.Buffer

	out.WriteString("\n")
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.Transpile())
	}

	out.WriteString(ce.Function.Transpile())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}
