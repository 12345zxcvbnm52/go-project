package encrypt

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"kenshop/pkg/errors"

	"golang.org/x/crypto/sha3"

	"github.com/spaolacci/murmur3"
)

const defaultHashAlgorithm = "sha512"

var (
	ErrUndefinedHashFunc = errors.New("未定义的HashFunc")
)

var (
	SHA1      = "sha1"
	SHA256    = "sha256"
	SHA512    = "sha512"
	SHA3_256  = "sha3-256"
	SHA3_512  = "sha3-512"
	MD5       = "md5"
	Murmur32  = "murmur32"
	Murmur64  = "murmur64"
	Murmur128 = "murmur128"
)

// 该函数不会返回错误,但如果参数有误则默认使用sha512算法,返回值中string为算法名,避免出现参数有误
func NewHashFuncWithDefault(algorithm string) (hash.Hash, string) {
	switch algorithm {
	case SHA1:
		return sha1.New(), algorithm
	case SHA256:
		return sha256.New(), algorithm
	case SHA512:
		return sha512.New(), algorithm
	case SHA3_256:
		return sha3.New256(), algorithm
	case SHA3_512:
		return sha3.New512(), algorithm
	case MD5:
		return md5.New(), algorithm
	case Murmur64:
		return murmur3.New64(), algorithm
	case Murmur128:
		return murmur3.New128(), algorithm
	case "", Murmur32:
		return murmur3.New32(), algorithm
	default:
		return sha512.New(), defaultHashAlgorithm
	}
}

func NewHashFunc(algorithm string) (hash.Hash, error) {
	switch algorithm {
	case SHA1:
		return sha1.New(), nil
	case SHA256:
		return sha256.New(), nil
	case SHA512:
		return sha512.New(), nil
	case SHA3_256:
		return sha3.New256(), nil
	case SHA3_512:
		return sha3.New512(), nil
	case MD5:
		return md5.New(), nil
	case Murmur64:
		return murmur3.New64(), nil
	case Murmur128:
		return murmur3.New128(), nil
	case "", Murmur32:
		return murmur3.New32(), nil
	default:
		return nil, ErrUndefinedHashFunc
	}
}

// NewHashFuncForPBKDF2 生成适用于PBKDF2的哈希函数
func NewHashFuncForPBKDF2(algorithm string) (func() hash.Hash, error) {
	switch algorithm {
	case SHA1:
		return sha1.New, nil
	case SHA256:
		return sha256.New, nil
	case SHA512:
		return sha512.New, nil
	case SHA3_256:
		return sha3.New256, nil
	case SHA3_512:
		return sha3.New512, nil
	default:
		return nil, ErrUndefinedHashFunc
	}
}

// NewHashFuncForPBKDF2 生成适用于PBKDF2的哈希函数
func NewHashFuncForPBKDF2WithDefault(algorithm string) (func() hash.Hash, string) {
	switch algorithm {
	case SHA1:
		return sha1.New, algorithm
	case SHA256:
		return sha256.New, algorithm
	case SHA512:
		return sha512.New, algorithm
	case SHA3_256:
		return sha3.New256, algorithm
	case SHA3_512:
		return sha3.New512, algorithm
	default:
		return sha512.New, defaultHashAlgorithm
	}
}
