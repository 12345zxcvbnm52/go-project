package secure

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	//强制浏览器只能通过HTTPS连接到网站,includeSubdomains意味着这个策略也适用于当前域名的所有子域,
	stsHeader          = "Strict-Transport-Security"
	stsSubdomainString = "; includeSubdomains"

	//防止网站被嵌入到<iframe>中,避免点击劫持攻击,确保页面不能被嵌套在其他站点的框架中,
	//DENY 表示完全不允许嵌套,SAMEORIGIN 表示只有同源的页面才允许嵌套,
	frameOptionsHeader = "X-Frame-Options"
	frameOptionsValue  = "DENY"
	//防止浏览器MIME类型嗅探,nosniff可以防止浏览器基于文件内容来猜测MIME类型，这有助于防止某些类型的攻击如执行恶意JavaScript
	contentTypeHeader = "X-Content-Type-Options"
	contentTypeValue  = "nosniff"

	xssProtectionHeader = "X-XSS-Protection"
	xssProtectionValue  = "1; mode=block"
	//允许指定内容如脚本,样式表可以加载,从而减少XSS和其他攻击的风险,
	//通过设置Content-Security-Policy可以限制哪些资源能被加载或者仅允许来自特定来源的资源
	cspHeader = "Content-Security-Policy"
)

func defaultBadHostHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Bad Host", http.StatusInternalServerError)
}

// Options is a struct for specifying configuration options for the secure.Secure middleware.
type Options struct {
	// AllowedHosts is a list of fully qualified domain names that are allowed. Default is empty list, which allows any and all host names.
	AllowedHosts []string
	// If SSLRedirect is set to true, then only allow https requests. Default is false.
	SSLRedirect bool
	// If SSLTemporaryRedirect is true, the a 302 will be used while redirecting. Default is false (301).
	SSLTemporaryRedirect bool
	// SSLHost is the host name that is used to redirect http requests to https. Default is "", which indicates to use the same host.
	SSLHost string
	// SSLProxyHeaders is set of header keys with associated values that would indicate a valid https request. Useful when using Nginx: `map[string]string{"X-Forwarded-Proto": "https"}`. Default is blank map.
	SSLProxyHeaders map[string]string
	// STSSeconds is the max-age of the Strict-Transport-Security header. Default is 0, which would NOT include the header.
	STSSeconds int64
	// If STSIncludeSubdomains is set to true, the `includeSubdomains` will be appended to the Strict-Transport-Security header. Default is false.
	STSIncludeSubdomains bool
	// If FrameDeny is set to true, adds the X-Frame-Options header with the value of `DENY`. Default is false.
	FrameDeny bool
	// CustomFrameOptionsValue allows the X-Frame-Options header value to be set with a custom value. This overrides the FrameDeny option.
	CustomFrameOptionsValue string
	// If ContentTypeNosniff is true, adds the X-Content-Type-Options header with the value `nosniff`. Default is false.
	ContentTypeNosniff bool
	// If BrowserXssFilter is true, adds the X-XSS-Protection header with the value `1; mode=block`. Default is false.
	BrowserXssFilter bool
	// ContentSecurityPolicy allows the Content-Security-Policy header value to be set with a custom value. Default is "".
	ContentSecurityPolicy string
	// When developing, the AllowedHosts, SSL, and STS options can cause some unwanted effects. Usually testing happens on http, not https, and on localhost, not your production domain... so set this to true for dev environment.
	// If you would like your development environment to mimic production with complete Host blocking, SSL redirects, and STS headers, leave this as false. Default if false.
	IsDevelopment bool

	// Handlers for when an error occurs (ie bad host).
	BadHostHandler http.Handler
}

// Secure is a middleware that helps setup a few basic security features. A single secure.Options struct can be
// provided to configure which features should be enabled, and the ability to override a few of the default values.
type secure struct {
	// Customize Secure with an Options struct.
	opt Options
}

// Constructs a new Secure instance with supplied options.
func New(options Options) *secure {
	if options.BadHostHandler == nil {
		options.BadHostHandler = http.HandlerFunc(defaultBadHostHandler)
	}

	return &secure{
		opt: options,
	}
}

func (s *secure) process(w http.ResponseWriter, r *http.Request) error {
	// Allowed hosts check.
	if len(s.opt.AllowedHosts) > 0 && !s.opt.IsDevelopment {
		isGoodHost := false
		for _, allowedHost := range s.opt.AllowedHosts {
			if strings.EqualFold(allowedHost, r.Host) {
				isGoodHost = true
				break
			}
		}

		if !isGoodHost {
			s.opt.BadHostHandler.ServeHTTP(w, r)
			return fmt.Errorf("Bad host name: %s", r.Host)
		}
	}

	// SSL check.
	if s.opt.SSLRedirect && s.opt.IsDevelopment == false {
		isSSL := false
		if strings.EqualFold(r.URL.Scheme, "https") || r.TLS != nil {
			isSSL = true
		} else {
			for k, v := range s.opt.SSLProxyHeaders {
				if r.Header.Get(k) == v {
					isSSL = true
					break
				}
			}
		}

		if isSSL == false {
			url := r.URL
			url.Scheme = "https"
			url.Host = r.Host

			if len(s.opt.SSLHost) > 0 {
				url.Host = s.opt.SSLHost
			}

			status := http.StatusMovedPermanently
			if s.opt.SSLTemporaryRedirect {
				status = http.StatusTemporaryRedirect
			}

			http.Redirect(w, r, url.String(), status)
			return fmt.Errorf("Redirecting to HTTPS")
		}
	}

	// Strict Transport Security header.
	if s.opt.STSSeconds != 0 && !s.opt.IsDevelopment {
		stsSub := ""
		if s.opt.STSIncludeSubdomains {
			stsSub = stsSubdomainString
		}

		w.Header().Add(stsHeader, fmt.Sprintf("max-age=%d%s", s.opt.STSSeconds, stsSub))
	}

	// Frame Options header.
	if len(s.opt.CustomFrameOptionsValue) > 0 {
		w.Header().Add(frameOptionsHeader, s.opt.CustomFrameOptionsValue)
	} else if s.opt.FrameDeny {
		w.Header().Add(frameOptionsHeader, frameOptionsValue)
	}

	// Content Type Options header.
	if s.opt.ContentTypeNosniff {
		w.Header().Add(contentTypeHeader, contentTypeValue)
	}

	// XSS Protection header.
	if s.opt.BrowserXssFilter {
		w.Header().Add(xssProtectionHeader, xssProtectionValue)
	}

	// Content Security Policy header.
	if len(s.opt.ContentSecurityPolicy) > 0 {
		w.Header().Add(cspHeader, s.opt.ContentSecurityPolicy)
	}

	return nil

}

func Secure(options Options) gin.HandlerFunc {
	s := New(options)

	return func(c *gin.Context) {
		err := s.process(c.Writer, c.Request)
		if err != nil {
			if c.Writer.Written() {
				c.AbortWithStatus(c.Writer.Status())
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		}
	}

}
