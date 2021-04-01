package lexer

import (
	"jsonParser/token"
	"fmt"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte // 現在検知中の文字
}

// readChar 文字を読み込む関数
func (l *Lexer) readChar() {

	// 読み込み位置が終端に達したかをチェックする
	if l.readPosition >= len(l.input) {
		l.ch = 0 //EOF
	} else {
		l.ch = l.input[l.readPosition]
	}

	// 現在の位置と次に読み込む位置を更新
	l.position = l.readPosition
	l.readPosition++
}

// NextToken 指定された位置のtokenの認識を行う
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	// スペースをスキップする
	l.skipWhitespace()

	// トークンの認識を行う
	switch l.ch {
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case ':':
		tok = newToken(token.COLON,l.ch)
	case '"':
		tok = newToken(token.DOUBLEQUOT,l.ch)
	case ',':
		tok = newToken(token.COMMA,l.ch)
	case '.':
		tok = newToken(token.PRIOD,l.ch)
	case '+':
		tok = newToken(token.PLUS,l.ch)
	case '-':
		tok = newToken(token.MINUS,l.ch)
	case '[':
		tok = newToken(token.LBRACKET,l.ch)
	case ']':
		tok = newToken(token.RBRACKET,l.ch)
	
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		// キーワード(予約語)や変数を認識する
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			fmt.Printf("%+v\n",tok)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.DIGITS
			tok.Literal = l.readNumber()
			fmt.Printf("%+v\n",tok)
			return tok
		}
	}
	l.readChar()

	return tok
}

// peekChar
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

// skipWhitespace スペースをスキップする
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// readIdetifier 連続する文字(文字列)を特定し, 文字列を返す
func (l *Lexer) readIdentifier() string {
	// 文字列の開始位置の保存
	position := l.position

	// 連続する文字列の判別
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]

}

// readNumber 連続する数字を特定し,実数か整数かの判断を行い連続した数字とトークンタイプとして返す
func (l *Lexer) readNumber() (string) {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// newToken 新しいtokenを生成する
func newToken(tokenType token.TokenType, ch byte) token.Token {
	newTK := token.Token{Type: tokenType, Literal: string(ch)}
	fmt.Printf("%+v\n",newTK)
	return newTK
}

// isLetter 文字かどうかを判断する
func isLetter(ch byte) bool {
	return 'a' <= ch && 'z' >= ch || 'A' <= ch && 'Z' >= ch || ch == '_'
}

// isDigit 数字かどうかを判断する
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// New 字句解析器の生成(インスタンス)
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}
