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

func SlicePluck[S ~[]E, E, V any, K comparable](s S, f func(E) (K, V)) map[K]V {
	m := make(map[K]V, len(s))
	for _, e := range s {
		k, v := f(e)
		m[k] = v
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

func KeyValues[M ~map[K]V, K comparable, V any](m M) ([]K, []V) {
	keys := make([]K, 0, len(m))
	values := make([]V, 0, len(m))
	for k, v := range m {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}
