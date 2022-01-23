package fn

type ValueOrErr[TVal any, TErr error] struct {
	Value TVal
	Err   TErr
}
