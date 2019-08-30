package email

import (
	"crypto/sha1"
	"fmt"
	"strings"
)

type Security interface {
	Sha1Digest(src string) string
	Sha1DigestSalt(src string, salt []byte) string
}

type security struct {
}

func NewSecurity() Security {
	return &security{}
}

func (c *security) Sha1Digest(src string) string {
	md := sha1.New()
	//md.Write([]byte(src))
	return c.byte2HexStr(md.Sum([]byte(src)))
}

func (c *security) Sha1DigestSalt(src string, salt []byte) string {
	md := sha1.New()
	md.Write(salt)
	md.Write([]byte(src))

	return c.byte2HexStr(md.Sum(nil))
}

func (c *security) byte2HexStr(b []byte) string {
	var sb strings.Builder
	for _, bt := range b {
		s := fmt.Sprintf("%x", bt)
		if len(s) == 1 {
			sb.WriteByte('0')
		}
		sb.WriteString(strings.ToUpper(s))
	}
	return sb.String()
}
