package main

// b2f is 1 for true, 0 for false.
func b2f(v bool) float32 {
	if v {
		return 1
	}
	return 0
}

// pick is a when v, else b.
func pick[T any](v bool, a, b T) T {
	if v {
		return a
	}
	return b
}
