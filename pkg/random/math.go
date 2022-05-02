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

type IWeight interface {
	Weight() int
}

// PickWithRate 从带权重的slice中根据指定概率随机抽取一项，未抽中返回兜底项
// items池子，out未抽中兜底项，n/d概率，weight<=0的item不会被抽中
func PickWithRate[T IWeight](items []T, out T, n, d int) T {
	if n <= 0 {
		return out
	}
	total := 0
	for _, v := range items {
		if w := v.Weight(); w > 0 {
			total += w
		}
	}
	if total == 0 {
		return out
	}
	var num int
	if d <= n {
		num = rand.IntN(total)
	} else {
		num = rand.IntN(total * d / n)
	}
	sum := 0
	for _, v := range items {
		if w := v.Weight(); w > 0 {
			sum += w
			if num < sum {
				return v
			}
		}
	}
	return out
}
