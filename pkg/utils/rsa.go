package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

func GetRSAKeys() (string, string, error) {
	// 生成 RSA 密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	// 将私钥转换为 ASN.1 PKCS#1 DER 编码
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)
	//fmt.Println(hex.EncodeToString(privDER))

	// 将 DER 编码的私钥转换为 PEM 格式
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privDER,
	})
	// 将公钥提取为 *rsa.PublicKey 类型
	publicKey := &privateKey.PublicKey
	// 将公钥转换为 ASN.1 PKIX DER 编码
	pubDER, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", "", err
	}
	//fmt.Println(hex.EncodeToString(pubDER))
	// 将 DER 编码的公钥转换为 PEM 格式
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubDER,
	})
	// 返回私钥和公钥
	return string(pubPEM), string(privPEM), nil
}

// RSAEncrypt RSA加密
func RSAEncrypt(publicKeyStr, text string) (string, error) {
	block, _ := pem.Decode([]byte(publicKeyStr))
	if block == nil {
		return "", errors.New("decode public key error")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}
	res, err1 := rsa.EncryptPKCS1v15(rand.Reader, pub.(*rsa.PublicKey), []byte(text))
	// 返回 base64 编码的密文
	return base64.StdEncoding.EncodeToString(res), err1
}

// RSADecrypt RSA解密
func RSADecrypt(privateKeyStr, chipText string) (string, error) {
	// 解码 base64 编码的密文
	decodedChipText, err := base64.StdEncoding.DecodeString(chipText)
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode([]byte(privateKeyStr))
	if block == nil {
		return "", errors.New("decode private key error")
	}
	prv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	res, err1 := rsa.DecryptPKCS1v15(rand.Reader, prv, decodedChipText)
	return string(res), err1
}
