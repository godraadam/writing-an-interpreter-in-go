package eval

import (
	"fmt"
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

func error(format string, a ...any) *object.ErrorObj {
	return &object.ErrorObj{Message: fmt.Sprintf(format, a...)}
}

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node.Stmts)
	case *ast.ExprStmt:
		return Eval(node.Expr)

	// Expressions
	case *ast.NumberLiteral:
		return &object.Number{Value: node.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToMonkeyBool(node.Value)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.PrefixExpr:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalPrefixExpr(node.Operator, right)
	case *ast.InfixExpr:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}
		right := Eval(node.Right)
		if isError(right) {
			return right
		}
		return evalInfixExpr(node.Operator, left, right)
	case *ast.BlockStmt:
		return evalBlockStmt(node)
	case *ast.IfExpr:
		return evalIfExpr(node)
	case *ast.ReturnStmt:
		val := Eval(node.Value)
		if isError(val) {
			return val
		}
		return &object.ReturnObj{Value: val}
	}

	return nil
}

func evalProgram(stmts []ast.Stmt) object.Object {
	var result object.Object
	for _, stmt := range stmts {
		result = Eval(stmt)
		switch result := result.(type) {
		case *object.ReturnObj:
			return result.Value
		case *object.ErrorObj:
			return result
		}
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
		return error("Unknown operator %s %s", op, operand.Type())
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
		return error("Unknown operator %s %s", "-", operand.Type())
	}
	value := operand.(*object.Number).Value
	return &object.Number{Value: -value}
}

func evalInfixExpr(op string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() != right.Type():
		return error("Type mismatch %s %s %s", left.Type(), op, right.Type())
	case left.Type() == object.NUMBER_OBJ:
		return evalNumberInfixExpr(op, left, right)
	case op == "==":
		return nativeBoolToMonkeyBool(left == right)
	case op == "!=":
		return nativeBoolToMonkeyBool(left != right)
	default:
		return error("Unknown operator %s %s %s", left.Type(), op, right.Type())
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
		return error("Unknown operator %s %s %s", _left.Type(), op, _right.Type())
	}
}

func evalIfExpr(ie *ast.IfExpr) object.Object {
	cond := Eval(ie.Condition)
	if isError(cond) {
		return cond
	}
	if truthy(cond) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval((ie.Alternative))
	}
	return NIL
}

func evalBlockStmt(bs *ast.BlockStmt) object.Object {
	var result object.Object

	for _, stmt := range bs.Stmts {
		result = Eval(stmt)
		if result != nil && (result.Type() == object.RETURN_OBJ || result.Type() == object.ERROR_OBJ) {
			return result
		}
	}
	return result
}

func truthy(obj object.Object) bool {
	switch obj {
	case NIL:
		return false
	case FALSE:
		return false
	default:
		if obj.Type() == object.NUMBER_OBJ {
			return obj.(*object.Number).Value != 0
		}
	}
	return true
}

func isError(obj object.Object) bool {
	return obj != nil && obj.Type() == object.ERROR_OBJ
}
