package jwt

import (
	"kenshop/pkg/cache"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// GinJWTMiddleware 提供了一个 Json-Web-Token 认证实现,失败时返回 401 HTTP 响应,
// 成功时调用该包装的中间件后可以通过c.Get("userID").(string)获取用户ID,
// 用户可以通过向LoginHandler发送json请求来获取token,然后需要在http-header的Authentication中传递该token,
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

	Cache *cache.MultiCache

	//TokenInside是一个字符串,用于指定token在请求中的位置,默认值为"header",允许有多个值,用,分隔
	TokenInside string

	// TokenHeadName是header中标识Token字段的字符串,默认值为"Autorization",
	TokenHeadName string

	// TimeFunc提供当前时间,主要用于不同时区之间的连接
	//可以覆盖它以使用其他时间值,这对于测试或如果服务器使用与token不同的时区非常有用,
	TimeFunc func() time.Time

	// 是否使用gin.Context的abort(),
	UseAbort bool

	// 可选地将token作为cookie返回,默认关闭
	SendCookie bool

	// cookie的有效时长,默认等于Timeout值,
	CookieMaxAge time.Duration

	// 允许在使用不安全的cookie,默认开启
	SecureCookie bool

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
	// 允许最多5秒时间偏差
	claims["nbf"] = mw.TimeFunc().Add(-time.Second * 5).Unix()
	claims["iat"] = claims["orig_iat"]
	claims["aud"] = mw.Audience
	refreshExpire := expire.Add(mw.MaxRefreshFunc(claims))
	claims["rfs_exp"] = refreshExpire.Unix()

	//把自定义的kv perload加入其中
	for i := 0; i+1 < len(keyvalue); i += 2 {
		claims[keyvalue[i]] = keyvalue[i+1]
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
	t, err := mw.parseTokenString(token)
	if err != nil {
		return nil, ErrTokenStringInvalid
	}
	return t, nil
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

// 从gin.Context中得到MapClaims,如果ctx中存在refresh-token(jwt被refresh了)则优先从此中找到token和claims
// 否则从GinJwtMiddleware指定的位置寻找
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
	claims, exists := c.Get("jwt-claims")
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

// 将Token加入到Cookie,仅当 SendCookie 开启时有效
func (mw *GinJWTMiddleware) SetCookie(c *gin.Context, token string) {
	if !mw.SendCookie {
		return
	}

	expireCookie := mw.TimeFunc().Add(mw.CookieMaxAge)
	maxage := int(expireCookie.Unix() - mw.TimeFunc().Unix())

	// 设置 Cookie
	c.SetCookie(
		mw.CookieName,
		token,
		maxage,
		"/",
		"",
		mw.SecureCookie,
		mw.SecureCookie,
	)

	if mw.CookieSameSite != 0 {
		c.SetSameSite(mw.CookieSameSite)
	}
}
