package isc

type Pair[A any, B any] struct {
	First  A
	Second B
}

type Triple[A any, B any, C any] struct {
	First  A
	Second B
	Third  C
}

func NewPair[A any, B any](a A, b B) Pair[A, B] {
	return Pair[A, B]{
		First:  a,
		Second: b,
	}
}

func NewTriple[A any, B any, C any](a A, b B, c C) Triple[A, B, C] {
	return Triple[A, B, C]{
		First:  a,
		Second: b,
		Third:  c,
	}
}
