package types

func SliceToSet[S ~[]E, E comparable](s S) map[E]struct{} {
	m := make(map[E]struct{}, len(s))
	for _, v := range s {
		m[v] = struct{}{}
	}
	return m
}

func SliceToCountSet[S ~[]E, E comparable](s S) map[E]int {
	m := make(map[E]int, len(s))
	for _, v := range s {
		m[v]++
	}
	return m
}

func Deduplicate[S ~[]E, E comparable](s S) S {
	return Keys(SliceToSet(s))
}

func UnionSet[S ~[]E, E comparable](a, b S) S {
	m := SliceToSet(a)
	for _, v := range b {
		m[v] = struct{}{}
	}
	return Keys(m)
}

func IntersectionSet[S ~[]E, E comparable](a, b S) S {
	m := make(map[E]bool, len(a))
	for _, v := range a {
		m[v] = false
	}
	for _, v := range b {
		if _, ok := m[v]; ok {
			m[v] = true
		}
	}
	result := make([]E, 0, len(m))
	for k, flag := range m {
		if flag {
			result = append(result, k)
		}
	}
	return result
}

func DifferentSet[S ~[]E, E comparable](parent, child S) S {
	m := SliceToSet(parent)
	for _, v := range child {
		delete(m, v)
	}
	return Keys(m)
}

func CheckSubset[S ~[]E, E comparable](parent, child S) bool {
	m := SliceToSet(parent)
	for _, v := range child {
		if _, ok := m[v]; !ok {
			return false
		}
	}
	return true
}

func CheckEqualSet[S ~[]E, E comparable](a, b S) bool {
	m := make(map[E]bool, len(a))
	for _, v := range a {
		m[v] = false
	}
	for _, v := range b {
		if _, ok := m[v]; ok {
			m[v] = true
		} else {
			return false
		}
	}
	for _, flag := range m {
		if !flag {
			return false
		}
	}
	return true
}
