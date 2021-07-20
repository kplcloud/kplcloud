/**
 * @Time : 2020/11/16 6:37 PM
 * @Author : solacowa@gmail.com
 * @File : shortid
 * @Software: GoLand
 */

package util

import (
	"strings"

	"github.com/teris-io/shortid"
)

func GenShortId(num uint8) (id string, err error) {
	short, err := shortid.New(num, shortid.DefaultABC, 1)
	if err != nil {
		return
	}
	id, err = short.Generate()
	if err != nil {
		return
	}
	alphabet := strings.ReplaceAll(id, "_", string(Krand(1, KC_RAND_KIND_ALL)))
	alphabet = strings.ReplaceAll(id, "-", string(Krand(1, KC_RAND_KIND_ALL)))
	return alphabet, err
}
