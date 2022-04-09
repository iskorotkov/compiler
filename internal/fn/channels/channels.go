package channels

func FromSlice[T any](values []T) <-chan T {
	c := make(chan T, len(values))
	defer close(c)

	for _, value := range values {
		c <- value
	}

	return c
}

func ToSlice[T any](c <-chan T) []T {
	var s []T
	for value := range c {
		s = append(s, value)
	}

	return s
}
