package types

func SliceDivide[S ~[]E, E any](s S, parts int) [][]E {
	n := len(s)
	if n == 0 || parts <= 1 {
		return [][]E{s}
	}
	if n <= parts {
		result := make([][]E, n)
		for i := range n {
			result[i] = s[i : i+1]
		}
		return result
	}
	base := n / parts
	extra := n % parts
	left := 0
	result := make([][]E, parts)
	for i := range parts {
		size := base
		if i < extra {
			size++
		}
		right := left + size
		result[i] = s[left:right]
		left = right
	}
	return result
}

func SliceToMap[S ~[]E, E comparable](s S) map[E]struct{} {
	m := make(map[E]struct{}, len(s))
	for _, v := range s {
		m[v] = struct{}{}
	}
	return m
}

func SliceToCountMap[S ~[]E, E comparable](s S) map[E]int {
	m := make(map[E]int, len(s))
	for _, v := range s {
		m[v]++
	}
	return m
}

func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func Values[M ~map[K]V, K comparable, V any](m M) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func Deduplicate[S ~[]E, E comparable](s S) S {
	return Keys(SliceToMap(s))
}

func UnionSet[S ~[]E, E comparable](a, b S) S {
	m := SliceToMap(a)
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
	m := SliceToMap(parent)
	for _, v := range child {
		delete(m, v)
	}
	return Keys(m)
}

func CheckSubset[S ~[]E, E comparable](parent, child S) bool {
	m := SliceToMap(parent)
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
