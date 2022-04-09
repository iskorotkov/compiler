package options

import (
	"fmt"
	"reflect"
)

type Option[T any] struct {
	ok  T
	err error
}

func (o Option[T]) Unwrap() (T, error) {
	return o.ok, o.err
}

func (o Option[T]) String() string {
	if reflect.ValueOf(&o).Elem().FieldByName("err").IsNil() {
		return fmt.Sprintf("ok: %v", o.ok)
	}

	return fmt.Sprintf("err: %v", o.err)
}

func Ok[T any](val T) Option[T] {
	return Option[T]{ok: val}
}

func Err[T any](err error) Option[T] {
	return Option[T]{err: err}
}
