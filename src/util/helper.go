/**
 * @Time : 2020/7/20 4:20 PM
 * @Author : solacowa@gmail.com
 * @File : verify
 * @Software: GoLand
 */

package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	// 中文正则匹配，合法字符为中文
	displayNamePattern = "[\u4e00-\u9fa5]*[a-z0-9A-Z-_]*$"
	// 名称的正则匹配, 合法的字符有 0-9, A-Z, a-z,-,_
	namePattern = `^[a-z0-9]([-a-z0-9])?([a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`
	// 带点的名称的正则匹配, 合法的字符有 0-9, A-Z, a-z,-,_
	namePointPattern = `^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`
)

var (
	displayNameRegexp = regexp.MustCompile(displayNamePattern)
	nameRegexp        = regexp.MustCompile(namePattern)
	namePointRegexp   = regexp.MustCompile(namePointPattern)
)

// CheckDisplayNameByString 验证中文名
func CheckDisplayNameByString(str string) bool {
	if len(str) == 0 {
		return false
	}
	return displayNameRegexp.MatchString(str)
}

// CheckNameByString 验证英文名称
func CheckNameByString(str string) bool {
	if len(str) == 0 {
		return false
	}
	return nameRegexp.MatchString(str)
}

// CheckNamePointRegexp 验证带点的英文名称
func CheckNamePointRegexp(str string) bool {
	if len(str) == 0 {
		return false
	}
	return namePointRegexp.MatchString(str)
}

// VerifyEmailFormat email verify
func VerifyEmailFormat(email string) bool {
	//pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`

	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// VerifyMobileFormat mobile verify
func VerifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

func Md5Str(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Hide(str string) (result string) {
	if len(str) == 0 {
		return
	}
	if strings.Contains(str, "@") {
		res := strings.Split(str, "@")
		star := ""
		if len(res[0]) < 3 {
			star = "***"
		} else {
			star = Substr(str, 0, 3) + "***"
		}
		result = star + "@" + res[1]
		return
	}
	reg := `^1[0-9]\d{9}$`
	rgx := regexp.MustCompile(reg)
	mobileMatch := rgx.MatchString(str)
	if mobileMatch {
		result = Substr(str, 0, 3) + "****" + Substr(str, 7, 11)
	} else {
		nameRune := []rune(str)
		lens := len(nameRune)
		if lens <= 1 {
			result = "***"
		} else if lens == 2 {
			result = string(nameRune[:1]) + "*"
		} else if lens == 3 {
			result = string(nameRune[:1]) + "*" + string(nameRune[2:3])
		} else if lens == 4 {
			result = string(nameRune[:1]) + "**" + string(nameRune[lens-1:lens])
		} else if lens > 4 {
			result = string(nameRune[:2]) + "***" + string(nameRune[lens-2:lens])
		}
	}
	return
}

// Substr 截取字符
func Substr(str string, start int, end int) string {
	rs := []rune(str)
	return string(rs[start:end])
}

func AmountToString(value int64) string {
	if value < 0 {
		return "0.00"
	}
	amountString := strconv.FormatInt(value, 10)
	switch len(amountString) {
	case 1:
		amountString = "00" + amountString
	case 2:
		amountString = "0" + amountString
	}
	return amountString[:len(amountString)-2] + "." + amountString[len(amountString)-2:]
}

var (
	headerNums    = [...]string{"139", "138", "137", "136", "135", "134", "159", "158", "157", "150", "151", "152", "188", "187", "182", "183", "184", "178", "130", "131", "132", "156", "155", "186", "185", "176", "133", "153", "189", "180", "181", "177"}
	headerNumsLen = len(headerNums)
)

func RandomPhone() string {
	rand.Seed(time.Now().UTC().UnixNano())
	header := headerNums[rand.Intn(headerNumsLen)]
	body := fmt.Sprintf("%08d", rand.Intn(99999999))
	phone := header + body
	return phone
}

func Difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func Decimal(value float64) float64 {
	return math.Trunc(value*1e2+0.5) * 1e-2
}

// RemoveDuplicateElement 去重 空间换时间
func RemoveDuplicateElement(args []string) []string {
	result := make([]string, 0, len(args))
	temp := map[string]struct{}{}
	for _, item := range args {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// FormatFileSize 字节的单位转换 保留两位小数
func FormatFileSize(fileSize int64) (size string) {
	if fileSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKi", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMi", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGi", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTi", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEi", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}
