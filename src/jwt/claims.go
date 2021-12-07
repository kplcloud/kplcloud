/**
 * @Time : 2019-06-24 11:22
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package jwt

import (
	"github.com/dgrijalva/jwt-go"
)

// ArithmeticCustomClaims 自定义声明
type ArithmeticCustomClaims struct {
	UserId  int64 `json:"userId"`
	IsAdmin bool  `json:"isAdmin"`
	jwt.StandardClaims
}

type ArithmeticTerminalClaims struct {
	UserId    int64  `json:"userId"`
	Cluster   string `json:"cluster"`
	Namespace string `json:"namespace"`
	PodName   string `json:"podName"`
	Container string `json:"container"`
	jwt.StandardClaims
}

//func Sign(email string, uid string, sessionTimeout int64) (string, error) {
//
//	expAt := time.Now().Add(time.Duration(sessionTimeout)).Unix()
//
//	fmt.Println("expAt", expAt)
//
//	// 创建声明
//	claims := ArithmeticCustomClaims{
//		UserId: uid,
//		Name:   email,
//		StandardClaims: jwt.StandardClaims{
//			ExpiresAt: expAt,
//			Issuer:    "system",
//		},
//	}
//
//	//创建token，指定加密算法为HS256
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//
//	//生成token
//	return token.SignedString([]byte(GetJwtKey()))
//}
