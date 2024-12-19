package pkg

import (
	"crypto/sha512"
	"fmt"

	passwd "github.com/anaskhan96/go-password-encoder"
)

func EncryptPassword(pd string) string {
	options := &passwd.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt, codePwd := passwd.Encode(pd, options)
	passWord := fmt.Sprintf("pbkdf2-sha512$%s$%s", salt, codePwd)
	return passWord
}

func UnencryptPassword(raw, salt, encryptPassword *string) bool {
	opt := &passwd.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	return passwd.Verify(*raw, *salt, *encryptPassword, opt)
}
