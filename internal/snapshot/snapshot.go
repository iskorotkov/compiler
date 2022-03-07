package snapshot

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Snapshot []string

func New(value interface{}) Snapshot {
	return Snapshot([]string{fmt.Sprintf("%v", value)})
}

func NewSlice[T any](items []T) Snapshot {
	var s []string
	for _, item := range items {
		s = append(s, fmt.Sprintf("%v", item))
	}

	return Snapshot(s)
}

func Load(filename string) Snapshot {
	filename = fmt.Sprintf("%s.%s", filename, "snapshot.txt")

	b, err := os.ReadFile(filename)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}

	s := strings.Split(string(b), "\n")

	return Snapshot(s)
}

func (s Snapshot) Available() bool {
	return s != nil
}

func (s Snapshot) Save(filename string) {
	filename = fmt.Sprintf("%s.%s", filename, "snapshot.txt")
	content := []byte(strings.Join(s, "\n"))

	if err := os.WriteFile(filename, content, 0744); err != nil {
		panic(err)
	}
}
