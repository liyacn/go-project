package ades

import (
	"project/pkg/random"
	"testing"
)

func FuzzNewAesCipher(f *testing.F) {
	key := []byte("a0EYDkdjIvn4N92U")
	iv := []byte("Te76o8pw0h9wRW1e")
	c, err := NewAesCTR(key, iv)
	if err != nil {
		f.Fatal(err)
	}
	text := []string{random.Strings(10), "13800138000"}
	for _, v := range text {
		f.Add(v)
	}
	f.Fuzz(func(t *testing.T, src string) {
		encrypted := c.Encrypt([]byte(src))
		decrypted := c.Decrypt(encrypted)
		if string(decrypted) == src {
			t.Logf("src:%s, encrypted:%s", src, encrypted)
		} else {
			t.Errorf("src:%s, encrypted:%s, decrypted:%s\n", src, encrypted, decrypted)
		}
	})
}
