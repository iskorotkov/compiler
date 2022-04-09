package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/iskorotkov/compiler/internal/contexts"
	"github.com/iskorotkov/compiler/internal/modules/reader"
	"github.com/iskorotkov/compiler/internal/modules/scanner"
	"github.com/iskorotkov/compiler/internal/modules/syntax_analyzer"
	"github.com/iskorotkov/compiler/internal/modules/syntax_neutralizer"
)

func main() {
	ctx := contexts.NewEnvContext(context.Background())

	if len(os.Args) > 1 {
		file, err := os.OpenFile(os.Args[1], os.O_RDONLY, 0666)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		compile(ctx, file)

		return
	}

	compile(ctx, os.Stdin)
}

func compile(ctx contexts.FullContext, r io.Reader) {
	buffer := 0

	rd := reader.New(buffer)
	literals := rd.Read(ctx, r)

	sc := scanner.New(buffer)
	tokens := sc.Scan(ctx, literals)

	neutralizer := syntax_neutralizer.New(1)

	sa := syntax_analyzer.New(buffer)
	syntaxConstructions := sa.Analyze(struct {
		contexts.LoggerContext
		contexts.NeutralizerContext
	}{ctx, contexts.NewNeutralizerContext(neutralizer)}, tokens)

	_, err := (<-syntaxConstructions).Unwrap()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("compiled successfully")
}
