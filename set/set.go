package set

type StrSet map[string]int

func New() StrSet {
	return make(StrSet)
}

func (s StrSet) Set(key string) {
	s[key] = 1
}

func (s StrSet) Has(key string) bool {
	_, ok := s[key]
	return ok
}

func (s StrSet) Rem(key string) {
	delete(s, key)
}

func (s StrSet) Keys() []string {
	var key []string
	for k := range s {
		key = append(key, k)
	}
	return key
}

func (s StrSet) Clr() {
	for k := range s {
		delete(s, k)
	}
}
