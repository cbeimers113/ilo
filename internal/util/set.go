package util

type Set map[any]struct{}

func NewSet(items ...any) Set {
	s := make(map[any]struct{}, 0)

	for _, item := range items {
		s[item] = struct{}{}
	}

	return s
}

func (s Set) Insert(item any) {
	s[item] = struct{}{}
}

func (s Set) Contains(item any) bool {
	_, ok := s[item]
	return ok
}
