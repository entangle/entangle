package main

import (
	"./errors"
	"./parser"
	"./source"
	"fmt"
	"os"
)

func main() {
	f, err := os.Open("test.entangle")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	src, err := source.FromReader(f, "test.entangle")

	if err != nil {
		panic(err)
	}

	interfaceDecl, err := parser.Parse(src)

	if parseErr, ok := err.(errors.ParseError); ok {
		parser.PrintError(parseErr)
	} else {
		fmt.Println(err)
	}

	fmt.Println(interfaceDecl)
}
