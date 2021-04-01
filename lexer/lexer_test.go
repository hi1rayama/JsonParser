package lexer

import (
	"jsonParser/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `{
		"A":true,
		"B":"+Hello",
		"C" :[
			1.1,
			"A",
			false,
			null
		]
	}
			  `

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LBRACE, "{"},
		{token.DOUBLEQUOT, "\""},
		{token.IDENT, "A"},
		{token.DOUBLEQUOT, "\""},
		{token.COLON, ":"},
		{token.TRUE, "true"},
		{token.COMMA,","},
		{token.DOUBLEQUOT, "\""},
		{token.IDENT,"B"},
		{token.DOUBLEQUOT, "\""},
		{token.COLON, ":"},
		{token.DOUBLEQUOT, "\""},
		{token.PLUS,"+"},
		{token.IDENT,"Hello"},
		{token.DOUBLEQUOT, "\""},
		{token.COMMA,","},
		{token.DOUBLEQUOT, "\""},
		{token.IDENT,"C"},
		{token.DOUBLEQUOT, "\""},
		{token.COLON, ":"},
		{token.LBRACKET,"["},
		{token.DIGITS,"1"},
		{token.PRIOD,"."},
		{token.DIGITS,"1"},
		{token.COMMA,","},
		{token.DOUBLEQUOT, "\""},
		{token.IDENT,"A"},
		{token.DOUBLEQUOT, "\""},
		{token.COMMA,","},
		{token.FALSE,"false"},
		{token.COMMA,","},
		{token.NULL,"null"},
		{token.RBRACKET,"]"},
		{token.RBRACE, "}"},
	}
	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected = %q got = %q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected = %q got = %q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
