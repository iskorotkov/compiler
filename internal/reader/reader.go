package reader

import (
	"bufio"
	"errors"
	"io"
	"os"
	"regexp"

	"github.com/iskorotkov/compiler/internal/data/literal"
)

var (
	// wordBoundaryRegex is used for finding boundaries between two literals or other boundaries.
	// We match word boundaries so that we can extract all symbols for future analysis.
	wordBoundaryRegex = regexp.MustCompile(`\W`)
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

func (s Reader) ReadFile(filename string) (<-chan literal.Option, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return s.Read(file), nil
}

func (s Reader) Read(r io.Reader) <-chan literal.Option {
	ch := make(chan literal.Option, s.buffer)

	go func() {
		defer close(ch)

		lineNumber := literal.LineNumber(1)
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			if errors.Is(scanner.Err(), io.EOF) {
				break
			}
			if err := scanner.Err(); err != nil {
				ch <- literal.Err(err)
				return
			}

			line := scanner.Text()
			s.splitLine(line, lineNumber, ch)

			lineNumber++
		}
	}()

	return ch
}

func (s Reader) splitLine(input string, lineNumber literal.LineNumber, ch chan<- literal.Option) {
	inputLength := literal.ColNumber(len(input))
	offset := literal.ColNumber(0)
	rest := input

	for {
		boundary := wordBoundaryRegex.FindStringIndex(rest)
		if boundary == nil {
			if len(rest) > 0 {
				// Add the rest of the line.
				ch <- literal.Ok(literal.New(rest, lineNumber, offset, inputLength))
			}

			break
		}

		boundaryStart, boundaryEnd := literal.ColNumber(boundary[0]), literal.ColNumber(boundary[1])

		complexBoundary := complexOperatorRegex.FindStringIndex(rest)
		if complexBoundary != nil && boundary[0] == complexBoundary[0] {
			boundaryEnd = literal.ColNumber(complexBoundary[1])
		}

		if boundaryStart > 0 {
			// Add discovered literal.
			ch <- literal.Ok(literal.New(rest[:boundaryStart], lineNumber, offset, offset+boundaryStart))
		}

		// Add discovered boundary between two literals or other boundaries.
		ch <- literal.Ok(literal.New(rest[boundaryStart:boundaryEnd], lineNumber, offset+boundaryStart, offset+boundaryEnd))

		offset += boundaryEnd
		rest = rest[boundaryEnd:]
	}

	// Add newline.
	ch <- literal.Ok(literal.New("\n", lineNumber, inputLength, inputLength+1))
}
