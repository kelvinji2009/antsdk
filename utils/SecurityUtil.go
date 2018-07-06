package utils

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"sort"
)

var hash crypto.Hash

func SetHash(t string) {
	switch t {
	case "RSA":
		hash = crypto.SHA1
	case "RSA2":
		hash = crypto.SHA256
	}
}

func GetSignMap(requestHolder *RequestParametersHolder) map[string]string {
	singleMap := make(map[string]string)

	if requestHolder.ApplicationParams.Length > 0 {
		for k, v := range requestHolder.ApplicationParams.GetMap() {
			singleMap[k] = v
		}
	}

	if requestHolder.ProtocalMustParams.Length > 0 {
		for k, v := range requestHolder.ProtocalMustParams.GetMap() {
			singleMap[k] = v
		}
	}

	if requestHolder.ProtocalOptParams.Length > 0 {
		for k, v := range requestHolder.ProtocalOptParams.GetMap() {
			singleMap[k] = v
		}
	}

	return singleMap
}

func GetSignStr(m map[string]string) string {
	// 对 key 进行升序排序
	sortedKeys := make([]string, 0)
	for k, _ := range m {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	// 对 key=value 的键值对用 & 连接起来，略过空值
	sbSignStr := NewStringBuilder()
	for i, k := range sortedKeys {
		if m[k] != "" {
			sbSignStr.Append(k)
			sbSignStr.Append("=")
			sbSignStr.Append(m[k])
			if i != (len(sortedKeys) - 1) {
				sbSignStr.Append("&")
			}
		}
	}
	return sbSignStr.ToString()
}

func Sign(mReq map[string]string, privateKey []byte) (string, error) {

	// 获取待签名参数
	signStr := GetSignStr(mReq)

	block, _ := pem.Decode(privateKey)
	if block == nil {
		return "", errors.New("Sign private key decode error")
	}

	prk8, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	pKey := prk8.(*rsa.PrivateKey)

	return RSASign(signStr, pKey)
}

func RSASign(origData string, privateKey *rsa.PrivateKey) (string, error) {
	h := hash.New()
	h.Write([]byte(origData))
	digest := h.Sum(nil)

	s, err := rsa.SignPKCS1v15(nil, privateKey, hash, []byte(digest))
	if err != nil {
		return "", err
	}
	data := base64.StdEncoding.EncodeToString(s)
	return string(data), nil
}

// 同步返回验签 https://docs.open.alipay.com/200/106120
func SyncVerifySign(sign string, body, alipayPublicKey []byte) (bool, error) {
	return RSAVerify(body, []byte(sign), alipayPublicKey)
}

// 异步返回验签
func AsyncVerifySign(body, alipayPublicKey []byte) (bool, error) {
	data, err := url.ParseQuery(string(body))
	if err != nil {
		return false, err
	}

	var m map[string]string
	m = make(map[string]string, 0)

	for k, v := range data {
		if k == "sign" || k == "sign_type" { //不要'sign'和'sign_type'
			continue
		}
		m[k] = v[0]
	}

	sign := data["sign"][0]

	//获取要进行计算哈希的sign string
	signStr := GetSignStr(m)

	return RSAVerify([]byte(signStr), []byte(sign), alipayPublicKey)
}

func RSAVerify(src, sign, alipayPublicKey []byte) (bool, error) {
	// 加载RSA的公钥
	block, _ := pem.Decode(alipayPublicKey)
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, err
	}
	rsaPub, _ := pub.(*rsa.PublicKey)

	// 计算待签名字串的哈希
	t := hash.New()
	io.WriteString(t, string(src))
	digest := t.Sum(nil)

	// base64 decode,必须步骤，支付宝对返回的签名做过base64 encode必须要反过来decode才能通过验证
	data, _ := base64.StdEncoding.DecodeString(string(sign))

	// 调用rsa包的VerifyPKCS1v15验证签名有效性
	err = rsa.VerifyPKCS1v15(rsaPub, hash, digest, data)
	if err != nil {
		return false, err
	}

	return true, nil
}

func ReadPemFile(path string) []byte {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return fd
}
