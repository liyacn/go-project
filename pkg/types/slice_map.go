package types

func SliceToMap[A ~[]T, T comparable](arr A) map[T]struct{} {
	m := make(map[T]struct{}, len(arr))
	for _, v := range arr {
		m[v] = struct{}{}
	}
	return m
}

func SliceToCountMap[A ~[]T, T comparable](arr A) map[T]int {
	m := make(map[T]int, len(arr))
	for _, v := range arr {
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

func Deduplicate[A ~[]T, T comparable](arr A) A {
	return Keys(SliceToMap(arr))
}

func UnionSet[A ~[]T, T comparable](a, b A) A {
	m := SliceToMap(a)
	for _, v := range b {
		m[v] = struct{}{}
	}
	return Keys(m)
}

func IntersectionSet[A ~[]T, T comparable](a, b A) A {
	m := make(map[T]bool, len(a))
	for _, v := range a {
		m[v] = false
	}
	for _, v := range b {
		if _, ok := m[v]; ok {
			m[v] = true
		}
	}
	result := make([]T, 0, len(m))
	for k, flag := range m {
		if flag {
			result = append(result, k)
		}
	}
	return result
}

func DifferentSet[A ~[]T, T comparable](parent, child A) A {
	m := SliceToMap(parent)
	for _, v := range child {
		delete(m, v)
	}
	return Keys(m)
}

func CheckSubset[A ~[]T, T comparable](parent, child A) bool {
	m := SliceToMap(parent)
	for _, v := range child {
		if _, ok := m[v]; !ok {
			return false
		}
	}
	return true
}

func CheckEqualSet[A ~[]T, T comparable](a, b A) bool {
	m := make(map[T]bool, len(a))
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
