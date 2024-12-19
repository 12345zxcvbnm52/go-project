package pkg

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	ID       string
	UserName string
	jwt.RegisteredClaims
}

var (
	ErrTokenExpired      = errors.New("Token已过期")
	ErrTokenMalformed    = errors.New("Token认证格式错误,这或许是一个残缺的Token")
	ErrTokenInvalid      = errors.New("无效的Token")
	ErrTokenNotActiveYet = errors.New("Token还没有生效")
	ErrTokenUndefined    = errors.New("未知的Token错误,建议检查Token是否正确")
)

type JwtAuth struct {
	jwtSign []byte
}

func JWTAuthCreate(claim CustomClaims, sign string) (string, error) {
	jwtAuth := JwtAuth{jwtSign: []byte(sign)}
	return jwtAuth.CreateToken(claim)
}

func JWTAuthCheck(tokenStr string, sign string) (*CustomClaims, error) {
	jwtAuth := JwtAuth{jwtSign: []byte(sign)}
	return jwtAuth.ParseToken(tokenStr)
}

func (j *JwtAuth) CreateToken(claim CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(j.jwtSign)
}

func (j *JwtAuth) ParseToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return j.jwtSign, nil
	})
	if err != nil {
		//细分出现的错误
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrTokenNotActiveYet
			} else {
				return nil, ErrTokenInvalid
			}
		} else {
			return nil, err
		}
	}
	//这里加个保险
	if token != nil {
		if claims, ok := token.Claims.(CustomClaims); ok && token.Valid {
			return &claims, nil
		}
	}
	return nil, ErrTokenUndefined
}
