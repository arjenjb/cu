package widget

func ifElse[T any](c bool, a T, b T) T {
	if c {
		return a
	} else {
		return b
	}
}
