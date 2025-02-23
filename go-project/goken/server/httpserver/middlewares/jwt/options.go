package jwt

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type JwtOption func(*GinJWTMiddleware)

func new(opts ...JwtOption) (*GinJWTMiddleware, error) {
	mw := &GinJWTMiddleware{
		TokenLookup:      "header:Authorization",
		SigningAlgorithm: "HS256",
		Timeout:          time.Hour * 24,
		MaxRefresh:       time.Hour * 24,
		TokenHeadName:    "Bearer",
		IdentityKey:      "identity",
		Realm:            "gin-jwt",
		CookieMaxAge:     time.Hour * 24,
		CookieName:       "jwt-cookie",
		ExpField:         "exp",
		SendCookie:       false,
		SecureCookie:     true,
		CookieHTTPOnly:   false,
		DisabledAbort:    false,
		TimeFunc:         time.Now,
		CookieSameSite:   http.SameSiteDefaultMode,
		Authorizator:     func(data interface{}, c *gin.Context) bool { return true },
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			c.JSON(http.StatusOK, gin.H{
				"code":   http.StatusOK,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		},
		LogoutResponse: func(c *gin.Context, code int) {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
			})
		},
		PayloadFunc: func(keyvalue ...interface{}) jwt.MapClaims {
			c := make(jwt.MapClaims)
			//每隔一个设置一对键值对
			for i := 0; i < len(keyvalue)-1; i += 2 {
				if key, ok := keyvalue[i].(string); ok {
					c[key] = keyvalue[i+1]
				}
			}
			return c
		},
		HTTPStatusMessageFunc: func(e error, c *gin.Context) string {
			return e.Error()
		},
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

	if mw.Authenticator == nil {
		return nil, ErrMissingAuthenticatorFunc
	}

	if mw.IdentityHandler == nil {
		mw.IdentityHandler = func(ctx *gin.Context) interface{} {
			return mw.IdentityKey
		}
	}

	//如果有能获得Key的Func就无需再判断有没有Key了
	if mw.KeyFunc != nil {
		return mw, nil
	}

	if mw.Key == nil {
		return nil, ErrMissingSecretKey
	}

	return mw, nil
}

func WithTokenLookup(tokenLookup string) JwtOption {
	return func(mw *GinJWTMiddleware) {
		mw.TokenLookup = tokenLookup
	}
}

func WithSecureKey(key string) JwtOption {
	return func(mw *GinJWTMiddleware) {
		mw.Key = []byte(key)
	}
}

func WithAuthenticator(Authenticator func(ctx *gin.Context) (interface{}, error)) JwtOption {
	return func(mw *GinJWTMiddleware) {
		mw.Authenticator = Authenticator
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

func WithTokenHeadName(tokenHeadName string) JwtOption {
	return func(mw *GinJWTMiddleware) {
		mw.TokenHeadName = strings.TrimSpace(mw.TokenHeadName)
	}
}

func WithAuthorizator(authorizator func(data interface{}, c *gin.Context) bool) JwtOption {
	return func(mw *GinJWTMiddleware) {
		mw.Authorizator = authorizator
	}
}

func WithUnauthorized(unauthorized func(c *gin.Context, code int, message string)) JwtOption {
	return func(mw *GinJWTMiddleware) {
		mw.Unauthorized = unauthorized
	}
}

func WithPayloadFunc(payloadFunc func(keyvalue ...interface{}) jwt.MapClaims) JwtOption {
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

func WithIdentityKey(identityKey string) JwtOption {
	return func(mw *GinJWTMiddleware) {
		mw.IdentityKey = identityKey
	}
}

func WithIdentityHandler(identityHandler func(c *gin.Context) interface{}) JwtOption {
	return func(mw *GinJWTMiddleware) {
		mw.IdentityHandler = identityHandler
	}
}

func WithHTTPStatusMessageFunc(httpStatusMessageFunc func(e error, c *gin.Context) string) JwtOption {
	return func(mw *GinJWTMiddleware) {
		mw.HTTPStatusMessageFunc = httpStatusMessageFunc
	}
}

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

func WithKey(key []byte) JwtOption {
	return func(mw *GinJWTMiddleware) {
		mw.Key = key
	}
}
