package caster

func ChanType[T any](in <-chan any) <-chan T {
	out := make(chan T, 1)
	go func() {
		defer close(out)
		for val := range in {
			out <- val.(T)
		}
	}()
	return out
}
