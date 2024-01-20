package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alidevjimmy/parsers/bnf"
	"github.com/alidevjimmy/parsers/ll1"
)

func main() {
	grammerFilePath := "Grammar_X.txt"
	fileio, err := os.Open(grammerFilePath)
	if err != nil {
		log.Fatalln(err)
	}
	
	bnf, err := bnf.NewBnf(fileio)
	if err != nil {
		log.Fatalln(err)
	}
	parser := ll1.NewLL1Parser(bnf.Grammar, "<palindrome>", "ε")

	tokensByte, err := os.ReadFile("Test_X.txt")
	if err != nil {
		log.Fatalln(err)
	}
	tokens := string(tokensByte)
	data, err := parser.Parse(tokens)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(data)
	fmt.Printf("==== %s Parsed ==== \n", tokens)
}
