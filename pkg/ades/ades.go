package ades

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type Cipher interface {
	Encrypt(src []byte) string
	Decrypt(s string) []byte
}

type (
	param struct {
		block cipher.Block
		iv    []byte
	}
	ctr param
)

func NewAesCTR(key, iv []byte) (Cipher, error) {
	block, ivb, err := newParam(aes.NewCipher, key, iv)
	return ctr{block, ivb}, err
}

func newParam(f func([]byte) (cipher.Block, error), key, iv []byte) (cipher.Block, []byte, error) {
	block, err := f(key)
	if err != nil {
		return nil, nil, err
	}
	ivb := make([]byte, block.BlockSize())
	copy(ivb, iv)
	return block, ivb, nil
}

func (c ctr) Encrypt(src []byte) string {
	dst := make([]byte, len(src))
	stream := cipher.NewCTR(c.block, c.iv)
	stream.XORKeyStream(dst, src)
	return base64.StdEncoding.EncodeToString(dst)
}
func (c ctr) Decrypt(s string) []byte {
	src, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil
	}
	dst := make([]byte, len(src))
	stream := cipher.NewCTR(c.block, c.iv)
	stream.XORKeyStream(dst, src)
	return dst
}
