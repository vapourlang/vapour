package walker

import (
	"fmt"

	"github.com/vapourlang/vapour/ast"
	"github.com/vapourlang/vapour/environment"
)

type Function struct {
	name       string
	returnType ast.Types
	arguments  []ast.Types
}

func (w *Walker) typesExist(types ast.Types) (*ast.Type, bool) {
	for _, t := range types {
		// bit hacky but NA, NULL, etc. have a blank type
		// need to fix upstream in parser
		if t.Name == "" {
			return nil, true
		}

		_, te := w.env.GetType(t.Package, t.Name)
		_, fe := w.env.GetSignature(t.Name)

		if !te && !fe {
			return t, false
		}
	}

	return nil, true
}

func (w *Walker) allTypesIdentical(types []*ast.Type) bool {
	if len(types) == 0 {
		return true
	}

	types, ok := w.getNativeTypes(types)

	if !ok {
		return false
	}

	var previousType *ast.Type
	for i, t := range types {
		if i == 0 {
			previousType = t
			continue
		}

		if t.Name != previousType.Name || t.List != previousType.List {
			return false
		}
	}

	return true
}

func typeIdentical(t1, t2 *ast.Type) bool {
	return t1.Name == t2.Name && t1.List == t2.List
}

func acceptAny(types ast.Types) bool {
	for _, t := range types {
		if t.Name == "any" {
			return true
		}
	}
	return false
}

func (w *Walker) typesValid(valid, actual ast.Types) bool {
	validNative, _ := w.getNativeTypes(valid)
	actualNative, _ := w.getNativeTypes(actual)

	validNative = append(validNative, valid...)
	actualNative = append(actualNative, actual...)

	// we don't have the type
	if len(validNative) == 0 {
		return true
	}

	if acceptAny(validNative) {
		return true
	}

	for _, l := range actualNative {
		if w.typeValid(l, validNative) {
			continue
		}

		return false
	}

	return true
}

func (w *Walker) typeValid(t *ast.Type, valid ast.Types) bool {
	// we just don't have the type
	// could be base R dataset
	if t.Name == "" {
		return true
	}

	for _, v := range valid {
		if typeIdentical(t, v) {
			return true
		}

		if v.Name == "num" && t.Name == "int" && v.List == t.List {
			return true
		}
	}

	return false
}

func (w *Walker) validAccessType(types ast.Types) bool {
	for _, t := range types {
		if t.Name == "any" {
			return true
		}

		obj, exists := w.env.GetType(t.Package, t.Name)

		// we don't have the type
		// assume it's an error on our end?
		if !exists {
			return true
		}

		if !contains(obj.Object, []string{"dataframe", "object", "struct", "environment"}) {
			return false
		}
	}
	return true
}

func (w *Walker) validMathTypes(types ast.Types) bool {
	types, ok := w.getNativeTypes(types)

	if !ok {
		return false
	}

	for _, t := range types {
		if !contains(t.Name, []string{"int", "num", "na", "date", "posixct", "posixlt"}) {
			return false
		}
	}
	return true
}

func contains(value string, arr []string) bool {
	for _, a := range arr {
		if value == a {
			return true
		}
	}
	return false
}

func (w *Walker) retrieveNativeTypes(types, nativeTypes ast.Types) (ast.Types, bool) {
	for _, t := range types {
		if environment.IsNativeType(t.Name) {
			nativeTypes = append(nativeTypes, t)
			continue
		}

		customType, exists := w.env.GetType(t.Package, t.Name)

		if customType.Object == "struct" || customType.Object == "object" {
			return append(nativeTypes, t), false
		}

		if exists && customType.Object == "vector" {
			return w.retrieveNativeTypes(customType.Type, nativeTypes)
		}

		if exists && customType.Object == "list" {
			return w.retrieveNativeTypes(customType.Type, nativeTypes)
		}

		if exists && customType.Object == "impliedList" {
			return w.retrieveNativeTypes(customType.Type, nativeTypes)
		}

		return append(nativeTypes, t), false
	}

	return nativeTypes, true
}

func (w *Walker) getNativeTypes(types ast.Types) (ast.Types, bool) {
	return w.retrieveNativeTypes(types, ast.Types{})
}

func (w *Walker) validIteratorTypes(types ast.Types) bool {
	types, _ = w.getNativeTypes(types)
	var valid []bool
	for _, t := range types {
		if contains(t.Name, []string{"int", "num", "char", "any"}) {
			valid = append(valid, true)
			continue
		}

		custom, exists := w.env.GetType(t.Package, t.Name)
		if !exists {
			valid = append(valid, false)
			continue
		}

		if len(custom.Type) > 0 && allLists(custom.Type) {
			valid = append(valid, true)
			continue
		}

		if custom.Object == "list" {
			valid = append(valid, true)
			continue
		}

		valid = append(valid, false)
	}
	return allTrue(valid)
}

func allLists(types ast.Types) bool {
	for _, t := range types {
		if !t.List {
			return false
		}
	}

	return true
}

func allTrue(values []bool) bool {
	for _, b := range values {
		if !b {
			return false
		}
	}
	return true
}

func (w *Walker) checkIdentifier(node *ast.Identifier) {
	v, exists := w.env.GetVariable(node.Value, true)

	if exists {
		if v.CanMiss {
			w.addWarnf(
				node.Token,
				"`%v` might be missing",
				node.Token.Value,
			)
		}

		return
	}

	_, exists = w.env.GetType("", node.Value)

	if exists {
		w.env.SetTypeUsed("", node.Value)
		return
	}

	_, exists = w.env.GetFunction(node.Value, true)

	if exists {
		// we currently don't set/check used functions
		// they may be exported from a package (and thus not used)
		return
	}

	_, exists = w.env.GetSignature(node.Value)

	if exists {
		w.env.SetSignatureUsed(node.Value)
		return
	}

	// we are actually declaring variable in a call
	if !w.isInNamespace() {
		w.addWarnf(
			node.Token,
			"`%v` not found",
			node.Value,
		)
	}
}

type callback func(*ast.Identifier)

func (w *Walker) callIfIdentifier(node ast.Node, fn callback) {
	switch n := node.(type) {
	case *ast.Identifier:
		fn(n)
	}
}

func (w *Walker) checkIfIdentifier(node ast.Node) {
	switch n := node.(type) {
	case *ast.Identifier:
		w.checkIdentifier(n)
	}
}

func (w *Walker) getAttribute(name string, attrs []*ast.TypeAttributesStatement) (ast.Types, bool) {
	for _, a := range attrs {
		if a.Name == name {
			return a.Type, true
		}
	}
	return nil, false
}

func (w *Walker) attributeMatch(arg ast.Argument, inc ast.Types, t environment.Type) bool {
	a, ok := w.getAttribute(arg.Name, t.Attributes)

	if !ok {
		w.addFatalf(
			arg.Token,
			"attribute `%v` not found",
			arg.Name,
		)
		return false
	}

	ok = w.typesValid(a, inc)

	if !ok {
		w.addFatalf(
			arg.Token,
			"attribute `%v` expects `%v`, got `%v`",
			arg.Name,
			a,
			inc,
		)
		return false
	}

	return true
}

func (w *Walker) warnUnusedTypes() {
	for k, v := range w.env.Types() {
		if v.Used {
			continue
		}
		w.addInfof(
			v.Token,
			"type `%v` is never used",
			k,
		)
	}
}

func (w *Walker) warnUnusedVariables() {
	for k, v := range w.env.Variables() {
		if v.Used {
			continue
		}
		w.addInfof(
			v.Token,
			"`%v` is never used",
			k,
		)
	}
}

func (w *Walker) canBeFunction(types ast.Types) (environment.Signature, bool) {
	for _, t := range types {
		signature, exists := w.env.GetSignature(t.Name)

		if exists {
			return signature, true
		}
	}
	return environment.Signature{}, false
}

func (w *Walker) getFunctionSignatureFromFunctionLiteral(fn *ast.FunctionLiteral) Function {
	fnc := Function{
		name:       fn.Name,
		returnType: fn.ReturnType,
	}

	for _, a := range fn.Parameters {
		fnc.arguments = append(fnc.arguments, a.Type)
	}

	return fnc
}

func (w *Walker) getFunctionSignatureFromIdentifier(n *ast.Identifier) Function {
	fnc := Function{
		name: n.Value,
	}

	fn, exists := w.env.GetFunction(n.Value, true)

	if !exists {
		return fnc
	}

	return w.getFunctionSignatureFromFunctionLiteral(fn.Value)
}

func (w *Walker) signaturesMatch(valid environment.Signature, actual Function) (string, bool) {
	ok := w.typesValid(valid.Value.Return, actual.returnType)

	if !ok {
		return fmt.Sprintf("expected return `%v`, got `%v`", valid.Value.Return, actual.returnType), false
	}

	if len(valid.Value.Arguments) != len(actual.arguments) {
		return fmt.Sprintf("number of arguments do not match signature `%v`, expected %v, got %v", valid.Value.Name, len(valid.Value.Arguments), len(actual.arguments)), false
	}

	for index, validArg := range valid.Value.Arguments {
		actualArg := actual.arguments[index]

		ok := w.typesValid(validArg, actualArg)

		if !ok {
			return fmt.Sprintf("argument #%v, expected `%v`, got `%v`", index, validArg, actualArg), false
		}
	}

	return "", true
}

func (w *Walker) comparisonsValid(valid, actual ast.Types) bool {
	validNative, _ := w.getNativeTypes(valid)
	actualNative, _ := w.getNativeTypes(actual)

	validNative = append(validNative, valid...)
	actualNative = append(actualNative, actual...)

	// we don't have the type
	if len(validNative) == 0 {
		return true
	}

	if acceptAny(validNative) {
		return true
	}

	for _, l := range actualNative {
		if w.comparisonValid(l, validNative) {
			continue
		}

		return false
	}

	return true
}

func (w *Walker) comparisonValid(t *ast.Type, valid ast.Types) bool {
	// we just don't have the type
	// could be base R dataset
	if t.Name == "" {
		return true
	}

	for _, v := range valid {
		if typeIdentical(t, v) {
			return true
		}

		if v.Name == "num" && t.Name == "int" && v.List == t.List {
			return true
		}

		if v.Name == "int" && t.Name == "num" && v.List == t.List {
			return true
		}
	}

	return false
}
