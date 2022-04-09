package main

import (
	"fmt"
	"io"
	"os"

	"github.com/iskorotkov/compiler/internal/logger"
	"github.com/iskorotkov/compiler/internal/modules/reader"
	"github.com/iskorotkov/compiler/internal/modules/scanner"
	"github.com/iskorotkov/compiler/internal/modules/syntax_analyzer"
)

//goland:noinspection GoUnusedGlobalVariable
var log = logger.New().Named("main")

func main() {
	if len(os.Args) > 1 {
		file, err := os.OpenFile(os.Args[1], os.O_RDONLY, 0666)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		compile(file)

		return
	}

	compile(os.Stdin)
}

func compile(r io.Reader) {
	buffer := 0

	rd := reader.New(buffer)
	literals := rd.Read(r)

	sc := scanner.New(buffer)
	tokens := sc.Scan(literals)

	sa := syntax_analyzer.New(buffer, 1)
	syntaxConstructions := sa.Analyze(tokens)

	_, err := (<-syntaxConstructions).Unwrap()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("compiled successfully")
}
