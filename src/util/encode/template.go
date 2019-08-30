/**
 * @Time : 2019-06-26 11:41
 * @Author : solacowa@gmail.com
 * @File : template
 * @Software: GoLand
 */

package encode

import (
	"bytes"
	"encoding/json"
	"html/template"
)

func EncodeTemplate(name string, tempContent string, paramContent interface{}) (string, error) {
	tmpl, err := template.New(name).Parse(tempContent)
	if err != nil {
		return "", err
	}
	var w bytes.Buffer
	err = tmpl.Execute(&w, paramContent)
	if err != nil {
		return "", err
	}
	paramContentJson, err := json.Marshal(paramContent)
	p := make([]byte, (len(tempContent)*2)+(len(string(paramContentJson))*2))
	n, err := w.Read(p)
	return string(p[:n]), nil
}

func String(n int32) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}
