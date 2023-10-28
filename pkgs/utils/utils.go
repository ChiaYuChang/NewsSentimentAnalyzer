package utils

func Unique[T comparable](v []T) []T {
	set := map[T]struct{}{}
	for _, e := range v {
		set[e] = struct{}{}
	}
	u := make([]T, 0, len(set))
	for k := range set {
		u = append(u, k)
	}
	return u
}
