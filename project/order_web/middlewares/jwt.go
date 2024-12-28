package middlewares

import (
	"errors"
	"net/http"
	gb "order_web/global"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	ID       uint32
	UserName string
	jwt.RegisteredClaims
}

/*
	jwt的数据应当去中心化保存在客户端手中,每次要检查时直接读取过期时间对比即可
*/

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

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Request.Header.Get("x-token")
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "请先登入账号",
			})
			c.Abort()
			return
		}
		j := NewJwtAuth()
		claims, err := j.ParseToken(tokenStr)
		if err != nil {
			//可以考虑双重jwt保证用户不会在使用中jwt过期
			if err == ErrTokenExpired {
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"msg":  "授权登录出错",
					"详细信息": err.Error(),
				})
				c.Abort()
				return
			}
		}
		c.Set("claims", claims)
		c.Next()
	}
}

func NewJwtAuth() *JwtAuth {
	return &JwtAuth{jwtSign: []byte(gb.ServerConfig.JwtSign)}
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
