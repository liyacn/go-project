package logger

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"log"
	"regexp"
	"strings"
	"sync"
)

var (
	once        sync.Once
	block       cipher.Block
	pattern     *regexp.Regexp
	globalLevel = levelAll
)

type Config struct {
	Level        int8
	Topic        string
	Output       string
	CipherKey    string
	CipherFields []string
}

func Initialize(cfg *Config) {
	once.Do(func() {
		if len(cfg.CipherFields) > 0 && cfg.CipherKey != "" {
			b, err := aes.NewCipher([]byte(cfg.CipherKey))
			if err != nil {
				log.Fatal(err)
			}
			block = b
			valid := regexp.MustCompile(`^\w+$`) //字段名仅可包含字母数字下划线
			for _, f := range cfg.CipherFields {
				if !valid.MatchString(f) {
					log.Fatal("invalid log.cipher.field: ", f)
				}
			}
			fields := strings.Join(cfg.CipherFields, "|")
			pattern = regexp.MustCompile(`(?i)\\*"\w*(?:` + fields + `)\w*\\*"\s*:\s*\\*"(.*?)\\*"`)
		}
		switch cfg.Output {
		case "file":
			setLogToFile(cfg.Topic)
		case "std":
			setLogToStdout()
		case "fmt":
			setLogToFormat()
		}
		if cfg.Level >= levelNone && cfg.Level <= levelAll {
			globalLevel = cfg.Level
		}
	})
}

func matchReplace(src []byte) []byte {
	if pattern == nil {
		return src
	}
	matches := pattern.FindAllSubmatchIndex(src, -1)
	if len(matches) == 0 {
		return src
	}
	buf := bytes.NewBuffer(nil)
	pre := 0
	for _, v := range matches {
		if v[2] < v[3] { // 仅1个捕获组：索引[0,1]是整个匹配，索引[2,3]是字段值，2<3说明字段值不为空
			buf.Write(src[pre:v[2]])
			buf.WriteString(encrypt(src[v[2]:v[3]]))
			pre = v[3]
		}
	}
	buf.Write(src[pre:])
	return buf.Bytes()
}

func encrypt(src []byte) string {
	bs := block.BlockSize()
	length := len(src)
	pad := bs - length%bs
	padText := bytes.Repeat([]byte{byte(pad)}, pad)
	plain := make([]byte, length, length+pad)
	copy(plain, src)
	plain = append(plain, padText...)
	length = len(plain)
	dst := make([]byte, length)
	for i := 0; i < length; i += bs {
		block.Encrypt(dst[i:i+bs], plain[i:i+bs])
	}
	return base64.StdEncoding.EncodeToString(dst)
}
