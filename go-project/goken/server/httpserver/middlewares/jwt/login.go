package jwt

import (
	"kenshop/pkg/common/httputil"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginHandler作为Middleware能用于客户端获取jwt-token
// 有效载荷为JSON格式,形如{"username": "用户名", "password": "密码"},
// 响应将形如{"token": "令牌"},
// 回调函数,应该根据登录信息执行用户认证,这个函数不会默认生成,必须提供一个认证函数,
// 必须返回用户数据作为用户标识符,该标识符将被存储在Claim数组中,
func (mw *GinJWTMiddleware) LoginHandler(Authenticator func(*gin.Context) bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ok := Authenticator(ctx); ok {
			tokenString, expire, err := mw.NewToken()
			if err != nil {
				httputil.WriteResponse(ctx, http.StatusUnauthorized, err.Error(), nil, mw.UseAbort)
			}
			mw.SetCookie(ctx, tokenString)
			httputil.WriteResponse(ctx, http.StatusOK, "", gin.H{"token": tokenString, "expire": expire}, mw.UseAbort)
		} else {
			httputil.WriteResponse(ctx, http.StatusUnauthorized, ErrFailedAuthentication.Error(), nil, mw.UseAbort)
		}
	}
}

// LogoutHandler作为Middleware可供客户端使用,用于移除jwt的cookie
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
			"",
			mw.SecureCookie,
			mw.SecureCookie,
		)
	}
}
