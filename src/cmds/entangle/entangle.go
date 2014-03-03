package main

import (
	"entangle/errors"
	"entangle/parser"
	"entangle/source"
	"entangle/generators"
	"entangle/generators/golang"
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

	//
	generators := make([]generators.Generator, 0)

	gen, err := golang.NewGenerator()
	if err != nil {
		panic(err)
	}

	generators = append(generators, gen)

	//
	outputPath := "_testing/src/users"

	// Make sure the output directory exists.
	var outputPathStat os.FileInfo
	var statErr error

	if outputPathStat, statErr = os.Stat(outputPath); statErr != nil {
		if !os.IsNotExist(statErr) {
			return
		}

		if statErr = os.MkdirAll(outputPath, 0777); statErr != nil {
			return
		}
	} else if !outputPathStat.Mode().IsDir() {
		err = fmt.Errorf("output path '%s' is not a directory", outputPath)
		return
	}

	for _, gen := range generators {
		err = gen.Generate(interfaceDecl, outputPath)
		if err != nil {
			panic(err)
		}
	}
}
