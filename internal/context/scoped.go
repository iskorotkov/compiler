package context

func Scoped[T LoggerContext](ctx T, scope string) (T, func()) {
	old := ctx.Logger()

	ctx.SetLogger(old.Named(scope))
	return ctx, func() {
		ctx.SetLogger(old)
	}
}
