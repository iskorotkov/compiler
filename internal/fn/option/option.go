package option

type Option[T any, U error] struct {
	Some T
	None U
}

type Factory[T any, U error] struct{}

func (f Factory[T, U]) Some(val T) Option[T, U] {
	return Option[T, U]{Some: val}
}

func (f Factory[T, U]) None(err U) Option[T, U] {
	return Option[T, U]{None: err}
}
