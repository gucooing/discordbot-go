package pgk

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"
)

var ec2b []byte

func Protoxor(users, types, cause string) []byte {
	// 读取 EC2B 客户端首包密钥
	var err error
	ec2b, err = os.ReadFile("data/ec2b.bin")
	if err != nil {
		return []byte("读取ec2b错误")
	}
	fmt.Println("使用的ec2b是：", base64.StdEncoding.EncodeToString(ec2b))
	discordbot := &Discordbot{
		Type:  types,
		Cause: cause,
		User:  users,
	}
	data, err := proto.Marshal(discordbot)
	if err != nil {
		fmt.Printf("protobuf序列化失败:", err)
	}
	//fmt.Println("protobuf序列化:", data)
	//对序列化结果的异或
	Xorec2b(data)
	//fmt.Println("xor加密结果:", data)
	rsadata := Rsaen(data)
	fmt.Println("rsa加密结果:", rsadata)
	SendMessage(rsadata)
	return rsadata
}

func Xorec2b(data []byte) {
	for i := 0; i < len(data); i++ {
		data[i] ^= ec2b[i%4096]
	}
}

type Rsa struct {
	Content string `json:"content"`
	Sign    string `json:"sign"`
}

func Rsaen(plaintext []byte) []byte {
	newplaintext := plaintext
	// 读取私钥文件
	privateKeyFile, err := ioutil.ReadFile("data/private.pem")
	if err != nil {
		fmt.Println("读取私钥文件失败：", err)
		return []byte("读取私钥文件失败")
	}
	// 解析私钥
	block, _ := pem.Decode(privateKeyFile)
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("解析私钥失败：", err)
		return []byte("解析私钥失败,密钥格式错误")
	}
	//私钥签名
	hashed := sha256.Sum256(plaintext)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), crypto.SHA256, hashed[:])
	if err != nil {
		fmt.Println("签名失败：", err)
		return []byte("签名失败")
	}

	// 读取公钥文件
	publicKeyFile, err := ioutil.ReadFile("data/public.pem")
	if err != nil {
		fmt.Println("读取公钥文件失败：", err)
		return []byte("读取公钥文件失败")
	}
	// 解析公钥
	block, _ = pem.Decode(publicKeyFile)
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		fmt.Println("解析公钥失败：", err)
		return []byte("解析公钥失败，格式错误")
	}
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	// 使用公钥加密
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, newplaintext)
	if err != nil {
		fmt.Println("RSA加密失败：", err)
		return []byte("RSA加密失败")
	}
	//整合
	newciphertext := base64.StdEncoding.EncodeToString(ciphertext)
	newsignature := base64.StdEncoding.EncodeToString(signature)
	rsadata := &Rsa{
		Content: newciphertext,
		Sign:    newsignature,
	}
	jsonrsadata, err := json.Marshal(rsadata)
	//fmt.Println("原神同款加密：", string(jsonrsadata))
	return jsonrsadata
}
