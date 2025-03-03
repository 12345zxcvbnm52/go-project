{{- $.ServiceComment -}}
{{- range $k,$v := $.SwaggerApi}}
// @{{$k}} {{$v}}
{{- end}}
type {{ $.Name }}HttpServer struct{
	Server 		*httpserver.Server
	{{$.Name}}Data {{$.Name | lower}}data.{{$.Name}}DataService
	Logger   *otelzap.Logger
}

{{ range $k1,$v1:= $.AllRequestForm}}
type {{$k1}} struct{
{{- range $k2,$v2 := $v1}}
	{{$k2}}	{{ $v2.Type }} `json:"{{ $v2.Snack }}" form:"{{ $v2.Snack }}" binding:"{{ $v2.Required }}"`
{{- end}}
}
{{ end}}

// 默认使用otelzap.Logger以及Grpc{{$.Name}}Data
func MustNew{{ $.Name }}HTTPServer(s *httpserver.Server,opts ...OptionFunc) *{{ $.Name }}HttpServer {
	ss := &{{ $.Name }}HttpServer{
		Server:		s,
	}
	for _, opt := range opts {
		opt(ss)
	}
	if ss.Logger == nil {
		ss.Logger = log.MustNewOtelLogger()
	}
	if ss.{{$.Name}}Data == nil {
		cli, err := s.GrpcCli.Dial()
		if err != nil {
			panic(err)
		}
		ss.{{ $.Name }}Data = {{ $.Name | lower}}data.MustNewGrpc{{ $.Name }}Data(cli)
	}
	return ss
}

{{- range $k,$v := .Methods}}
{{ $v.MethodComment }}
	{{- range $k2,$v2 := $v.SwaggerApi}}
// @{{$k2}} {{$v2}}
	{{- end}}
	{{- range $k2,$v2 := $v.Params}}
// @Param {{$v2}}
	{{- end}}
func (s *{{ $.Name }}HttpServer){{ $v.HandlerName }}(c *gin.Context) {
	u := &form.{{$v.RequestFormName}}{}
	if err := c.ShouldBind(u); err != nil {
		s.ValidatorErrorHandle(c, err)
		return
	}

	res, err := s.{{$.Name}}Data.{{ $v.HandlerName }}DB(s.Server.Ctx, &proto.{{$v.RequestType}}{
		{{- range $k2,$v2:=$v.RequestParams}}
		{{$v2.Camel}}: u.{{$v2.Camel}},
		{{- end}}
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		RpcErrorHandle(c, err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
{{- end}}

func (s *{{ $.Name }}HttpServer) Execute()error {
{{- range $k,$v := .Methods}}
		s.Server.Engine.{{$v.Method}}("{{$v.Path2Http}}", s.{{ $v.HandlerName }})
{{- end}}
	return s.Server.Serve()
}

type OptionFunc func(*UserHttpServer)

func WithLogger(l *otelzap.Logger) OptionFunc {
	return func(s *UserHttpServer) {
		s.Logger = l
	}
}

func WithUserDataService(s userdata.UserDataService) OptionFunc {
	return func(h *UserHttpServer) {
		h.UserData = s
	}
}
