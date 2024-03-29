package context

func Scoped[T LoggerContext](ctx T, scope string) (T, func()) {
	if scope == "" {
		return ctx, func() {}
	}

	old := ctx.Logger()

	ctx.setLogger(old.Named(scope))
	return ctx, func() {
		ctx.setLogger(old)
	}
}
