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
	fmt.Printf("%+v in ParseProgram \n", *json)
	fmt.Printf("%+v in ParseProgram \n", p)

	// 正しく終了したかを確認する
	if p.curToken.Type != token.EOF {
		msg := fmt.Sprintf("not curToken EOF got='%s'", p.curToken.Literal)
		p.errors = append(p.errors, msg)
	}

	return json
}

func (p *Parser) ParseValue() ast.Value {
	var value ast.Value
	fmt.Printf("%+v in ParseVale before \n", p)
	switch p.curToken.Type {
	case token.DOUBLEQUOT:
		value = p.ParseString()
	case token.MINUS:
		p.nextToken()
		value = p.ParseInteger(true)
	case token.DIGITS:
		// 数値型(現在は整数)

		value = p.ParseInteger(false)
	case token.TRUE:
		// true
		value = &ast.BooleanValue{Token: p.curToken, Value: true}
	case token.FALSE:
		// false
		value = &ast.BooleanValue{Token: p.curToken, Value: false}
	case token.LBRACE:
		// オブジェクト
		value = p.ParseObject()

	case token.LBRACKET:
		// 配列
		value = p.ParseArray()
	default:
		fmt.Println("mll")
		value = nil

	}
	// トークンの更新
	p.nextToken()
	fmt.Printf("%+v in ParseVale after \n", p)

	return value

}

func (p *Parser) ParseInteger(minus bool) *ast.NumberValue {
		// 文字列型から整数型に変換
		integer, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
		if err != nil {
			msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
			p.errors = append(p.errors, msg)
			return nil
		}
		if minus {
			return &ast.NumberValue{Token: p.curToken, Value: -1 * integer}
		}else {
			return &ast.NumberValue{Token: p.curToken, Value: integer}
		}

}

func (p *Parser) ParseString() *ast.StringValue {
	// Tokenの開始位置
	//   ↓
	//   "key"

	// Tokenの終了位置
	//       ↓
	//   "key"

	if p.expectPeek(token.IDENT, true) {
		value := &ast.StringValue{Token: p.curToken, Value: p.curToken.Literal}
		if p.expectPeek(token.DOUBLEQUOT, true) {
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
func (p *Parser) ParseObject() *ast.ObjectValue {
	// Tokenの開始位置
	//  ↓
	//  {"key" : Value}

	// Tokenの終了位置
	//                ↓
	//  {"key" : Value}

	// '{' <WS> '}' の場合
	if p.expectPeek(token.RBRACE, false) {
		return &ast.ObjectValue{Token: token.Token{token.WS, ""}, Members: []ast.Member{}}
	}

	// <MEMBERS>の作成
	members := []ast.Member{}
	p.nextToken()

	// トークンを読み進めてMember解析していき、tokenが"","ではなかったら読み込みを終了させる
	for {
		// <MEMBER>の作成
		member := p.ParseMember()
		if member == nil {
			return nil
		}
		members = append(members, *member)
		fmt.Printf("%+vmember after\n",p)
		if p.curToken.Type != token.COMMA{
			break
		}
		p.nextToken()
	}
	// p.nextToken()
	fmt.Printf("%+v parseMember after\n",p)
	return &ast.ObjectValue{Token: token.Token{token.LBRACE, "{"}, Members: members}
}

// ParseArray 配列の構文解析を行う
func (p *Parser) ParseArray() *ast.ArrayValue {
	fmt.Printf("%+v in ParseArray before \n", p)
	// '[' <WS> ']' の場合
	if p.expectPeek(token.RBRACKET, false) {
		fmt.Printf("%+v in ParseArray after WS \n", p)
		return &ast.ArrayValue{Token: token.Token{token.WS, ""}, Elements: []ast.Element{}}
	}

	// <ELEMENTS>の作成
	av := &ast.ArrayValue{Token: token.Token{token.LBRACKET, "["}}
	elements := []ast.Element{}
	p.nextToken()
	for {
		// <ELEMENT>の作成
		el := ast.Element{Token:p.curToken}
		val := p.ParseValue()
		if val == nil {
			msg := fmt.Sprintf("構文エラ-")
			p.errors = append(p.errors, msg)
			return nil
		}
		el.Value = val
		elements = append(elements,el)
		if p.curToken.Type != token.COMMA{
			break
		}
		p.nextToken()
	}

	av.Elements = elements
	fmt.Printf("%+v in ParseArray after \n", p)
	return av

}

// ParseMember オブジェクトの要素(Key-Value)の構文解析を行う
func (p *Parser) ParseMember() *ast.Member {
	// Tokenの位置
	//   ↓
	// { "key": Value}

	// <MEMBER>の作成
	// <WS> <STRING> <WS> ':' <ELEMENT> ; <WS>は字句解析の時点で除外されている

	// Key(<STRING>)部分の構文解析
	fmt.Printf("%+v ParseMember Start\n",p)
	if p.curToken.Type != token.DOUBLEQUOT {
		msg := fmt.Sprintf("構文エラー 文字列は,「\"」で囲んでください2")
		p.errors = append(p.errors, msg)
		return nil
	}
	// Keyが文字列か判断
	str := p.ParseString()
	if str == nil {
		msg := fmt.Sprintf("構文エラー StringValueの生成に失敗")
		p.errors = append(p.errors, msg)
		return nil
	}
	member := &ast.Member{Token: p.curToken, Key: str}

	// コロンかどうかを確認
	if !p.expectPeek(token.COLON, true) {
		return nil
	}
	fmt.Printf("%+v156\n",p)
	p.nextToken()
	// Value<ELEMENT>部分の構文解析 <<<<<<- この箇所は再帰的な処理
	val := p.ParseValue()
	fmt.Printf("%+v val\n",val)
	member.Element.Value = val
	fmt.Printf("member %+v \n", *member)
	return member
}

// peekTokneIs 次に読み込むトークンのタイプがt(引数で与えられたトークンタイプ)と一致しているかの確認
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek 次に読み込むトークンのタイプがt(引数で与えられたトークンタイプ)と一致していたらtokenを読み込む
func (p *Parser) expectPeek(t token.TokenType, errFlag bool) bool {
	// peekTokenの型をチェックし, その型が正しい場合はnextTokenを呼んでトークンを進める
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		if errFlag {
			p.peekError(t) // 消せばTestLetStatementsのエラーが出なくなる
		}
		return false
	}
}
