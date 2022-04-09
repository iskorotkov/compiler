package reader

import (
	"bufio"
	"errors"
	"io"
	"regexp"

	"github.com/iskorotkov/compiler/internal/data/literal"
	"github.com/iskorotkov/compiler/internal/fn/options"
)

var (
	// wordBoundaryRegex is used for finding boundaries between two literals or other boundaries.
	// We match word boundaries so that we can extract all symbols for future analysis.
	wordBoundaryRegex = regexp.MustCompile(`\W`)

	// doubleConstantRegex is used for finding double constants.
	doubleConstantRegex = regexp.MustCompile(`^\d+\.\d+`)

	// complexOperatorRegex matches complex operators that consist of 2 characters.
	complexOperatorRegex = regexp.MustCompile(`[<>][<>=]|:=`)
)

type Reader struct {
	buffer int
}

func New(buffer int) *Reader {
	return &Reader{
		buffer: buffer,
	}
}

func (s Reader) Read(ctx interface{}, r io.Reader) <-chan options.Option[literal.Literal] {
	ch := make(chan options.Option[literal.Literal], s.buffer)

	go func() {
		defer close(ch)

		lineNumber := literal.LineNumber(1)
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			if errors.Is(scanner.Err(), io.EOF) {
				break
			}
			if err := scanner.Err(); err != nil {
				ch <- options.Err[literal.Literal](err)
				return
			}

			line := scanner.Text()
			s.splitLine(ctx, line, lineNumber, ch)

			lineNumber++
		}
	}()

	return ch
}

//goland:noinspection GoUnusedParameter
func (s Reader) splitLine(ctx interface{}, input string, lineNumber literal.LineNumber, ch chan<- options.Option[literal.Literal]) {
	inputLength := literal.ColNumber(len(input))
	offset := literal.ColNumber(0)
	rest := input

	for {
		boundary := wordBoundaryRegex.FindStringIndex(rest)
		if boundary == nil {
			if len(rest) > 0 {
				// Add the rest of the line.
				ch <- options.Ok(literal.New(rest, lineNumber, offset+1, inputLength+1))
			}

			break
		}

		boundaryStart, boundaryEnd := literal.ColNumber(boundary[0]), literal.ColNumber(boundary[1])

		// Extend selection for complex operators.
		complexBoundary := complexOperatorRegex.FindStringIndex(rest)
		if complexBoundary != nil && boundary[0] == complexBoundary[0] {
			boundaryEnd = literal.ColNumber(complexBoundary[1])
		} else {
			// Expand selection for double constants.
			doubleConstantBoundary := doubleConstantRegex.FindStringIndex(rest)
			if doubleConstantBoundary != nil {
				boundaryStart = literal.ColNumber(doubleConstantBoundary[1])
				boundaryEnd = boundaryStart + 1
			}
		}

		if boundaryStart > 0 {
			// Add discovered literal.
			ch <- options.Ok(literal.New(rest[:boundaryStart], lineNumber, offset+1, offset+boundaryStart+1))
		}

		// Add discovered boundary between two literals or other boundaries.
		ch <- options.Ok(literal.New(rest[boundaryStart:boundaryEnd], lineNumber, offset+boundaryStart+1, offset+boundaryEnd+1))

		offset += boundaryEnd
		rest = rest[boundaryEnd:]
	}

	// Add newline.
	ch <- options.Ok(literal.New("\n", lineNumber, inputLength+1, inputLength+2))
}
