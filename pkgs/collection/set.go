package collection

type Set[T comparable] map[T]bool

func NewSet[T comparable](x ...T) Set[T] {
	s := make(Set[T])
	for i := 0; i < len(x); i++ {
		s.Add(x[i])
	}
	return s
}

func (s Set[T]) Add(x T) {
	s[x] = true
}

func (s Set[T]) Has(x T) bool {
	isVadid, ok := s[x]
	return isVadid && ok
}

func (s Set[T]) Del(x T) {
	s[x] = false
}

func (s Set[T]) Key() []T {
	key := []T{}
	for k, v := range s {
		if v {
			key = append(key, k)
		}
	}
	return key
}

func (s1 Set[T]) Merge(s2 Set[T]) Set[T] {
	s := NewSet[T]()
	for x := range s1 {
		s.Add(x)
	}

	for x := range s2 {
		s.Add(x)
	}
	return s
}
