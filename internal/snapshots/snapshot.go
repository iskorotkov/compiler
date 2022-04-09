package snapshots

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Snapshot []string

//goland:noinspection GoUnusedExportedFunction
func New(value interface{}) Snapshot {
	return []string{fmt.Sprintf("%v", value)}
}

func NewSlice[T any](items []T) Snapshot {
	var s []string
	for _, item := range items {
		s = append(s, fmt.Sprintf("%v", item))
	}

	return s
}

func Load(filename string) Snapshot {
	filename = fmt.Sprintf("%s.%s", filename, "snapshot.txt")

	b, err := os.ReadFile(filename)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}

	if len(b) == 0 {
		return nil
	}

	s := strings.Split(string(b), "\n")

	return s
}

func (s Snapshot) Available() bool {
	return len(s) != 0
}

func (s Snapshot) Save(filename string) {
	filename = fmt.Sprintf("%s.%s", filename, "snapshot.txt")
	content := []byte(strings.Join(s, "\n"))

	if err := os.WriteFile(filename, content, 0744); err != nil {
		panic(err)
	}
}
