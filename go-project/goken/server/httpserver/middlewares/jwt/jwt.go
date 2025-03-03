package jwt

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// gin-jwt中间件的对外使用api,该中间件会对每一个被拦截的请求都检测jwt是否认证成功
// 亦包括对用户权限的检测
func (mw *GinJWTMiddleware) JwtAuthHandler(c *gin.Context) {
	claims, err := mw.GetClaimsFromContext(c)
	if err != nil {
		mw.unauthorized(c, http.StatusUnauthorized, err.Error())
		return
	}

	//判断Token是否Expire
	switch v := claims[mw.ExpField].(type) {
	case nil:
		mw.unauthorized(c, http.StatusBadRequest, ErrMissingExpField.Error())
		return
	case float64:
		if int64(v) < mw.TimeFunc().Unix() {
			mw.unauthorized(c, http.StatusUnauthorized, ErrExpiredToken.Error())
			return
		}
	case json.Number:
		n, err := v.Int64()
		if err != nil {
			mw.unauthorized(c, http.StatusBadRequest, ErrWrongFormatOfExp.Error())
			return
		}
		if n < mw.TimeFunc().Unix() {
			mw.unauthorized(c, http.StatusUnauthorized, ErrExpiredToken.Error())
			return
		}
	default:
		mw.unauthorized(c, http.StatusBadRequest, ErrWrongFormatOfExp.Error())
		return
	}

	//获得登入者claims和身份并保存到context中方便后续使用
	c.Set("JWT_PAYLOAD", claims)
	c.Next()
}

// RefreshHandler作为middleware可用于验证,刷新token,刷新的token仍然是有效的
// 刷新策略是生成的新token字符串会先放到Context中,最终会在
func (mw *GinJWTMiddleware) RefreshHandler(c *gin.Context) {

	tokenString, _, err := mw.RefreshToken(c)
	if err != nil {
		mw.unauthorized(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Set("refresh-token", tokenString)
	c.Header("refresh-token", tokenString)
}

// 刷新token并检查token是否已过期
func (mw *GinJWTMiddleware) RefreshToken(c *gin.Context) (string, time.Time, error) {
	claims, err := mw.CheckIfTokenExpire(c)
	//如果refresh token仍旧超时则直接返回错误
	if err != nil {
		return "", time.Now(), err
	}

	// 创建一个新的Token,考虑到安全性还是生成新的token而不是复用
	newToken := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	newClaims := newToken.Claims.(jwt.MapClaims)

	for key := range claims {
		newClaims[key] = claims[key]
	}

	expire := mw.TimeFunc().Add(mw.TimeoutFunc(claims))
	newClaims[mw.ExpField] = expire.Unix()
	newClaims["orig_iat"] = mw.TimeFunc().Unix()
	newClaims["nbf"] = claims["orig_iat"]
	newClaims["iat"] = claims["orig_iat"]
	newClaims["rfs_exp"] = expire.Add(mw.MaxRefreshFunc(claims)).Unix()
	c.Set("JWT_PAYLOAD", newClaims)

	tokenString, err := newToken.SignedString(mw.Key)
	if err != nil {
		return "", time.Now(), err
	}
	mw.SetCookie(c, tokenString)
	return tokenString, expire, nil
}

// 检查Token是否超时
func (mw *GinJWTMiddleware) CheckIfTokenExpire(c *gin.Context) (jwt.MapClaims, error) {
	token, err := mw.getTokenFromCtx(c)
	if err != nil {
		// 如果收到一个错误,且该错误不是ValidationErrorExpired,则返回该错误,
		// 如果错误只是 ValidationErrorExpired,则继续执行直至检查是否能refresh

		validationErr, ok := err.(*jwt.ValidationError)
		if !ok || validationErr.Errors != jwt.ValidationErrorExpired {
			return nil, err
		}
	}

	claims := token.Claims.(jwt.MapClaims)
	expIat := claims["rfs_exp"].(float64)

	if expIat < float64(mw.TimeFunc().Unix()) {
		return nil, ErrExpiredRefreshToken
	}

	return claims, nil
}
