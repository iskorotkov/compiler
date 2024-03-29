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
	"github.com/iskorotkov/compiler/internal/module/typechecker"
	"github.com/iskorotkov/compiler/internal/module/wasm"
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

	syn := syntax_analyzer.New(buffer)
	programs := syn.Analyze(ctx, tokens)

	checker := typechecker.NewTypeChecker(buffer)
	results := checker.Check(ctx, programs)

	generator := wasm.NewGenerator()
	if m, ok := <-generator.Generate(ctx, results); ok {
		if len(os.Args) > 2 {
			file, err := os.Create(os.Args[2])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if _, err = file.WriteString(m.String()); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if _, err := generator.WAT2WASM(file.Name()); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			fmt.Println(m.String())
		}
	}

	writeReport(ctx)
}

func writeReport(ctx context.FullContext) {
	if len(ctx.Errors()) == 0 {
		fmt.Println("compiled successfully")
		return
	}

	fmt.Printf("finished with %d error(s)\n", len(ctx.Errors()))
	for _, err := range ctx.Errors() {
		fmt.Println(err)
	}
}
