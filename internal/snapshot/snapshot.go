package snapshot

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

type Snapshot string

func New(value interface{}) Snapshot {
	return Snapshot(fmt.Sprintf("%+v", value))
}

func NewSlice[T any](items []T) Snapshot {
	var s []string
	for _, item := range items {
		s = append(s, fmt.Sprintf("%+v", item))
	}

	return Snapshot(strings.Join(s, "\n"))
}

func Load(filename string) Snapshot {
	filename = fmt.Sprintf("%s.%s", filename, "snapshot.txt")

	b, err := os.ReadFile(filename)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}

	return Snapshot(b)
}

func (s Snapshot) Equal(value interface{}) bool {
	if s == "" {
		log.Println("no snapshot available")
		return true
	}

	other := New(value)
	return s == other
}

func (s Snapshot) Save(filename string) {
	filename = fmt.Sprintf("%s.%s", filename, "snapshot.txt")
	if err := os.WriteFile(filename, []byte(s), 0744); err != nil {
		panic(err)
	}
}
