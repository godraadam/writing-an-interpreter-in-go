package eval

import (
	"monkey/ast"
	"monkey/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NIL   = &object.Nil{}
)

func nativeBoolToMonkeyBool(value bool) object.Object {
	if value {
		return TRUE
	}
	return FALSE
}

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalStmts(node.Stmts)
	case *ast.ExprStmt:
		return Eval(node.Expr)

	// Expressions
	case *ast.NumberLiteral:
		return &object.Number{Value: node.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToMonkeyBool(node.Value)
	case *ast.PrefixExpr:
		right := Eval(node.Right)
		return evalPrefixExpr(node.Operator, right)
	case *ast.InfixExpr:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpr(node.Operator, left, right)
	}

	return nil
}

func evalStmts(stmts []ast.Stmt) object.Object {
	var result object.Object
	for _, stmt := range stmts {
		result = Eval(stmt)
	}
	return result
}

func evalPrefixExpr(op string, operand object.Object) object.Object {
	switch op {
	case "!":
		return evalBangOpExpr(operand)
	case "-":
		return evalMinusOpExpr(operand)
	default:
		return NIL
	}
}

func evalBangOpExpr(operand object.Object) object.Object {
	switch operand {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NIL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusOpExpr(operand object.Object) object.Object {
	if operand.Type() != object.NUMBER_OBJ {
		return NIL
	}
	value := operand.(*object.Number).Value
	return &object.Number{Value: -value}
}

func evalInfixExpr(op string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.NUMBER_OBJ && right.Type() == object.NUMBER_OBJ:
		return evalNumberInfixExpr(op, left, right)
	case op == "==":
		return nativeBoolToMonkeyBool(left == right)
	case op == "!=":
		return nativeBoolToMonkeyBool(left != right)
	default:
		return NIL
	}
}

func evalNumberInfixExpr(op string, _left object.Object, _right object.Object) object.Object {
	left := _left.(*object.Number).Value
	right := _right.(*object.Number).Value

	switch op {
	case "+":
		return &object.Number{Value: left + right}
	case "-":
		return &object.Number{Value: left - right}
	case "/":
		return &object.Number{Value: left / right}
	case "*":
		return &object.Number{Value: left * right}
	case "==":
		return nativeBoolToMonkeyBool(left == right)
	case "!=":
		return nativeBoolToMonkeyBool(left != right)
	case "<":
		return nativeBoolToMonkeyBool(left < right)
	case ">":
		return nativeBoolToMonkeyBool(left > right)
	case "<=":
		return nativeBoolToMonkeyBool(left <= right)
	case ">=":
		return nativeBoolToMonkeyBool(left >= right)
	default:
		return NIL
	}
}
