/**
 * @Time : 2020/11/10 10:46 AM
 * @Author : solacowa@gmail.com
 * @File : encrypt_test
 * @Software: GoLand
 */

package util

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestGenRsaKey(t *testing.T) {
	publicKey, privateKey, err := GenRsaKey()
	if err != nil {
		t.Error(err)
	}

	bb, err := RsaEncrypt([]byte("hello"), publicKey)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(publicKey))
	fmt.Println(string(privateKey))

	cc, err := RsaDecrypt(bb, privateKey)
	if err != nil {
		t.Error(err)
	}

	println(string(cc))
}

func TestRsaDecrypt(t *testing.T) {
	data := []byte("eE0Ns5O9KvER8sTCn280WiK68jGrcT8fvG8fnebC1D0w69v7l/o86rD3EpsyH2BJYESk4TC+NrsyoxyLY8OuLR8ZvnBJ6mXkrMEy3JWKhywlV0I0qZjbT37fU3oavlNWql8x2b8p8FfmoTZKFri581Trq7hLP8wvDEk6CFcAGCs=")
	data, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		t.Error(err)
	}
	pvk := []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDJGSdrHkdJhMiBgBGBqTLil8eJe5l3gWe6WXfOUGgsjk7tUL3j
FApnv951BWIB/Z+LfAABDEauLUEQyuTss6npqhedmQZ4LpH2RrvN9WX4D4llOeF4
068S3W1exkIqY/QrERKsrCXqFD5/RtYCdtRtApm59ymxnqLfmeXCIg8W7QIDAQAB
AoGBAJ1lqbumrEc3vbPWaF1i8Cf4gj3yVtD5oRVy91mtB4xwKgiHHMjCI87Glzhi
aS6Dsz96Y4pucFfdpcKd+4XkrYSJIlcPiFnMYbrQ6KRi0HiwzvLOgEBtQgQHVvC3
XL+MxDJMYUC0vlAU7C3xysz33RKqk4dHw1mlSQXPT8MBBhohAkEA1H895Oa/N+RV
gG3GIfa5IoonrrcKhLsOd5Pt19TPZwcmmj/woVr8lqXi7/3+p1zq/otPP4NGKLM0
jjwB7r6TtQJBAPJEg3KrFld3TD+OuqMI+mBmU3h+vqv36JUqr2OSjoy8O43WdQs9
4+g/N5x8YZAqAQgsNlkcwhY3KPcadw6m6VkCQQCDnmCm9GnCY9K13siXZuubQjl8
FXIVbotyc5UhV3YzmZFGf447U1EauptLDWb7ISmJCp7Gdzgwo3dNFkwYJcD1AkEA
tf1vvSD2bIgKeCgw3a4t32Key4Jym15kkkF5dVQvz1rLZfY3AFXisaFjliL9az2S
fuAvh2uKBQQ0usNfslsCKQJAZcVtGtPTVQBOAu4bWs+x6xolZIuPJrnA2h5NUyjH
g92i1IYJsdDicqrXmzu3RVRUlPwVR7r+CkychUZCyBq94g==
-----END RSA PRIVATE KEY-----`)
	res, err := RsaDecrypt(data, pvk)
	if err != nil {
		t.Error(err)
	}
	print(string(res))
}
