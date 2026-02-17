package object

import "fmt"

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

const (
	NUMBER_OBJ  = "NUMBER"
	BOOLEAN_OBJ = "BOOLEAN"
	NIL_OBJ     = "NULL"
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

type Nil struct{}

func (b *Nil) Inspect() string {
	return "nil"
}

func (n *Nil) Type() ObjectType {
	return NIL_OBJ
}
