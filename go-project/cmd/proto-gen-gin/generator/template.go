package generator

import (
	"bytes"
	_ "embed"
	"html/template"
	"strings"
)

//go:embed template.go.tpl
var tpl string

// rpc GetDemoName(*Req, *Resp)
type method struct {
	MethodComment string
	HandlerName   string
	RequestType   string
	ReplyType     string
	SwaggerApi    map[string]string
	// 路由
	Path string
	//路由参数
	PathParams []string
	//方法类型
	Method       string
	Body         string
	ResponseBody string
}

// HasPathParams 检测路由是否存在路径参数,例如/user/{id}或/product/:name
func (m *method) HasPathParams() bool {
	paths := strings.Split(m.Path, "/")
	for _, p := range paths {
		if len(p) > 0 && (p[0] == '{' && p[len(p)-1] == '}' || p[0] == ':') {
			return true
		}
	}

	return false
}

// 将所有{xx}的路径参数转为:xx形式的路由参数
func (m *method) initPathParams() {
	paths := strings.Split(m.Path, "/")
	for i, p := range paths {
		if p != "" && (p[0] == '{' && p[len(p)-1] == '}' || p[0] == ':') {
			paths[i] = ":" + p[1:len(p)-1]
			m.PathParams = append(m.PathParams, paths[i][1:])
		}
	}

	m.Path = strings.Join(paths, "/")
}

type service struct {
	Name           string
	FullName       string
	ServiceComment string

	SwaggerApi map[string]string
	Methods    []*method
}

func (s *service) execute() string {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("http").Parse(strings.TrimSpace(tpl))
	if err != nil {
		panic(err)
	}

	if err := tmpl.Execute(buf, s); err != nil {
		panic(err)
	}

	return buf.String()
}

func (s *service) ServiceName() string {
	return s.Name + "Server"
}

func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// 此函数能将str转化为go的驼峰命名格式
func (s *service) GoCamelCase(str string) string {
	var b []byte
	for i := 0; i < len(str); i++ {
		c := str[i]
		switch {
		case c == '.' && i+1 < len(str) && isASCIILower(str[i+1]):
		case c == '.':
			b = append(b, '_')
		case c == '_' && (i == 0 || str[i-1] == '.'):

			b = append(b, 'X')
		case c == '_' && i+1 < len(str) && isASCIILower(str[i+1]):
		case isASCIIDigit(c):
			b = append(b, c)
		default:
			if isASCIILower(c) {
				c -= 'a' - 'A'
			}
			b = append(b, c)
			for ; i+1 < len(str) && isASCIILower(str[i+1]); i++ {
				b = append(b, str[i+1])
			}
		}
	}
	return string(b)
}
