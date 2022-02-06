package option

type Factory[T any, U error] struct{}

func (f Factory[T, U]) Ok(val T) Option[T, U] {
	return Ok[T, U](val)
}

func (f Factory[T, U]) Err(err U) Option[T, U] {
	return Err[T, U](err)
}
