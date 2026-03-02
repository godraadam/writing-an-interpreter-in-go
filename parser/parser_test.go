package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      any
	}{
		{"let x = 5;", "x", 5.0},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.Parse()
		checkParserErrors(t, p)

		if len(program.Stmts) != 1 {
			t.Fatalf("program.Stmts does not contain 1 statements. got=%d",
				len(program.Stmts))
		}

		stmt := program.Stmts[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStmt).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue any
	}{
		{"return 5;", 5.0},
		{"return true;", true},
		{"return foobar;", "foobar"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.Parse()
		checkParserErrors(t, p)

		if len(program.Stmts) != 1 {
			t.Fatalf("program.Stmts does not contain 1 statements. got=%d",
				len(program.Stmts))
		}

		stmt := program.Stmts[0]
		returnStmt, ok := stmt.(*ast.ReturnStmt)
		if !ok {
			t.Fatalf("stmt not *ast.ReturnStmt. got=%T", stmt)
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}
		if testLiteralExpression(t, returnStmt.Value, tt.expectedValue) {
			return
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Stmts) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Stmts))
	}
	stmt, ok := program.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("program.Stmts[0] is not ast.ExprStatement. got=%T",
			program.Stmts[0])
	}

	ident, ok := stmt.Expr.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expr)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar",
			ident.TokenLiteral())
	}
}

func TestNumberLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Stmts) != 1 {
		fmt.Print(program.Stmts)
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Stmts))
	}
	stmt, ok := program.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("program.Stmts[0] is not ast.ExprStmt. got=%T",
			program.Stmts[0])
	}

	literal, ok := stmt.Expr.(*ast.NumberLiteral)
	if !ok {
		t.Fatalf("exp not *ast.NumberLiteral. got=%T", stmt.Expr)
	}
	if literal.Value != 5.0 {
		t.Errorf("literal.Value not %f. got=%f", 5.0, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5",
			literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    any
	}{
		{"!5;", "!", 5.0},
		{"-15;", "-", 15.0},
		{"!foobar;", "!", "foobar"},
		{"-foobar;", "-", "foobar"},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.Parse()
		checkParserErrors(t, p)

		if len(program.Stmts) != 1 {
			t.Fatalf("program.Stmts does not contain %d statements. got=%d\n",
				1, len(program.Stmts))
		}

		stmt, ok := program.Stmts[0].(*ast.ExprStmt)
		if !ok {
			t.Fatalf("program.Stmts[0] is not ast.ExprStatement. got=%T",
				program.Stmts[0])
		}

		exp, ok := stmt.Expr.(*ast.PrefixExpr)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpr. got=%T", stmt.Expr)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  any
		operator   string
		rightValue any
	}{
		{"5 + 5;", 5.0, "+", 5.0},
		{"5 - 5;", 5.0, "-", 5.0},
		{"5 * 5;", 5.0, "*", 5.0},
		{"5 / 5;", 5.0, "/", 5.0},
		{"5 > 5;", 5.0, ">", 5.0},
		{"5 < 5;", 5.0, "<", 5.0},
		{"5 == 5;", 5.0, "==", 5.0},
		{"5 != 5;", 5.0, "!=", 5.0},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.Parse()
		checkParserErrors(t, p)

		if len(program.Stmts) != 1 {
			t.Fatalf("program.Stmts does not contain %d statements. got=%d\n",
				1, len(program.Stmts))
		}

		stmt, ok := program.Stmts[0].(*ast.ExprStmt)
		if !ok {
			t.Fatalf("program.Stmts[0] is not ast.ExprStmt. got=%T",
				program.Stmts[0])
		}

		if !testInfixExpression(t, stmt.Expr, tt.leftValue,
			tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5.5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5.5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
	}

	for i, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.Parse()
		checkParserErrors(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("Test [%d]: expected=%q, got=%q", i, tt.expected, actual)
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.Parse()
		checkParserErrors(t, p)

		if len(program.Stmts) != 1 {
			t.Fatalf("program has not enough statements. got=%d",
				len(program.Stmts))
		}

		stmt, ok := program.Stmts[0].(*ast.ExprStmt)
		if !ok {
			t.Fatalf("program.Stmts[0] is not ast.ExprStatement. got=%T",
				program.Stmts[0])
		}

		boolean, ok := stmt.Expr.(*ast.BooleanLiteral)
		if !ok {
			t.Fatalf("exp not *ast.BooleanLiteral. got=%T", stmt.Expr)
		}
		if boolean.Value != tt.expectedBoolean {
			t.Errorf("boolean.Value not %t. got=%t", tt.expectedBoolean,
				boolean.Value)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Stmts) != 1 {
		t.Fatalf("program.Stmts does not contain %d statements. got=%d\n",
			1, len(program.Stmts))
	}

	stmt, ok := program.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("program.Stmts[0] is not ast.ExprStmt. got=%T",
			program.Stmts[0])
	}

	exp, ok := stmt.Expr.(*ast.IfExpr)
	if !ok {
		t.Fatalf("stmt.Expr is not ast.Expr. got=%T",
			stmt.Expr)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Stmts) != 1 {
		t.Log(exp.Consequence.Stmts)
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Stmts))
	}

	consequence, ok := exp.Consequence.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("Stmts[0] is not ast.ExprStmt. got=%T",
			exp.Consequence.Stmts[0])
	}

	if !testIdentifier(t, consequence.Expr, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Stmts was not nil. got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Stmts) != 1 {
		t.Log(program.Stmts)
		t.Fatalf("program.Stmts does not contain %d statements. got=%d\n",
			1, len(program.Stmts))
	}

	stmt, ok := program.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("program.Stmts[0] is not ast.ExprStmt. got=%T",
			program.Stmts[0])
	}

	exp, ok := stmt.Expr.(*ast.IfExpr)
	if !ok {
		t.Fatalf("stmt.Expr is not ast.IfExpr. got=%T", stmt.Expr)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Stmts) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Stmts))
	}

	consequence, ok := exp.Consequence.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("Stmts[0] is not ast.ExprStmt. got=%T",
			exp.Consequence.Stmts[0])
	}

	if !testIdentifier(t, consequence.Expr, "x") {
		return
	}

	if len(exp.Alternative.Stmts) != 1 {
		t.Errorf("exp.Alternative.Stmts does not contain 1 statements. got=%d\n",
			len(exp.Alternative.Stmts))
	}

	alternative, ok := exp.Alternative.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("Stmts[0] is not ast.ExprStmt. got=%T",
			exp.Alternative.Stmts[0])
	}

	if !testIdentifier(t, alternative.Expr, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Stmts) != 1 {
		t.Fatalf("program.Stmts does not contain %d statements. got=%d\n",
			1, len(program.Stmts))
	}

	stmt, ok := program.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("program.Stmts[0] is not ast.ExprStatement. got=%T",
			program.Stmts[0])
	}

	function, ok := stmt.Expr.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expr is not ast.FunctionLiteral. got=%T",
			stmt.Expr)
	}

	if len(function.Params) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
			len(function.Params))
	}

	testLiteralExpression(t, function.Params[0], "x")
	testLiteralExpression(t, function.Params[1], "y")

	if len(function.Body.Stmts) != 1 {
		t.Fatalf("function.Body.Stmts has not 1 statements. got=%d\n",
			len(function.Body.Stmts))
	}

	bodyStmt, ok := function.Body.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("function.Body.Stmt is not ast.ExprStmt. got=%T",
			function.Body.Stmts[0])
	}

	testInfixExpression(t, bodyStmt.Expr, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.Parse()
		checkParserErrors(t, p)

		stmt := program.Stmts[0].(*ast.ExprStmt)
		function := stmt.Expr.(*ast.FunctionLiteral)

		if len(function.Params) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n",
				len(tt.expectedParams), len(function.Params))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Params[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := New(l)
	program := p.Parse()
	checkParserErrors(t, p)

	if len(program.Stmts) != 1 {
		t.Fatalf("program.Stmts does not contain %d statements. got=%d\n",
			1, len(program.Stmts))
	}

	stmt, ok := program.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("stmt is not ast.ExprStmt. got=%T",
			program.Stmts[0])
	}

	exp, ok := stmt.Expr.(*ast.CallExpr)
	if !ok {
		t.Fatalf("stmt.Expr is not ast.CallExpr. got=%T",
			stmt.Expr)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Args) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Args))
	}

	testLiteralExpression(t, exp.Args[0], 1)
	testInfixExpression(t, exp.Args[1], 2, "*", 3)
	testInfixExpression(t, exp.Args[2], 4, "+", 5)
}

func TestCallExpressionParameterParsing(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectedArgs  []string
	}{
		{
			input:         "add();",
			expectedIdent: "add",
			expectedArgs:  []string{},
		},
		{
			input:         "add(1);",
			expectedIdent: "add",
			expectedArgs:  []string{"1"},
		},
		{
			input:         "add(1, 2 * 3, 4 + 5);",
			expectedIdent: "add",
			expectedArgs:  []string{"1", "(2 * 3)", "(4 + 5)"},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.Parse()
		checkParserErrors(t, p)

		stmt := program.Stmts[0].(*ast.ExprStmt)
		exp, ok := stmt.Expr.(*ast.CallExpr)
		if !ok {
			t.Fatalf("stmt.Expr is not ast.CallExpr. got=%T",
				stmt.Expr)
		}

		if !testIdentifier(t, exp.Function, tt.expectedIdent) {
			return
		}

		if len(exp.Args) != len(tt.expectedArgs) {
			t.Fatalf("wrong number of arguments. want=%d, got=%d",
				len(tt.expectedArgs), len(exp.Args))
		}

		for i, arg := range tt.expectedArgs {
			if exp.Args[i].String() != arg {
				t.Errorf("argument %d wrong. want=%q, got=%q", i,
					arg, exp.Args[i].String())
			}
		}
	}
}

func testLetStatement(t *testing.T, s ast.Stmt, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStmt)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	ident, ok := letStmt.Target.(*ast.Identifier)

	if !ok {
		return false
	}
	if ident.Value != name {
		t.Errorf("letStmt.Target.Value not '%s'. got=%s", name, ident.Value)
		return false
	}

	if ident.TokenLiteral() != name {
		t.Errorf("letStmt.Target.TokenLiteral() not '%s'. got=%s",
			name, ident.TokenLiteral())
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, exp ast.Expr, left any,
	operator string, right any) bool {

	opExp, ok := exp.(*ast.InfixExpr)
	if !ok {
		t.Errorf("expr is not ast.InfixExpr. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("expr.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testLiteralExpression(
	t *testing.T,
	expr ast.Expr,
	expected any,
) bool {
	switch v := expected.(type) {
	case float64:
		return testNumberLiteral(t, expr, v)
	case float32:
		return testNumberLiteral(t, expr, float64(v))
	case int32:
		return testNumberLiteral(t, expr, float64(v))
	case int64:
		return testNumberLiteral(t, expr, float64(v))
	case int:
		return testNumberLiteral(t, expr, float64(v))
	case string:
		return testIdentifier(t, expr, v)
	case bool:
		return testBooleanLiteral(t, expr, v)
	}
	t.Errorf("type of expr not handled. got=%T", expr)
	return false
}

func testNumberLiteral(t *testing.T, il ast.Expr, value float64) bool {
	num, ok := il.(*ast.NumberLiteral)
	if !ok {
		t.Errorf("il not *ast.NumberLiteral. got=%T", il)
		return false
	}

	if num.Value != value {
		t.Errorf("number.Value not %f. got=%f", value, num.Value)
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expr, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expr, value bool) bool {
	bo, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("exp not *ast.BooleanLiteral. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s",
			value, bo.TokenLiteral())
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
