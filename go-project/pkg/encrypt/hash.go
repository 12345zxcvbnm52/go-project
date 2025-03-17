package encrypt

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"kenshop/pkg/errors"
	"strings"

	"github.com/buger/jsonparser"
	uuid "github.com/google/uuid"
	"golang.org/x/crypto/sha3"

	"github.com/spaolacci/murmur3"
)

const defaultHashAlgorithm = "sha512"

// GenerateToken 生成token,如果哈希算法为空,则不会加密而是直接裸露连接
func GenerateToken(orgID, keyID, hashAlgorithm string) (string, error) {
	if keyID == "" {
		keyID = strings.ReplaceAll(uuid.Must(uuid.NewV7()).String(), "-", "")
	}

	if hashAlgorithm != "" {
		_, err := hashFunction(hashAlgorithm)
		if err != nil {
			hashAlgorithm = defaultHashAlgorithm
		}

		jsonToken := fmt.Sprintf(`{"org":"%s","id":"%s","h":"%s"}`, orgID, keyID, hashAlgorithm)

		return base64.StdEncoding.EncodeToString([]byte(jsonToken)), err
	}

	// Legacy keys
	return orgID + keyID, nil
}

// B64JSONPrefix stand for `{"` in base64.
const B64JSONPrefix = "ey"

// TokenHashAlgo ...
func TokenHashAlgo(token string) string {
	// Legacy tokens not b64 and not JSON records
	if strings.HasPrefix(token, B64JSONPrefix) {
		if jsonToken, err := base64.StdEncoding.DecodeString(token); err == nil {
			hashAlgo, _ := jsonparser.GetString(jsonToken, "h")

			return hashAlgo
		}
	}

	return ""
}

// TokenOrg ...
func TokenOrg(token string) string {
	if strings.HasPrefix(token, B64JSONPrefix) {
		if jsonToken, err := base64.StdEncoding.DecodeString(token); err == nil {
			// Checking error in case if it is a legacy tooken which just by accided has the same b64JSON prefix
			if org, err := jsonparser.GetString(jsonToken, "org"); err == nil {
				return org
			}
		}
	}

	// 24 is mongo bson id length
	if len(token) > 24 {
		return token[:24]
	}

	return ""
}

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

func hashFunction(algorithm string) (hash.Hash, error) {
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
		return nil, errors.Errorf("unknown key hash function: %s", algorithm)
	}
}

// HashStr return hash the give string and return.
func HashStr(in string) string {
	h, _ := hashFunction(TokenHashAlgo(in))
	_, _ = h.Write([]byte(in))

	return hex.EncodeToString(h.Sum(nil))
}

// HashKey return hash the give string and return.
func HashKey(in string) string {
	return HashStr(in)
}
