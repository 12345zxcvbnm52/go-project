package jwt

import (
	"encoding/json"
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

	// 回调函数,应该根据登录信息执行用户认证,这个函数不会默认生成,必须提供一个认证函数,
	//必须返回用户数据作为用户标识符,该标识符将被存储在Claim数组中,
	Authenticator func(c *gin.Context) (interface{}, error)

	// 回调函数,在用户认证成功后调用,用于判断用户是否有权限访问某些资源,默认返回true(表示始终允许访问)
	Authorizator func(data interface{}, c *gin.Context) bool

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

	// 回调函数,用于获取身份标识(如用户ID),
	//该函数用于从请求中提取身份信息,并将其作为userID提供给上下文,
	IdentityHandler func(*gin.Context) interface{}

	// 设置身份键
	IdentityKey string

	// TokenFormat 是一个字符串格式为 "<source>:<x-token>"用于从请求中提取 token,
	// 可选,默认值为 "header:Authorization",
	// 可能的值:
	// - "header:x-token"
	// - "query:x-token"
	// - "cookie:x-token"
	TokenFormat string

	// TokenHeadName是header中标识Token字段的字符串,默认值为"Bearer"
	TokenHeadName string

	// TimeFunc提供当前时间,主要用于不同时区之间的连接
	//可以覆盖它以使用其他时间值,这对于测试或如果服务器使用与token不同的时区非常有用,
	TimeFunc func() time.Time

	// 动态返回JWT中间件中内容失败时的HTTP状态错误消息,
	HTTPStatusMessageFunc func(e error, c *gin.Context) string

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

	// 禁用context的abort(),
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

	// ErrMissingAuthenticatorFunc indicates Authenticator is required
	ErrMissingAuthenticatorFunc = errors.New("ginJWTMiddleware.Authenticator func is undefined")

	// ErrMissingLoginValues indicates a user tried to authenticate without username or password
	ErrMissingLoginValues = errors.New("missing Username or Password")

	// ErrFailedAuthentication indicates authentication failed, could be faulty username or password
	ErrFailedAuthentication = errors.New("incorrect Username or Password")

	// ErrFailedTokenCreation indicates JWT Token failed to create, reason unknown
	ErrFailedTokenCreation = errors.New("failed to create JWT Token")

	// ErrExpiredToken indicates JWT token has expired. Can't refresh.
	ErrExpiredToken = errors.New("token is expired") // in practice, this is generated from the jwt library not by us

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

	// ErrNoPrivKeyFile indicates that the given private key is unreadable
	ErrNoPrivKeyFile = errors.New("private key file unreadable")

	// ErrNoPubKeyFile indicates that the given public key is unreadable
	ErrNoPubKeyFile = errors.New("public key file unreadable")

	// ErrInvalidPrivKey indicates that the given private key is invalid
	ErrInvalidPrivKey = errors.New("private key invalid")
)

// New for check error with GinJWTMiddleware
func NewGinJWTMiddleware(opts ...JwtOption) (*GinJWTMiddleware, error) {
	return new(opts...)
}

// gin-jwt中间件的对外使用api,该中间件会对每一个被拦截的请求都检测jwt是否认证成功
// 亦包括对用户权限的检测
func (mw *GinJWTMiddleware) MiddlewareFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		mw.middlewareImpl(c)
	}
}

func (mw *GinJWTMiddleware) middlewareImpl(c *gin.Context) {
	claims, err := mw.GetClaimsFromContext(c)
	if err != nil {
		mw.unauthorized(c, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(err, c))
		return
	}

	//判断Token是否Expire
	switch v := claims[mw.ExpField].(type) {
	case nil:
		mw.unauthorized(c, http.StatusBadRequest, mw.HTTPStatusMessageFunc(ErrMissingExpField, c))
		return
	case float64:
		if int64(v) < mw.TimeFunc().Unix() {
			mw.unauthorized(c, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(ErrExpiredToken, c))
			return
		}
	case json.Number:
		n, err := v.Int64()
		if err != nil {
			mw.unauthorized(c, http.StatusBadRequest, mw.HTTPStatusMessageFunc(ErrWrongFormatOfExp, c))
			return
		}
		if n < mw.TimeFunc().Unix() {
			mw.unauthorized(c, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(ErrExpiredToken, c))
			return
		}
	default:
		mw.unauthorized(c, http.StatusBadRequest, mw.HTTPStatusMessageFunc(ErrWrongFormatOfExp, c))
		return
	}

	//获得登入者claims和身份并保存到context中方便后续使用
	c.Set("JWT_PAYLOAD", claims)
	identity := mw.IdentityHandler(c)

	if identity != nil {
		c.Set(mw.IdentityKey, identity)
	}

	//检测该用户是否允许访问该资源
	if !mw.Authorizator(identity, c) {
		mw.unauthorized(c, http.StatusForbidden, mw.HTTPStatusMessageFunc(ErrForbidden, c))
		return
	}

	c.Next()
}

// LoginHandler作为Middleware能用于客户端获取jwt-token
// 有效载荷为JSON格式,形如{"username": "用户名", "password": "密码"},
// 响应将形如{"token": "令牌"},
func (mw *GinJWTMiddleware) LoginHandler(c *gin.Context) {
	data, err := mw.Authenticator(c)
	if err != nil {
		mw.unauthorized(c, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(err, c))
		return
	}

	tokenString, expire, err := mw.NewToken("id", data)
	if err != nil {
		mw.unauthorized(c, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(ErrFailedTokenCreation, c))
		return
	}

	mw.SetCookie(c, tokenString)
	mw.LoginResponse(c, http.StatusOK, tokenString, expire)
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

func (mw *GinJWTMiddleware) signedString(token *jwt.Token) (string, error) {
	var tokenString string
	var err error
	tokenString, err = token.SignedString(mw.Key)

	return tokenString, err
}

// RefreshHandler作为middleware可用于验证,刷新token,刷新的token仍然是有效的
// 刷新策略是生成的新token字符串会先放到Context中,最终会在
func (mw *GinJWTMiddleware) RefreshHandler(c *gin.Context) {

	tokenString, _, err := mw.RefreshToken(c)
	if err != nil {
		mw.unauthorized(c, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(err, c))
		return
	}
	c.Set("refresh-token", tokenString)
	c.Header("refresh-token", tokenString)
}

// 刷新token并检查token是否已过期
func (mw *GinJWTMiddleware) RefreshToken(c *gin.Context) (string, time.Time, error) {
	claims, err := mw.CheckIfTokenExpire(c)
	//如果哪怕refresh了token仍旧超时则返回错误
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
	claims["rfs_exp"] = expire.Add(mw.MaxRefreshFunc(claims)).Unix()

	tokenString, err := mw.signedString(newToken)
	if err != nil {
		return "", time.Now(), err
	}

	mw.SetCookie(c, tokenString)

	return tokenString, expire, nil
}

// 检查Token是否超时
func (mw *GinJWTMiddleware) CheckIfTokenExpire(c *gin.Context) (jwt.MapClaims, error) {
	token, err := mw.ParseToken(c)
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
		return nil, ErrExpiredToken
	}

	return claims, nil
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

	fmt.Println(expire.Add(mw.MaxRefreshFunc(claims)).Unix())

	claims["rfs_exp"] = expire.Add(mw.MaxRefreshFunc(claims)).Unix()

	if mw.PayloadFunc != nil {
		for key, value := range mw.PayloadFunc(keyvalue...) {
			claims[key] = value
		}
	}

	tokenString, err := mw.signedString(token)
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
func (mw *GinJWTMiddleware) ParseToken(c *gin.Context) (*jwt.Token, error) {
	var token string
	var err error

	methods := strings.Split(mw.TokenFormat, ",")
	for _, method := range methods {
		if len(token) > 0 {
			break
		}
		//如果存在refresh-token就直接使用refresh-token
		if v, exist := c.Get("refresh-token"); exist {
			key := v.(string)
			token = string(key)
		} else {
			parts := strings.Split(strings.TrimSpace(method), ":")
			k := strings.TrimSpace(parts[0])
			v := strings.TrimSpace(parts[1])
			switch k {
			case "header":
				token, err = mw.jwtFromHeader(c, v)
			case "query":
				token, err = mw.jwtFromQuery(c, v)
			case "cookie":
				token, err = mw.jwtFromCookie(c, v)
			case "param":
				token, err = mw.jwtFromParam(c, v)
			case "form":
				token, err = mw.jwtFromForm(c, v)
			}
		}
	}

	if err != nil {
		return nil, err
	}

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

// 将TokenStr反序列为jwt.Token
func (mw *GinJWTMiddleware) ParseTokenString(token string) (*jwt.Token, error) {
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

func (mw *GinJWTMiddleware) unauthorized(c *gin.Context, code int, message string) {
	c.Header("WWW-Authenticate", "JWT realm="+mw.Realm)
	if !mw.DisabledAbort {
		c.Abort()
	}
	mw.Unauthorized(c, code, message)
}

// 从gin.Context中得到MapClaims
func (mw *GinJWTMiddleware) GetClaimsFromContext(c *gin.Context) (jwt.MapClaims, error) {
	token, err := mw.ParseToken(c)
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
