package parser

import (
	"fmt"
	"jsonParser/ast"
	"jsonParser/lexer"
	"jsonParser/token"
	"strconv"
)

// Parser 構文解析器
type Parser struct {
	l *lexer.Lexer // 字句解析器

	errors []string

	curToken  token.Token // 現在のトークン
	peekToken token.Token // 次のトークン
}

// New Parserのコンストラクタ
func New(l *lexer.Lexer) *Parser {

	// 構文解析器のインスタンス生成
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.nextToken()
	p.nextToken()
	fmt.Printf("%+v in New \n", p)
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// nextToken curTokenとpeekTokenを進めるヘルパーメソッド
func (p *Parser) nextToken() {

	// curTokenを進める
	p.curToken = p.peekToken
	// peekTokenを進める
	p.peekToken = p.l.NextToken()

}

func (p *Parser) ParseJson() *ast.Json {
	// ASTのルートノードを生成
	json := &ast.Json{}
	json.Element.Value = p.ParseValue()
	fmt.Printf("%+v in ParseProgram \n", json)

	// 正しく終了したかを確認する
	if p.curToken.Type != token.EOF {
		msg := fmt.Sprintf("not curToken EOF got='%s'", p.curToken.Literal)
		p.errors = append(p.errors, msg)
	}

	return json
}

func (p *Parser) ParseValue() ast.Value {
	var value ast.Value
	switch p.curToken.Type {
	case token.DOUBLEQUOT:
		value = p.ParseString()
		p.nextToken()
	case token.DIGITS:
		// 数値型(現在は整数)

		// 文字列型から整数型に変換
		integer, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
		if err != nil {
			msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
			p.errors = append(p.errors, msg)
			return nil
		}
		value = &ast.NumberValue{Token: p.curToken, Value: integer}
	case token.TRUE:
		// true
		value = &ast.BooleanValue{Token: p.curToken, Value: true}
	case token.FALSE:
		// false
		value = &ast.BooleanValue{Token: p.curToken, Value: false}
		// case token.LBRACKET:
		// 	// 配列
		// case token.LBRACE:
		// 	// オブジェクト
		// case token.NULL:
		// 	// null

	}
	// トークンの更新
	p.nextToken()

	return value

}

func (p *Parser) ParseString() *ast.StringValue {

	if p.expectPeek(token.IDENT) {
		value := &ast.StringValue{Token: p.curToken, Value: p.curToken.Literal}
		if p.expectPeek(token.DOUBLEQUOT) {
			return value
		} else {
			return nil
		}
	}
	msg := fmt.Sprintf("構文エラー 文字列は,「\"」で囲んでください")
	p.errors = append(p.errors, msg)
	return nil
}

// ParseObject オブジェクトを構文解析を行うメソッド(一番重要!!!!!!!!!!!!!)
// func (p *Parser) ParseObject() ast.ObjectValue {

// }

// peekTokneIs 次に読み込むトークンのタイプがt(引数で与えられたトークンタイプ)と一致しているかの確認
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek 次に読み込むトークンのタイプがt(引数で与えられたトークンタイプ)と一致していたらtokenを読み込む
func (p *Parser) expectPeek(t token.TokenType) bool {
	// peekTokenの型をチェックし, その型が正しい場合はnextTokenを呼んでトークンを進める
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t) // 消せばTestLetStatementsのエラーが出なくなる
		return false
	}
}
