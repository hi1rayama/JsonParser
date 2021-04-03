package ast

import (
	"bytes"
	"jsonParser/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Value interface {
	Node
	valueNode()
}

type Json struct {
	Element Element
}

type Element struct {
	Token token.Token
	Value Value
}

func (e *Element) TokenLiteral() string { return e.Token.Literal }
func (e *Element) String() string {

	if e.Value != nil {
		return e.Value.String()
	}
	return ""
}

// Object JSONのオブジェクトを表す構造体
// <OBJECT> ::= '{' <WS> '}' | '{' <MUMBERS> '}'
type ObjectValue struct {
	Token   token.Token // '{'トークン
	Members []Member    // オブジェクトの要素
}

func (ov *ObjectValue) valueNode()           {}
func (ov *ObjectValue) TokenLiteral() string { return ov.Token.Literal }
func (ov *ObjectValue) String() string {
	var out bytes.Buffer

	out.WriteString("{")

	for _, s := range ov.Members {
		out.WriteString(s.String())
	}
	out.WriteString("}")

	return out.String()

}

// MemberLiteral オブジェクトの 「key : value」 の を表現する構造体
// <MEMBER> ::= <WS> <STRING> <WS> ':' <ELEMENT>
type Member struct {
	Token   token.Token // 識別子(key)トークン
	Element Element
}

func (m *Member) valueNode()           {}
func (m *Member) TokenLiteral() string { return m.Token.Literal }
func (m *Member) String() string {
	var out bytes.Buffer
	out.WriteString(m.Token.Literal)
	out.WriteString(":")
	out.WriteString(m.Element.String())

	return out.String()
}

// NumberValue 数字を表す構造体(最初は整数型のみ)
// <NUMBER> ::= <INTEGER> <FRACTION> <EXPONENT>
type NumberValue struct {
	Token token.Token
	Value int64
}

func (nv *NumberValue) valueNode()           {}
func (nv *NumberValue) TokenLiteral() string { return nv.Token.Literal }
func (nv *NumberValue) String() string {
	return nv.Token.Literal
}

// StringValue 文字列リテラルを表す構造体
type StringValue struct {
	Token token.Token
	Value string
}

func (sv *StringValue) valueNode()           {}
func (sv *StringValue) TokenLiteral() string { return sv.Token.Literal }
func (sv *StringValue) String() string {
	return sv.Token.Literal
}

type BooleanValue struct {
	Token token.Token
	Value bool
}

func (b *BooleanValue) valueNode()           {}
func (b *BooleanValue) TokenLiteral() string { return b.Token.Literal }
func (b *BooleanValue) String() string       { return b.Token.Literal }
