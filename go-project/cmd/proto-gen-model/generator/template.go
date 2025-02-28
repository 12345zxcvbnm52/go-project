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
}

type service struct {
	Name           string
	FullName       string
	ServiceComment string

	Methods []*method
}

func (s *service) execute() string {
	var funcMap = template.FuncMap{
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
	}

	buf := new(bytes.Buffer)
	tmpl, err := template.New("text").Funcs(funcMap).Parse(strings.TrimSpace(tpl))
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
