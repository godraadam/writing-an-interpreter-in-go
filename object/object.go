package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
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
	ARRAY_OBJ    = "ARRAY"
	MAP_OBJ      = "MAP"
)

type Number struct {
	Value float64
}

func (n *Number) Inspect() string {
	return fmt.Sprintf("%.17g", n.Value)
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
	var out bytes.Buffer
	out.WriteString("\"")
	out.WriteString(s.Value)
	out.WriteString("\"")
	return out.String()
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

type ArrayObj struct {
	Elements []Object
}

func (arr *ArrayObj) Inspect() string {
	var out bytes.Buffer

	out.WriteString("[")
	elements := []string{}
	for _, el := range arr.Elements {
		elements = append(elements, el.Inspect())
	}
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

func (arr *ArrayObj) Type() ObjectType {
	return ARRAY_OBJ
}

type MapPair struct {
	Key   Object
	Value Object
}

type MapObj struct {
	Pairs map[MapKey]MapPair
}

func (mo *MapObj) Inspect() string {
	var out bytes.Buffer

	out.WriteString("{")
	pairs := []string{}
	for _, pair := range mo.Pairs {
		pairs = append(pairs, pair.Key.Inspect()+": "+pair.Value.Inspect())
	}
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

func (mo *MapObj) Type() ObjectType {
	return MAP_OBJ
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

func (f *FunctionObj) Type() ObjectType {
	return FUNCTION_OBJ
}

type MapKey struct {
	Type  ObjectType
	Value float64
}

type Hashable interface {
	Hash() MapKey
}

func (b *Boolean) Hash() MapKey {
	if b.Value {
		return MapKey{Type: BOOLEAN_OBJ, Value: 1}
	}
	return MapKey{Type: BOOLEAN_OBJ, Value: 0}

}

func (n *Number) Hash() MapKey {
	return MapKey{Type: NUMBER_OBJ, Value: n.Value}
}

func (s *String) Hash() MapKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return MapKey{Type: STR_OBJ, Value: float64(h.Sum64())}
}
