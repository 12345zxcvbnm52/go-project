package token

import (
	"encoding/base64"
	"net/http"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
)

type CookieClaims struct {
	ID       uint32 `json:"id,omitempty"`
	UserName string `json:"username,omitempty"`
	LifeTime int    `json:"lftime"`
}

func (cookie *CookieClaims) GenCookie(c *gin.Context) string {
	jstr, _ := sonic.Marshal(cookie)
	session_id := base64.StdEncoding.EncodeToString([]byte(jstr))

	c.SetCookie("x-cookie-"+strconv.Itoa(int(cookie.ID)), session_id, cookie.LifeTime, c.Request.URL.String(), "localhost", false, true)
	return session_id
}

func (cookie *CookieClaims) VerifyCookie(c *gin.Context) {
	jstr, _ := sonic.Marshal(cookie)
	session_id := base64.StdEncoding.EncodeToString([]byte(jstr))
	var exist = false
	for _, v := range c.Request.Cookies() {
		if v.Name == "x-cookie-"+strconv.Itoa(int(cookie.ID)) && v.Value == session_id {
			exist = true
			break
		}
	}
	if !exist {
		c.Redirect(http.StatusMovedPermanently, "/login")
		c.Abort()
		return
		//可以考虑再次进入Verify函数
	}

}
