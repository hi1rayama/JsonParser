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

func TestObjectValue(t *testing.T) {
	tests := []struct {
		input         string
		expectedKey   []string
		expectedValue []interface{}
	}{
		{"{}", []string{""}, []interface{}{""}},
		{"{\"abg\": 100}", []string{"abg"}, []interface{}{100}},
		{"{\"abg\": \"hira\"}", []string{"abg"}, []interface{}{"hira"}},
		{"{\"abg\": true}", []string{"abg"}, []interface{}{true}},
		{"{\"abg\": false}", []string{"abg"}, []interface{}{false}},
		{input: `{"a": 100,"b":200}`,
			expectedKey:   []string{"a", "b"},
			expectedValue: []interface{}{100, 200},
		},
	
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseJson()
		checkParserErrors(t, p)
		fmt.Printf("%+v in TestObjectValue\n", program)

		// ObjectValue型かを確認
		obj, ok := program.Element.Value.(*ast.ObjectValue)
		if !ok {
			t.Fatalf("stmt not *ast.ObjectValue. got=%T", program.Element.Value)
		}
		fmt.Printf("%+v\n", *obj)

		// '{' <WS> '}' の場合
		if len(obj.Members) == 0 && tt.expectedValue[0] == "" {
			continue
		}

		// 各メンバー(要素)を確認
		for i, mem := range obj.Members {
			fmt.Printf("%s\n", mem.String())
			str, ok := mem.Key.(*ast.StringValue)
			if !ok {
				t.Fatalf("stmt not *ast.StringValue. got=%T", mem.Key)
			}
			// Keyの値が正しいかを確認
			if tt.expectedKey[i] != str.Value {
				t.Errorf("keyが違います expect=%s got=%s", tt.expectedKey, str.Value)
			}

			// Valueが正しいリテラルかを確認
			if !testLiteralValue(t, mem.Element, tt.expectedValue[i]) {
				return
			}

		}
	}
}

func TestArrayValue(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue []interface{}
	}{
		{"[ ]", []interface{}{""}},
		{"[1]", []interface{}{1}},
		{"[\"a\"]", []interface{}{"a"}},
		{"[true]", []interface{}{true}},
		{"[true]", []interface{}{true}},
		{"[1,2]", []interface{}{1,2}},
		{`["a","b"]`, []interface{}{"a","b"}},
		{`["a","b",true]`, []interface{}{"a","b",true}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseJson()
		checkParserErrors(t, p)
		fmt.Printf("%+v in TestArrayValue\n", program)

		// ArrayValue型かを確認
		arr, ok := program.Element.Value.(*ast.ArrayValue)
		if !ok {
			t.Fatalf("stmt not *ast.ArrayValue. got=%T", program.Element.Value)
		}
		fmt.Printf("%+v\n", *arr)

		// '{' <WS> '}' の場合
		if len(arr.Elements) == 0 && tt.expectedValue[0] == "" {
			continue
		}

		// 各メンバー(要素)を確認
		for i, el := range arr.Elements {
			fmt.Printf("%s\n", el.String())

			// 配列の要素が正しいリテラルかを確認
			if !testLiteralValue(t, el, tt.expectedValue[i]) {
				return
			}

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
func testLiteralValue(
	t *testing.T,
	exp ast.Element,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}


func testBooleanLiteral(t *testing.T, el ast.Element, value bool) bool {
	bo, ok := el.Value.(*ast.BooleanValue)
	if !ok {
		t.Errorf("exp not *ast.BooleanValue. got=%T", el.Value)
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

func testIntegerLiteral(t *testing.T, el ast.Element, value int64) bool {
	integ, ok := el.Value.(*ast.NumberValue)
	if !ok {
		t.Errorf("il not *ast.NumberValue. got=%T", el)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value,
			integ.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, el ast.Element, value string) bool {
	ident, ok := el.Value.(*ast.StringValue)
	if !ok {
		t.Errorf("exp not *ast.String. got=%T", el)
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
