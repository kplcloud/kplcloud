package convert

import (
	"regexp"
)

const (
	// 中文正则匹配，合法字符为中文
	displayNamePattern = "[\u4e00-\u9fa5]*[a-z0-9A-Z-_]*$"
	// 名称的正则匹配, 合法的字符有 0-9, A-Z, a-z,-,_
	nameEnPattern = `^[a-z0-9]([-a-z0-9])?([a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`

	//egressName
	egressName = `^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`
)

var (
	displayNameRegexp = regexp.MustCompile(displayNamePattern)
	nameEnRegexp      = regexp.MustCompile(nameEnPattern)
	egressNameRegexp  = regexp.MustCompile(egressName)
)

//校验egress名称是否合法
func IsEngressName(str string) bool {
	if str != "" {
		return egressNameRegexp.MatchString(str)
	}
	return false
}

// 检验是否为合法的昵称, 汉字
func IsDisplayName(b []byte) bool {
	if len(b) == 0 {
		return false
	}
	return displayNameRegexp.Match(b)
}

// 同 func IsNickname(b []byte) bool
func IsDisplayNameString(str string) bool {
	if len(str) == 0 {
		return false
	}
	return displayNameRegexp.MatchString(str)
}

// 检验是否为合法的用户名, 合法的字符有 0-9, A-Z, a-z,
func IsEnName(b []byte) bool {
	if len(b) == 0 {
		return false
	}
	return nameEnRegexp.Match(b)
}

// 同 func IsName(b []byte) bool
func IsEnNameString(str string) bool {
	if len(str) == 0 {
		return false
	}
	return nameEnRegexp.MatchString(str)
}
