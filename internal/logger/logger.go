package logger

import (
	"fmt"
	"log"
	"os"
)

func New(tag string) *log.Logger {
	return log.New(os.Stdout, fmt.Sprintf("[%s] ", tag), log.LstdFlags|log.Lmsgprefix|log.Lshortfile)
}
