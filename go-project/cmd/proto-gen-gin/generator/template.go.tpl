// PathsMap中key为对应的handler函数名,{{$.ServiceComment }}
{{- range $k,$v := $.SwaggerApi}}
// @{{$k}} {{$v}}
{{- end}}
type {{ $.Name }}HttpServer struct{
	PathsMap	map[string] gin.HandlerFunc
}

func Register{{ $.Name }}HTTPServer() *{{ $.Name }}HttpServer {
	s := &{{ $.Name }}HttpServer{
		PathsMap:  make(map[string] gin.HandlerFunc),
	}
	return s
}


func (s *{{ $.Name }}HttpServer)RegisterHandlerFunc(method string,path string,handler gin.HandlerFunc) {
	s.PathsMap[method+";"+path] = handler
}

{{- range $k,$v := .Methods}}
{{ $v.MethodComment }}
	{{- range $k2,$v2 := $v.SwaggerApi}}
// @{{$k2}} {{$v2}}
	{{- end}}
func (s *{{ $.Name }}HttpServer){{ $v.HandlerName }}(c *gin.Context) {}
{{- end}}

func (s *{{ $.Name }}HttpServer) register() {
{{- range .Methods}}
		s.RegisterHandlerFunc("{{.Method}}","{{.Path}}", s.{{ .HandlerName }})
{{- end}}
}

func(s *{{ $.Name }}HttpServer) Execute(r gin.IRouter){
	s.register()
	for k,v := range s.PathsMap{
		path := strings.Split(k,";")
		path[0] = strings.ToUpper(path[0])
		switch path[0]{
			case "GET":
				r.GET(path[1],v)
			case "PUT":
				r.PUT(path[1],v)
			case "DELETE":
				r.DELETE(path[1],v)
			case "POST":
				r.POST(path[1],v)
			case "PATCH":
				r.PATCH(path[1],v)
		}
	}
}