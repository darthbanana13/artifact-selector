package funcdecorator

type FunctionDecorator[F any] func(F) F

func DecorateFunction[F any] (f F, decorators ...FunctionDecorator[F]) F {
	decorated := f
	for _, decorator := range decorators {
		decorated = decorator(decorated)
	}
	return decorated
}
