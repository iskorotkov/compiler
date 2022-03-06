package option

type Option[T any, U error] struct {
	ok  T
	err U
}

func (o Option[T, U]) Unwrap() (T, U) {
	return o.ok, o.err
}

func Ok[T any, U error](val T) Option[T, U] {
	return Option[T, U]{ok: val}
}

func Err[T any, U error](err U) Option[T, U] {
	return Option[T, U]{err: err}
}
