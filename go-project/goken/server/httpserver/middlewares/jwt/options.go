package jwt

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type JwtOption func(*GinJWTMiddleware)

func MustNewGinJWTMiddleware(key string, opts ...JwtOption) *GinJWTMiddleware {
	mw := &GinJWTMiddleware{
		TokenInside:      "header",
		TokenHeadName:    "Authorization",
		SigningAlgorithm: "HS256",
		Timeout:          time.Hour * 24,
		MaxRefresh:       time.Hour * 24,
		Realm:            "gin-jwt",
		CookieMaxAge:     time.Hour * 24,
		CookieName:       "jwt-cookie",
		ExpField:         "exp",
		SendCookie:       false,
		SecureCookie:     false,
		CookieHTTPOnly:   false,
		DisabledAbort:    false,
		TimeFunc:         time.Now,
		Key:              []byte(key),
		CookieSameSite:   http.SameSiteDefaultMode,
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":  code,
				"error": message,
			})
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(code, gin.H{
				"code":   code,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		},
		LogoutResponse: func(c *gin.Context, code int) {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
			})
		},
		PayloadFunc: func(keyvalue ...string) jwt.MapClaims {
			c := make(jwt.MapClaims)
			//每隔一个设置一对键值对
			for i := 0; i < len(keyvalue)-1; i += 2 {
				c[keyvalue[i]] = keyvalue[i+1]
			}
			return c
		},
		// HTTPStatusMessageFunc: func(e error, c *gin.Context) string {
		// 	return e.Error()
		// },
	}

	// 应用所有选项函数
	for _, JwtOption := range opts {
		JwtOption(mw)
	}

	if mw.TimeoutFunc == nil {
		mw.TimeoutFunc = func(data interface{}) time.Duration {
			return mw.Timeout
		}
	}

	if mw.MaxRefreshFunc == nil {
		mw.MaxRefreshFunc = func(data interface{}) time.Duration {
			return mw.MaxRefresh
		}
	}

	//如果有能获得Key的Func就无需再判断有没有Key了
	if mw.KeyFunc != nil {
		return mw
	}

	if mw.Key == nil {
		panic(ErrMissingSecretKey)
	}

	return mw
}

func WithTokenInside(t string) JwtOption {
	return func(mw *GinJWTMiddleware) {
		mw.TokenInside = t
	}
}

func WithTokenHeadName(h string) JwtOption {
	return func(mw *GinJWTMiddleware) {
		mw.TokenHeadName = h
	}
}

func WithSigningAlgorithm(algorithm string) JwtOption {
	return func(mw *GinJWTMiddleware) {
		mw.SigningAlgorithm = algorithm
	}
}

func WithTimeout(timeout time.Duration) JwtOption {
	return func(mw *GinJWTMiddleware) {
		mw.Timeout = timeout
	}
}

func WithTimeoutFunc(timeoutFunc func(data interface{}) time.Duration) JwtOption {
	return func(mw *GinJWTMiddleware) {
		mw.TimeoutFunc = timeoutFunc
	}
}

func WithTimeFunc(timeFunc func() time.Time) JwtOption {
	return func(mw *GinJWTMiddleware) {
		mw.TimeFunc = timeFunc
	}
}

func WithUnauthorized(unauthorized func(c *gin.Context, code int, message string)) JwtOption {
	return func(mw *GinJWTMiddleware) {
		mw.Unauthorized = unauthorized
	}
}

func WithPayloadFunc(payloadFunc func(keyvalue ...string) jwt.MapClaims) JwtOption {
	return func(mw *GinJWTMiddleware) {
		mw.PayloadFunc = payloadFunc
	}
}

func WithLoginResponse(loginResponse func(c *gin.Context, code int, token string, expire time.Time)) JwtOption {
	return func(mw *GinJWTMiddleware) {
		mw.LoginResponse = loginResponse
	}
}

func WithLogoutResponse(logoutResponse func(c *gin.Context, code int)) JwtOption {
	return func(mw *GinJWTMiddleware) {
		mw.LogoutResponse = logoutResponse
	}
}

// func WithHTTPStatusMessageFunc(httpStatusMessageFunc func(e error, c *gin.Context) string) JwtOption {
// 	return func(mw *GinJWTMiddleware) {
// 		mw.HTTPStatusMessageFunc = httpStatusMessageFunc
// 	}
// }

func WithRealm(realm string) JwtOption {
	return func(mw *GinJWTMiddleware) {
		mw.Realm = realm
	}
}

func WithCookieMaxAge(cookieMaxAge time.Duration) JwtOption {
	return func(mw *GinJWTMiddleware) {
		mw.CookieMaxAge = cookieMaxAge
	}
}

func WithCookieName(cookieName string) JwtOption {
	return func(mw *GinJWTMiddleware) {
		mw.CookieName = cookieName
	}
}

func WithExpField(expField string) JwtOption {
	return func(mw *GinJWTMiddleware) {
		mw.ExpField = expField
	}
}

func WithKeyFunc(keyFunc func(t *jwt.Token) (interface{}, error)) JwtOption {
	return func(mw *GinJWTMiddleware) {
		mw.KeyFunc = keyFunc
	}
}
