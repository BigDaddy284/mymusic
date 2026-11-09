package main

// lerp is the value t of the way from a to b.
func lerp(a, b, t float32) float32 { return a + (b-a)*t }

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
