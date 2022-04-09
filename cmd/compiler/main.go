package main

import (
	stdcontext "context"
	"fmt"
	"io"
	"os"

	"github.com/iskorotkov/compiler/internal/context"
	"github.com/iskorotkov/compiler/internal/module/reader"
	"github.com/iskorotkov/compiler/internal/module/scanner"
	"github.com/iskorotkov/compiler/internal/module/syntax_analyzer"
	"github.com/iskorotkov/compiler/internal/module/syntax_neutralizer"
)

func main() {
	ctx := context.NewEnvContext(stdcontext.Background())

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

func compile(ctx context.FullContext, r io.Reader) {
	buffer := 0

	rd := reader.New(buffer)
	literals := rd.Read(ctx, r)

	sc := scanner.New(buffer)
	tokens := sc.Scan(ctx, literals)

	neutralizer := syntax_neutralizer.New(1)

	sa := syntax_analyzer.New(buffer)
	syntaxConstructions := sa.Analyze(struct {
		context.LoggerContext
		context.NeutralizerContext
	}{ctx, context.NewNeutralizerContext(neutralizer)}, tokens)

	_, err := (<-syntaxConstructions).Unwrap()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("compiled successfully")
}
