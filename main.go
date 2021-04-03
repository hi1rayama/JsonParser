package main

import (
	"fmt"
	"io/ioutil"
	// "os"
	"jsonParser/lexer"
	"jsonParser/parser"
)

func main() {
	// ファイル読み込み
	input := useIoutilReadFile("test.json")

	l := lexer.New(input)
	p := parser.New(l)
	_ = p.ParseJson()

	if len(p.Errors()) != 0 {
		fmt.Println("Error", p.Errors())
	} else {
		fmt.Println("OK")
	}
}

func useIoutilReadFile(fileName string) string {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}