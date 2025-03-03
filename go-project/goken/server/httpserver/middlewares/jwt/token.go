package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

const (
	TokenInHeader = "header:x-token"
	TokenInQuery  = "query:x-token"
	TokenInCookie = "cookie:x-token"
	TokenInParam  = "param:x-token"
	TokenInForm   = "form:x-token"
)

// GinJWTMiddleware 提供了一个 Json-Web-Token 认证实现,失败时返回 401 HTTP 响应,
// 成功时调用该包装的中间件后可以通过c.Get("userID").(string)获取用户ID,
// 用户可以通过向LoginHandler发送json请求来获取token,然后需要在http-header的Authenticatio中传递该token,
// 例如:Authorization:Bearer XXX_TOKEN_XXX
type GinJWTMiddleware struct {
	// jwt所属的域名,也可以用于签发者姓名
	Realm string

	//受众,用于适配jwt的aud
	Audience []string

	// 签名算法-可能的值有HS256,HS384,HS512,RS256,RS384或RS512
	// 默认值HS256
	SigningAlgorithm string

	// 用于签名的密钥
	Key []byte

	// 用于动态获取签名密钥的回调函数,设置KeyFunc将绕过所有其他密钥设置,可适用于使用公私钥对的场景
	KeyFunc func(token *jwt.Token) (interface{}, error)

	// JWT token的有效时长,可选默认值为一天,
	Timeout time.Duration

	// 返回timeout时间,即源码内不直接使用Timeout字段,可以自定义逻辑处理Timeout
	TimeoutFunc func(data interface{}) time.Duration

	// 返回MaxRefresh时间,即源码内不直接使用MaxRefresh字段,可以自定义逻辑处理MaxRefresh
	MaxRefreshFunc func(data interface{}) time.Duration

	// 此字段允许客户端在MaxRefresh时间过去之前刷新其 token,
	// 默认为一天,传入数值为0则关闭Refresh模式
	MaxRefresh time.Duration

	WithCache bool

	// 在登录期间调用的回调函数,使用此函数可以向webtoken添加自定义的payload数据(负载数据),
	// 默认不会设置额外的PayloadFunc
	// 除此之外还可以通过c.Get("JWT_PAYLOAD")在请求期间得到这些数据,请注意payload在jwt存储中是不会加密的,
	PayloadFunc func(keyvalue ...string) jwt.MapClaims

	// 用户自定义的未授权时回调处理函数,默认功能类似于WriteErrorToResponse,
	//当认证失败时,调用该函数返回一个自定义的响应,默认返回401状态码和错误信息
	Unauthorized func(c *gin.Context, httpCode int, message string)

	// 用户自定义的登录响应回调函数,类似于AfterLogin钩子函数,
	//当用户登录成功并返回token时,会调用此函数来定制登录后的响应内容,
	LoginResponse func(c *gin.Context, httpCode int, message string, time time.Time)

	// 用户自定义的登出回调函数,类似于AfterLogout钩子函数,
	//当用户登录成功并返回token时,会调用此函数来定制登出后的响应内容,
	LogoutResponse func(c *gin.Context, httpCode int)

	//TokenInside是一个字符串,用于指定token在请求中的位置,默认值为"header",允许有多个值,用,分隔
	TokenInside string

	// TokenHeadName是header中标识Token字段的字符串,默认值为"Autorization",
	TokenHeadName string

	// TimeFunc提供当前时间,主要用于不同时区之间的连接
	//可以覆盖它以使用其他时间值,这对于测试或如果服务器使用与token不同的时区非常有用,
	TimeFunc func() time.Time

	// 动态返回JWT中间件中内容失败时的HTTP状态错误消息,
	//HTTPStatusMessageFunc func(e error, c *gin.Context) string

	// 可选地将token作为cookie返回
	SendCookie bool

	// cookie的有效时长,默认等于Timeout值,
	CookieMaxAge time.Duration

	// 允许在开发过程中通过http使用不安全的cookie
	SecureCookie bool

	// 允许在开发过程中客户端访问cookie
	CookieHTTPOnly bool

	// 允许在开发过程中更改cookie域名
	CookieDomain string

	// 禁用gin.Context的abort(),
	DisabledAbort bool

	// 允许在开发过程中更改cookie名称
	CookieName string

	// 允许使用 http.SameSite cookie参数(用于控制jwt-cookie的跨域传输),可选参数为:
	// SameSiteDefaultMode:默认行为,取决于浏览器,
	// SameSiteStrictMode:仅在同站请求中发送cookie,
	// SameSiteLaxMode:允许在跨域请求中发送cookie,
	CookieSameSite http.SameSite

	// 允许修改jwt的解析器方法
	ParseOptions []jwt.ParserOption

	// 默认值为"exp",是expire存储在MapClaims中的key
	ExpField string
}

var (
	// ErrMissingSecretKey indicates Secret key is required
	ErrMissingSecretKey = errors.New("secret key is required")

	// ErrForbidden when HTTP status 403 is given
	ErrForbidden = errors.New("you don't have permission to access this resource")

	// ErrFailedAuthentication indicates authentication failed, could be faulty username or password
	ErrFailedAuthentication = errors.New("incorrect Username or Password")

	// ErrFailedTokenCreation indicates JWT Token failed to create, reason unknown
	ErrFailedTokenCreation = errors.New("failed to create JWT Token")

	// ErrExpiredToken indicates JWT token has expired. Can't refresh.
	ErrExpiredToken = errors.New("token is expired") // in practice, this is generated from the jwt library not by us

	//ErrExpiredRefreshToken indicates JWT token has expired. Can't refresh.
	ErrExpiredRefreshToken = errors.New("refresh token is expired")

	// ErrEmptyToken can be thrown if authing with a HTTP header, the Auth header needs to be set
	ErrEmptyToken = errors.New("token header is empty")

	// ErrMissingExpField missing exp field in token
	ErrMissingExpField = errors.New("missing exp field")

	// ErrWrongFormatOfExp field must be float64 format
	ErrWrongFormatOfExp = errors.New("exp must be float64 format")

	// ErrInvalidToken indicates auth header is invalid, could for example have the wrong Realm name
	ErrInvalidToken = errors.New("auth header is invalid")

	// ErrEmptyQueryToken can be thrown if authing with URL Query, the query token variable is empty
	ErrEmptyQueryToken = errors.New("query token is empty")

	// ErrEmptyCookieToken can be thrown if authing with a cookie, the token cookie is empty
	ErrEmptyCookieToken = errors.New("cookie token is empty")

	// ErrEmptyParamToken can be thrown if authing with parameter in path, the parameter in path is empty
	ErrEmptyParamToken = errors.New("parameter token is empty")

	// ErrInvalidSigningAlgorithm indicates signing algorithm is invalid, needs to be HS256, HS384, HS512, RS256, RS384 or RS512
	ErrInvalidSigningAlgorithm = errors.New("invalid signing algorithm")
)

func (mw *GinJWTMiddleware) AuthorizationHandler(c *gin.Context) {
	//claims:=ExtractClaimsFromContext(c)
	//identity:=c.GetString(mw.IdentityKey)
	c.Next()
}

// 用于生成jwt.Token
func (mw *GinJWTMiddleware) NewToken(keyvalue ...string) (string, time.Time, error) {
	token := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)
	expire := mw.TimeFunc().Add(mw.TimeoutFunc(claims))

	claims[mw.ExpField] = expire.Unix()
	//orig-iat代表着jwt令牌的原始签发时间,这个时间不会因为refresh等更新
	claims["orig_iat"] = mw.TimeFunc().Unix()
	claims["iss"] = mw.Realm
	claims["nbf"] = claims["orig_iat"]
	claims["iat"] = claims["orig_iat"]
	claims["aud"] = mw.Audience
	claims["rfs_exp"] = expire.Add(mw.MaxRefreshFunc(claims)).Unix()
	if mw.PayloadFunc != nil {
		for key, value := range mw.PayloadFunc(keyvalue...) {
			claims[key] = value
		}
	}

	tokenString, err := token.SignedString(mw.Key)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expire, nil
}

func (mw *GinJWTMiddleware) jwtFromHeader(c *gin.Context, key string) (string, error) {
	token := c.Request.Header.Get(key)
	if token == "" {
		return "", ErrEmptyToken
	}
	return token, nil
}

func (mw *GinJWTMiddleware) jwtFromQuery(c *gin.Context, key string) (string, error) {
	token := c.Query(key)
	if token == "" {
		return "", ErrEmptyQueryToken
	}
	return token, nil
}

func (mw *GinJWTMiddleware) jwtFromCookie(c *gin.Context, key string) (string, error) {
	cookie, _ := c.Cookie(key)
	if cookie == "" {
		return "", ErrEmptyCookieToken
	}
	return cookie, nil
}

func (mw *GinJWTMiddleware) jwtFromParam(c *gin.Context, key string) (string, error) {
	token := c.Param(key)
	if token == "" {
		return "", ErrEmptyParamToken
	}
	return token, nil
}

func (mw *GinJWTMiddleware) jwtFromForm(c *gin.Context, key string) (string, error) {
	token := c.PostForm(key)
	if token == "" {
		return "", ErrEmptyParamToken
	}
	return token, nil
}

// 从gin.Context中获取jwt.Token
func (mw *GinJWTMiddleware) getTokenFromCtx(c *gin.Context) (*jwt.Token, error) {
	var token string
	var err error

	//如果存在refresh-token就直接使用refresh-token
	if v, exist := c.Get("refresh-token"); exist {
		key := v.(string)
		token = string(key)
	} else {
		inside := strings.Split(mw.TokenInside, ",")
		for _, v := range inside {
			if token != "" {
				break
			}
			switch v {
			case "header":
				token, err = mw.jwtFromHeader(c, mw.TokenHeadName)
			case "query":
				token, err = mw.jwtFromQuery(c, mw.TokenHeadName)
			case "cookie":
				token, err = mw.jwtFromCookie(c, mw.TokenHeadName)
			case "param":
				token, err = mw.jwtFromParam(c, mw.TokenHeadName)
			case "form":
				token, err = mw.jwtFromForm(c, mw.TokenHeadName)
			}
		}
	}
	if err != nil {
		return nil, err
	}

	return mw.parseTokenString(token)
}

// 将TokenStr反序列为jwt.Token
func (mw *GinJWTMiddleware) parseTokenString(token string) (*jwt.Token, error) {
	if mw.KeyFunc != nil {
		return jwt.Parse(token, mw.KeyFunc, mw.ParseOptions...)
	}

	return jwt.Parse(token,
		func(t *jwt.Token) (interface{}, error) {
			if jwt.GetSigningMethod(mw.SigningAlgorithm) != t.Method {
				return nil, ErrInvalidSigningAlgorithm
			}

			return mw.Key, nil
		},
		mw.ParseOptions...,
	)
}

// 从gin.Context中得到MapClaims
func (mw *GinJWTMiddleware) GetClaimsFromContext(c *gin.Context) (jwt.MapClaims, error) {
	token, err := mw.getTokenFromCtx(c)
	if err != nil {
		return nil, err
	}

	claims := jwt.MapClaims{}
	for key, value := range token.Claims.(jwt.MapClaims) {
		claims[key] = value
	}

	return claims, nil
}

// 从gin.Context中反解出对应的MapClaims
func ExtractClaimsFromContext(c *gin.Context) jwt.MapClaims {
	claims, exists := c.Get("JWT_PAYLOAD")
	if !exists {
		return make(jwt.MapClaims)
	}

	return claims.(jwt.MapClaims)
}

// 从gin.Token中反解出对应的MapClaims
func ExtractClaimsFromToken(token *jwt.Token) jwt.MapClaims {
	if token == nil {
		return make(jwt.MapClaims)
	}

	claims := jwt.MapClaims{}
	for key, value := range token.Claims.(jwt.MapClaims) {
		claims[key] = value
	}

	return claims
}

// 将Token加入到Cookie,只在SendCookie开启时有效
func (mw *GinJWTMiddleware) SetCookie(c *gin.Context, token string) {
	if mw.SendCookie {
		expireCookie := mw.TimeFunc().Add(mw.CookieMaxAge)
		maxage := int(expireCookie.Unix() - mw.TimeFunc().Unix())

		if mw.CookieSameSite != 0 {
			c.SetSameSite(mw.CookieSameSite)
		}

		c.SetCookie(
			mw.CookieName,
			token,
			maxage,
			"/",
			mw.CookieDomain,
			mw.SecureCookie,
			mw.CookieHTTPOnly,
		)
	}
}

func (mw *GinJWTMiddleware) unauthorized(c *gin.Context, code int, message string) {
	c.Header("WWW-Authenticate", fmt.Sprintf("JWT realm=%s", mw.Realm))
	if !mw.DisabledAbort {
		c.Abort()
	}
	mw.Unauthorized(c, code, message)
}
