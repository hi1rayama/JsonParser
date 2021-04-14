package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"jsonParser/lexer"
	"jsonParser/parser"
	"flag"
)

func main() {
	// コマンドライン引数のチェック
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("usage : go run main.go filename")
		os.Exit(1)
	}
	// ファイル読み込み
	input := useIoutilReadFile(args[0])

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseJson()

	if len(p.Errors()) != 0 {
		fmt.Println("Error", p.Errors())
	} else {
		fmt.Println(program.String())
	}
}

func useIoutilReadFile(fileName string) string {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}
