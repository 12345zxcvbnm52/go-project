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
	defaultSaltLen    = 16
	defaultIterations = 100
	defaultKeyLen     = 32
)

var defaultHashFunction = sha512.New()

// 用于自定义盐值长度,迭代次数,编码密钥长度以及所使用的哈希函数,
// 如果设置为 nil,则使用默认选项:EncryptOption{ 256, 10000, 512, "sha512" }
type EncryptOption struct {
	SaltLen      int
	Iterations   int
	KeyLen       int
	HashFunction func() hash.Hash
}

func NewDefaultConfig() *EncryptOption {
	return &EncryptOption{
		SaltLen:      defaultSaltLen,
		Iterations:   defaultIterations,
		KeyLen:       defaultKeyLen,
		HashFunction: sha512.New,
	}
}

// 默认使用sha512加密
func EncryptString(pd string) string {
	return EncryptStringWithHashFunc(pd, defaultHashAlgorithm)
}

// 若不指出或hashfunc错误则默认使用sha512
func EncryptStringWithHashFunc(pd string, hashfunc string) string {
	if pd == "" {
		return ""
	}
	hash, hashstr := NewHashFuncForPBKDF2WithDefault(hashfunc)
	options := &EncryptOption{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: hash}
	salt, codePwd := encode(pd, options)
	passWord := fmt.Sprintf("pbkdf2-%s$%s$%s", hashstr, salt, codePwd)
	return passWord
}

func UnencryptString(rawPassword string, encryptPassword string) bool {
	if rawPassword == "" || encryptPassword == "" {
		return false
	}
	passwords := strings.SplitN(encryptPassword, "$", 3)
	if len(passwords) != 3 {
		return false
	}
	hashstr := strings.SplitN(passwords[0], "-", 2)
	if len(hashstr) != 2 {
		return false
	}
	hash, err := NewHashFuncForPBKDF2(hashstr[1])
	if err != nil {
		return false
	}
	opt := &EncryptOption{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: hash}
	return verify(rawPassword, passwords[1], passwords[2], opt)
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

// encode函数接受两个参数:一个原始密码和一个指向Option结构体的指针,
// 如果希望使用默认选项,可以将第二个参数传递为 nil,它返回生成的盐值和用户的编码密钥,
func encode(rawPwd string, options *EncryptOption) (string, string) {
	if options == nil {
		options = NewDefaultConfig()
	}
	salt := generateSalt(options.SaltLen)
	encodedPwd := pbkdf2.Key([]byte(rawPwd), salt, options.Iterations, options.KeyLen, options.HashFunction)
	return string(salt), hex.EncodeToString(encodedPwd)
}

// Verify函数接受四个参数:原始密码,生成的盐值,编码后的密码以及一个指向Options结构体的指针,
// 它返回一个布尔值用于确定密码是否正确,如果将最后一个参数传递为 nil,则会使用默认选项,
func verify(rawPwd string, salt string, encodedPwd string, options *EncryptOption) bool {
	if options == nil {
		options = NewDefaultConfig()
	}
	buf := pbkdf2.Key([]byte(rawPwd), []byte(salt), options.Iterations, options.KeyLen, options.HashFunction)
	return encodedPwd == hex.EncodeToString(buf)
}
