package generator

import (
	_ "embed"
	"regexp"
	"strings"
)

//go:embed template.go.tpl
var tpl string

type service struct {
	Name           string
	FullName       string
	ServiceComment string

	//服务的api
	SwaggerApi     map[string]string
	Methods        []*method
	AllRequestForm map[string]map[string]*RequestParam
	AllRequestUsed map[string]int
}

type RequestParam struct {
	//请求体类型内的字段的Camel字符串
	Camel string
	//请求体类型内的字段的Snack字符串
	Snack string
	//请求体类型内的字段的类型
	Type string
	//json名
	Json   string
	Form   string
	Url    string
	Header string

	Xml              string
	Yaml             string
	Required         string
	IsSlice          bool
	SubRequestParams map[string]map[string]*RequestParam
}

// rpc GetDemoName(*Req, *Resp)
type method struct {
	//rpc函数的注解,只收录代码之上的注解
	MethodComment string
	// rpc函数名
	HandlerName string
	//请求体类型
	RequestType string
	//返回体类型
	ReplyType       string
	RequestParams   map[string]*RequestParam
	RequestFormName string
	//记录swagger的注解
	SwaggoInfo map[string]string
	//swagger的Params注解集
	SwaggoParams   []string
	SwaggoHeaders  []string
	SwaggoFailures []string
	//是否是必须的参数
	// 把路由参数转化为http适配的(:id)
	Path2Http string
	//把路由参数转为swagger适配的({id})
	Path2Swagger string
	//路径
	Path string
	//方法类型
	Method string
	//
	Body      string
	HeaderStr string
	UrlStr    string
	QueryStr  string
	FormStr   string
	BodyStr   string
}

// 将所有{xx}的路径参数转为:xx形式的路由参数
func (m *method) pathParams2Http() {
	paths := strings.Split(m.Path2Http, "/")
	for i, p := range paths {
		if p != "" && (p[0] == '{' && p[len(p)-1] == '}' || p[0] == ':') {
			paths[i] = ":" + p[1:len(p)-1]
		}
	}

	m.Path2Http = strings.Join(paths, "/")
}

// 将所有:xx形式的路由参数转换为{xx}形式的路径参数
func (m *method) pathParams2Swagger() {
	paths := strings.Split(m.Path2Swagger, "/")
	for i, p := range paths {
		if len(p) > 1 && p[0] == ':' {
			paths[i] = "{" + p[1:] + "}"
		}
	}

	m.Path2Swagger = strings.Join(paths, "/")
}

// ToExportedCamelCase 转换为 Go 的导出驼峰命名法
func GoExportedCamelCase(s string) string {
	return goCamelCase(s)
}

// ToUnexportedCamelCase 转换为 Go 的不导出驼峰命名法:首字母小写
func GoUnexportedCamelCase(s string) string {
	s = goCamelCase(s)
	return strings.ToLower(s[:1]) + s[1:]
}

func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// 此函数能将str转化为go的驼峰命名格式
func goCamelCase(str string) string {
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
	ss := string(b)
	return strings.ReplaceAll(ss, "_", "")
}

func SnakeCase(s string) string {
	// 使用正则表达式匹配大写字母并在前面加上下划线
	re := regexp.MustCompile(`([a-z0-9])([A-Z])`)
	snake := re.ReplaceAllString(s, `${1}_${2}`)
	// 将所有字符转换为小写
	return strings.ToLower(snake)
}

// HasPathParams 是否包含路由参数
func (m *method) HasPathParams() bool {
	paths := strings.Split(m.Path, "/")
	for _, p := range paths {
		if len(p) > 0 && (p[0] == '{' && p[len(p)-1] == '}' || p[0] == ':') {
			return true
		}
	}

	return false
}
