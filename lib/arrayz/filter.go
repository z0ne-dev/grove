package arrayz

func Filter[T any](a []T, fn func(T) bool) []T {
	var r []T
	for _, v := range a {
		if fn(v) {
			r = append(r, v)
		}
	}
	return r
}
