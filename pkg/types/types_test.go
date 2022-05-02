package types

import "testing"

func TestAllSet(t *testing.T) {
	cases := []struct {
		a, b    []int
		u, i, d []int
		s, e    bool
	}{
		{
			a: nil,
			b: nil,
			u: nil,
			i: nil,
			d: nil,
			s: true,
			e: true,
		},
		{
			a: []int{1, 2, 3},
			b: nil,
			u: []int{1, 2, 3},
			i: nil,
			d: []int{1, 2, 3},
			s: true,
			e: false,
		},
		{
			a: nil,
			b: []int{1, 2, 3},
			u: []int{1, 2, 3},
			i: nil,
			d: nil,
			s: false,
			e: false,
		},
		{
			a: []int{1, 2, 3},
			b: []int{1, 2, 3},
			u: []int{1, 2, 3},
			i: []int{1, 2, 3},
			d: nil,
			s: true,
			e: true,
		},
		{
			a: []int{1, 2, 3},
			b: []int{2, 3, 4},
			u: []int{1, 2, 3, 4},
			i: []int{2, 3},
			d: []int{1},
			s: false,
			e: false,
		},
		{
			a: []int{1, 2, 3, 4, 5},
			b: []int{2, 3, 4},
			u: []int{1, 2, 3, 4, 5},
			i: []int{2, 3, 4},
			d: []int{1, 5},
			s: true,
			e: false,
		},
	}
	for _, c := range cases {
		if u := UnionSet(c.a, c.b); !CheckEqualSet(c.u, u) {
			t.Errorf("UnionSet(%v,%v) = %v; want %v", c.a, c.b, u, c.u)
		}
		if i := IntersectionSet(c.a, c.b); !CheckEqualSet(c.i, i) {
			t.Errorf("IntersectionSet(%v,%v) = %v; want %v", c.a, c.b, i, c.i)
		}
		if d := DifferentSet(c.a, c.b); !CheckEqualSet(c.d, d) {
			t.Errorf("DifferentSet(%v,%v) = %v; want %v", c.a, c.b, d, c.d)
		}
		if s := CheckSubset(c.a, c.b); s != c.s {
			t.Errorf("Subset(%v,%v) = %v; want %v", c.a, c.b, s, c.s)
		}
		if e := CheckEqualSet(c.a, c.b); e != c.e {
			t.Errorf("CheckEqualSet(%v,%v) = %v; want %v", c.a, c.b, e, c.e)
		}
	}
}
