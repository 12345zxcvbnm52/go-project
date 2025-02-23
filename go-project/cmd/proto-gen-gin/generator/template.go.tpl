// PathsMap中key为对应的handler函数名
{{ $.ServiceComment }}
type {{ $.Name }}HttpServer struct{
	PathsMap	map[string] gin.HandlerFunc
	Router  gin.IRouter
}

func Register{{ $.Name }}HTTPServer(r gin.IRouter) *{{ $.Name }}HttpServer {
	s := &{{ $.Name }}HttpServer{
		Router:     r,
		PathsMap:  make(map[string] gin.HandlerFunc),
	}
	return s
}

func (s *{{ $.Name }}HttpServer)RegisterHandlerFunc(method string,path string,handler gin.HandlerFunc) {
	s.PathsMap[method+";"+path] = handler
}

{{range $k,$v := .Methods}}
{{ $v.MethodComment }}
func (s *{{ $.Name }}HttpServer){{ $v.HandlerName }}(c *gin.Context) {}
{{end}}

func (s *{{ $.Name }}HttpServer) register() {
{{range .Methods}}
		s.RegisterHandlerFunc("{{.Method}}","{{.Path}}", s.{{ .HandlerName }})
{{end}}
}

func(s *{{ $.Name }}HttpServer) Excute(){
	s.register()
	for k,v := range s.PathsMap{
		path := strings.Split(k,";")
		path[0] = strings.ToUpper(path[0])
		switch path[0]{
			case "GET":
				s.Router.GET(path[1],v)
			case "PUT":
				s.Router.PUT(path[1],v)
			case "DELETE":
				s.Router.DELETE(path[1],v)
			case "POST":
				s.Router.POST(path[1],v)
			case "PATCH":
				s.Router.PATCH(path[1],v)
		}
	}
}