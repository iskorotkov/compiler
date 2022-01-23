package fn

type Option[T any, U error] struct {
	Some T
	None U
}

type OptionFactory[T any, U error] struct{}

func (f OptionFactory[T, U]) Some(val T) Option[T, U] {
	return Option[T, U]{Some: val}
}

func (f OptionFactory[T, U]) None(err U) Option[T, U] {
	return Option[T, U]{None: err}
}
