package encrypt

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

const (
	defaultSaltLen    = 256
	defaultIterations = 10000
	defaultKeyLen     = 512
)

var defaultHashFunction = sha512.New

// 用于自定义盐值长度,迭代次数,编码密钥长度以及所使用的哈希函数,
// 如果设置为 nil,则使用默认选项:EncryptOption{ 256, 10000, 512, "sha512" }
type EncryptOption struct {
	SaltLen      int
	Iterations   int
	KeyLen       int
	HashFunction func() hash.Hash
}

func EncryptString(pd string) string {
	if pd == "" {
		return ""
	}
	options := &EncryptOption{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt, codePwd := Encode(pd, options)
	passWord := fmt.Sprintf("pbkdf2-sha512$%s$%s", salt, codePwd)
	return passWord
}

func UnencryptString(rawPassword string, encryptPassword string) bool {
	if rawPassword == "" || encryptPassword == "" {
		return false
	}
	passwords := strings.SplitN(encryptPassword, "$", 3)
	opt := &EncryptOption{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	return Verify(rawPassword, passwords[1], passwords[2], opt)
}

func generateSalt(length int) []byte {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	salt := make([]byte, length)
	rand.Read(salt)
	for key, val := range salt {
		salt[key] = alphanum[val%byte(len(alphanum))]
	}
	return salt
}

// Encode函数接受两个参数:一个原始密码和一个指向Option结构体的指针,
// 如果希望使用默认选项,可以将第二个参数传递为 nil,它返回生成的盐值和用户的编码密钥,
func Encode(rawPwd string, options *EncryptOption) (string, string) {
	if options == nil {
		salt := generateSalt(defaultSaltLen)
		encodedPwd := pbkdf2.Key([]byte(rawPwd), salt, defaultIterations, defaultKeyLen, defaultHashFunction)
		return string(salt), hex.EncodeToString(encodedPwd)
	}
	salt := generateSalt(options.SaltLen)
	encodedPwd := pbkdf2.Key([]byte(rawPwd), salt, options.Iterations, options.KeyLen, options.HashFunction)
	return string(salt), hex.EncodeToString(encodedPwd)
}

// Verify函数接受四个参数:原始密码,生成的盐值,编码后的密码以及一个指向Options结构体的指针,
// 它返回一个布尔值用于确定密码是否正确,如果将最后一个参数传递为 nil,则会使用默认选项,
func Verify(rawPwd string, salt string, encodedPwd string, options *EncryptOption) bool {
	if options == nil {
		return encodedPwd == hex.EncodeToString(pbkdf2.Key([]byte(rawPwd), []byte(salt), defaultIterations, defaultKeyLen, defaultHashFunction))
	}
	return encodedPwd == hex.EncodeToString(pbkdf2.Key([]byte(rawPwd), []byte(salt), options.Iterations, options.KeyLen, options.HashFunction))
}
