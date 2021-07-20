/**
 * @Time : 2020/11/10 9:20 AM
 * @Author : solacowa@gmail.com
 * @File : encrypt
 * @Software: GoLand
 */

package util

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"

	"github.com/pkg/errors"
)

//RSA公钥私钥产生
func GenRsaKey() (pubkey, prvkey []byte, err error) {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		err = errors.Wrap(err, "rsa.GenerateKey")
		return
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	prvkey = pem.EncodeToMemory(block)
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		err = errors.Wrap(err, "x509.MarshalPKIXPublicKey")
		return
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	pubkey = pem.EncodeToMemory(block)
	return
}

//签名
func RsaSignWithSha256(data []byte, keyBytes []byte) (signature []byte, err error) {
	h := sha256.New()
	h.Write(data)
	hashed := h.Sum(nil)
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		err = errors.Wrap(err, "private key error")
		return
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		err = errors.Wrap(err, "ParsePKCS8PrivateKey err")
		return
	}

	signature, err = rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		err = errors.Wrap(err, "rsa.SignPKCS1v15")
		return
	}

	return
}

//验证
func RsaVerySignWithSha256(data, signData, keyBytes []byte) (bool, error) {
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return false, errors.New("public key error")
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, errors.Wrap(err, "x509.ParsePKIXPublicKey")
	}

	hashed := sha256.Sum256(data)
	err = rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), crypto.SHA256, hashed[:], signData)
	if err != nil {
		return false, errors.Wrap(err, "rsa.VerifyPKCS1v15")
	}
	return true, nil
}

// 公钥加密
func RsaEncrypt(data, keyBytes []byte) ([]byte, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "x509.ParsePKIXPublicKey")
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, pub, data)
	if err != nil {
		return nil, errors.Wrap(err, "rsa.VerifyPKCS1v15")
	}
	return ciphertext, nil
}

// 私钥解密
func RsaDecrypt(ciphertext, keyBytes []byte) ([]byte, error) {
	//获取私钥
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, errors.New("private key error")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "x509.ParsePKCS1PrivateKey")
	}
	// 解密
	data, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	if err != nil {
		return nil, errors.Wrap(err, "rsa.DecryptPKCS1v15")
	}
	return data, nil
}
