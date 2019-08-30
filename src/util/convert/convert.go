/**
 * @Time : 2019-07-03 16:47 
 * @Author : soupzhb@gmail.com
 * @File : convert.go
 * @Software: GoLand
 */

package convert

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"reflect"
	"strings"
)

func InfToInt(inter interface{}) (i int) {
	switch inter.(type) {
	case int:
		i = inter.(int)
		break
	}
	return
}

func Struct2Map(obj interface{}) (map[string]interface{}) {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func Struct2Json2Map(obj interface{}) (result map[string]interface{}, err error) {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return
	}
	err = json.Unmarshal(jsonBytes, &result)
	return
}

func Map2String(data []string) (result string) {
	if len(data) <= 0 {
		return
	}
	for _, v := range data {
		if strings.Contains(v, "\"") {
			result += v
		} else {
			result += "\"" + v + "\""
		}
		result += ","
	}
	result = strings.Trim(result, ",")
	return
}

func HashString(byte []byte) string {
	h := sha256.New()
	h.Write(byte)
	bs := h.Sum(nil)
	return hex.EncodeToString(bs)
}
