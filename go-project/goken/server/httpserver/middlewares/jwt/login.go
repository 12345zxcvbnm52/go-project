package jwt

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginHandler作为Middleware能用于客户端获取jwt-token
// 有效载荷为JSON格式,形如{"username": "用户名", "password": "密码"},
// 响应将形如{"token": "令牌"},
// 回调函数,应该根据登录信息执行用户认证,这个函数不会默认生成,必须提供一个认证函数,
// 必须返回用户数据作为用户标识符,该标识符将被存储在Claim数组中,
func (mw *GinJWTMiddleware) LoginHandler(Authenticator func(c *gin.Context) (map[string]string, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := Authenticator(ctx)
		if err != nil {
			mw.unauthorized(ctx, http.StatusUnauthorized, err.Error())
			return
		}
		kv := make([]string, 0)
		for k, v := range data {
			kv = append(kv, k, v)
		}
		tokenString, expire, err := mw.NewToken(kv...)
		if err != nil {
			mw.unauthorized(ctx, http.StatusUnauthorized, err.Error())
			return
		}

		mw.SetCookie(ctx, tokenString)
		mw.LoginResponse(ctx, http.StatusOK, tokenString, expire)
	}
}

// LogoutHandler作为Middleware可供客户端使用,
// 用于移除jwt的cookie以及执行内部的LogoutResponse回调函数
func (mw *GinJWTMiddleware) LogoutHandler(c *gin.Context) {
	if mw.SendCookie {
		if mw.CookieSameSite != 0 {
			c.SetSameSite(mw.CookieSameSite)
		}

		//通过设置空值移除cookie
		c.SetCookie(
			mw.CookieName,
			"",
			-1,
			"/",
			mw.CookieDomain,
			mw.SecureCookie,
			mw.CookieHTTPOnly,
		)
	}

	mw.LogoutResponse(c, http.StatusOK)
}
