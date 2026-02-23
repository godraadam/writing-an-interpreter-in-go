package object

import (
	"bytes"
	"fmt"
	"monkey/ast"
	"strings"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

const (
	NUMBER_OBJ   = "NUMBER"
	BOOLEAN_OBJ  = "BOOLEAN"
	STR_OBJ      = "STRING"
	RETURN_OBJ   = "RETURN"
	NIL_OBJ      = "NULL"
	ERROR_OBJ    = "ERROR"
	FUNCTION_OBJ = "FUNCTION"
)

type Number struct {
	Value float64
}

func (n *Number) Inspect() string {
	return fmt.Sprintf("%f", n.Value)
}

func (n *Number) Type() ObjectType {
	return NUMBER_OBJ
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (n *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

type String struct {
	Value string
}

func (s *String) Inspect() string {
	return s.Value
}

func (n *String) Type() ObjectType {
	return STR_OBJ
}

type Nil struct{}

func (b *Nil) Inspect() string {
	return "nil"
}

func (n *Nil) Type() ObjectType {
	return NIL_OBJ
}

type ReturnObj struct {
	Value Object
}

func (obj *ReturnObj) Inspect() string {
	return obj.Value.Inspect()
}

func (obj *ReturnObj) Type() ObjectType {
	return RETURN_OBJ
}

type ErrorObj struct {
	Message string
}

func (err *ErrorObj) Inspect() string {
	return "ERROR: " + err.Message
}

func (err *ErrorObj) Type() ObjectType {
	return ERROR_OBJ
}

type FunctionObj struct {
	Params []*ast.Identifier
	Body   *ast.BlockStmt
	Env    *Environment
}

func (f *FunctionObj) Inspect() string {
	var out bytes.Buffer

	params := []string{}

	for _, p := range f.Params {
		params = append(params, p.String())
	}

	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

func (err *FunctionObj) Type() ObjectType {
	return FUNCTION_OBJ
}
