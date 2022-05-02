package random

import "math/rand/v2"

const (
	upper  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lower  = "abcdefghijklmnopqrstuvwxyz"
	number = "0123456789"
	char62 = upper + lower + number
)

func Uppers(n int) string  { return gen(upper, n) }
func Lowers(n int) string  { return gen(lower, n) }
func Numbers(n int) string { return gen(number, n) }
func Strings(n int) string { return gen(char62, n) }

func gen(chars string, n int) string {
	length := len(chars)
	buf := make([]byte, n)
	for i := range buf {
		buf[i] = chars[rand.IntN(length)]
	}
	return string(buf)
}

// RangeNum 生成[s,l]区间的随机数，可包含负数
func RangeNum(s, l int) int {
	if s > l {
		s, l = l, s
	}
	n := rand.IntN(l - s + 1)
	return n + s
}

// Shuffle 打乱一个数组
func Shuffle[T any](sli []T) {
	rand.Shuffle(len(sli), func(i, j int) {
		sli[i], sli[j] = sli[j], sli[i]
	})
}

// InRate n/d 的概率返回true
func InRate(n, d int) bool {
	if n <= 0 {
		return false
	}
	if d <= n {
		return true
	}
	return rand.IntN(d) < n
}
