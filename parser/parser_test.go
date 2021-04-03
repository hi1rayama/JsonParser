package parser

import (
	"fmt"
	"jsonParser/ast"
	"jsonParser/lexer"
	"testing"
)

// TestNumberValue 数値単体の構文解析の単体テスト
func TestNumberValue(t *testing.T) {
	input := "5"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseJson()
	checkParserErrors(t, p)
	fmt.Printf("%+v in TestNumberValue 19\n", program)
	number, ok := program.Element.Value.(*ast.NumberValue)
	if !ok {
		t.Fatalf("program.Element.Value is not ast.NumberValue. got=%T",
			program.Element.Value)
	}
	fmt.Printf("%+v\n", number)

	if number.Value != 5 {
		t.Errorf("value not %d. got=%d", 5, number.Value)
	}
	if number.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5",
			number.TokenLiteral())
	}
}

// TestStringValue 文字列単体の構文解析の単体テスト
func TestStringValue(t *testing.T) {
	input := "\"abc\""

	l := lexer.New(input)
	p := New(l)
	program := p.ParseJson()
	checkParserErrors(t, p)
	fmt.Printf("%+v in TestStringValue\n", program)
	str, ok := program.Element.Value.(*ast.StringValue)
	if !ok {
		t.Fatalf("program.Element.Value is not ast.StringValue. got=%T",
			program.Element.Value)
	}
	fmt.Printf("%+v\n", str)

	if str.Value != "abc" {
		t.Errorf("value not %s. got=%s", "abc", str.Value)
	}
	if str.TokenLiteral() != "abc" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "abc",
			str.TokenLiteral())
	}
}

// TestBooleanValue 文字列単体の構文解析の単体テスト
func TestBooleanValue(t *testing.T) {
	var tests = map[string]bool{"true": true, "false": false}

	for input, expect := range tests {
		l := lexer.New(input)
		p := New(l)
		program := p.ParseJson()
		checkParserErrors(t, p)
		fmt.Printf("%+v in TestTrueValue\n", program)
		bol, ok := program.Element.Value.(*ast.BooleanValue)
		if !ok {
			t.Fatalf("program.Element.Value is not ast.BooleanValue. got=%T",
				program.Element.Value)
		}
		fmt.Printf("%+v\n", bol)

		if bol.Value != expect {
			t.Errorf("value not %t. got=%t", expect, bol.Value)
		}
		if bol.TokenLiteral() != input {
			t.Errorf("literal.TokenLiteral not %s. got=%s", "abc",
				bol.TokenLiteral())
		}

	}

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
