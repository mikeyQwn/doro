package lib

type Selector[T any] struct {
	values     []T
	currentIdx int
}

// Initializes a selector with given values and initial value offset
// Panics if values slice is empty or if offseted index is out of bounds
func NewSelector[T any](values []T, offset int) *Selector[T] {
	if len(values) == 0 {
		panic("no values provided for the selector")
	}

	if offset < 0 || offset >= len(values) {
		panic("`offset` in selector out of bounds for array `values`")
	}

	return &Selector[T]{
		values:     values,
		currentIdx: offset,
	}
}

func (s *Selector[T]) Next() T {
	s.currentIdx += 1
	if s.currentIdx >= len(s.values) {
		s.currentIdx = 0
	}

	return s.values[s.currentIdx]
}

func (s *Selector[T]) Prev() T {
	s.currentIdx -= 1
	if s.currentIdx < 0 {
		s.currentIdx = len(s.values) - 1
	}

	return s.values[s.currentIdx]
}

func (s *Selector[T]) Curr() T {
	return s.values[s.currentIdx]
}
