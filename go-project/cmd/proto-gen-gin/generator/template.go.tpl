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
	{{$k2}}	{{ $v2.Type }} `json:"{{ $v2.Json }}" form:"{{ $v2.Form }}" uri:"{{$v2.Url}}" header:"{{$v2.Header}}" binding:"{{ $v2.Required }}"`
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
	{{- range $k2,$v2 := $v.SwaggoInfo}}
// @{{$k2}} {{$v2}}
	{{- end}}
	{{- range $k2,$v2 := $v.SwaggoHeaders}}
// @Header {{$v2}}
	{{- end}}
	{{- range $k2,$v2 := $v.SwaggoFailures}}
// @Failure {{$v2}}
	{{- end}}
	{{- range $k2,$v2 := $v.SwaggoParams}}
// @Param {{$v2}}
	{{- end}}
func (s *{{ $.Name }}HttpServer){{ $v.HandlerName }}(c *gin.Context) {
	{{if eq $v.RequestFormName "Empty"}}{{else}}u := &{{$.Name | lower}}form.{{$v.RequestFormName}}{}{{end}}

	{{ $v.QueryStr }}
	{{ $v.HeaderStr }}
	{{ $v.UrlStr }}
	{{ $v.BodyStr }}

	res, err := s.{{$.Name}}Data.{{ $v.HandlerName }}DB(s.Server.Ctx, 	&{{if eq $v.RequestFormName "Empty"}}emptypb.{{$v.RequestType}}{{else}}proto.{{$v.RequestType}}{{end}}{
		{{- range $k2,$v2:=$v.RequestParams}}
		{{$v2.Camel}}: u.{{$v2.Camel}},
		{{- end}}
	})
	if err != nil {
		s.Logger.Sugar().Errorw("微服务调用失败", "msg", err.Error())
		httputil.WriteRpcError(c, err, s.Server.UseAbort)
		return
	}

	httputil.WriteResponse(c, http.StatusOK,"", res, s.Server.UseAbort)
}
{{- end}}

func (s *{{ $.Name }}HttpServer) Execute()error {
{{- range $k,$v := .Methods}}
		s.Server.Engine.{{$v.Method}}("{{$v.Path2Http}}", s.{{ $v.HandlerName }})
{{- end}}
	return s.Server.Serve()
}

//-------------------------------------------------------

// {{$.Name}}DataService是提供{{$.Name}}底层相关数据操作的接口
type {{$.Name}}DataService interface {
{{- range $k,$v := .Methods}}
{{$v.MethodComment }}
	{{ $v.HandlerName }}DB(context.Context, {{if eq $v.RequestType "Empty"}}*emptypb.Empty{{else}}*proto.{{$v.RequestType}}{{end}}) ({{if eq $v.ReplyType "Empty"}}*emptypb.Empty{{else}}*proto.{{$v.ReplyType}}{{end}}, error)
{{- end}}
}

func MustNewGrpc{{$.Name}}Data(c *grpc.ClientConn) {{$.Name}}DataService {
	return &Grpc{{$.Name}}Data{Cli: proto.New{{$.Name}}Client(c)}
}

var _ {{$.Name}}DataService = (*Grpc{{$.Name}}Data)(nil)

// {{$.Name}}服务中的Data层,是数据操作的具体逻辑
type Grpc{{$.Name}}Data struct {
	Cli proto.{{$.Name}}Client
}

{{- range $k,$v := .Methods}}
{{$v.MethodComment }}
func(d *Grpc{{$.Name}}Data){{ $v.HandlerName }}DB(ctx context.Context, in {{if eq $v.RequestType "Empty"}}*emptypb.Empty{{else}}*proto.{{$v.RequestType}}{{end}}) ({{if eq $v.ReplyType "Empty"}}*emptypb.Empty{{else}}*proto.{{$v.ReplyType}}{{end}}, error){
	return d.Cli.{{ $v.HandlerName }}(ctx,in)
}
{{- end}}

type OptionFunc func(*{{$.Name}}HttpServer)

func WithLogger(l *otelzap.Logger) OptionFunc {
	return func(s *{{$.Name}}HttpServer) {
		s.Logger = l
	}
}

func With{{$.Name}}DataService(s {{$.Name | lower}}data.{{$.Name}}DataService) OptionFunc {
	return func(h *{{$.Name}}HttpServer) {
		h.{{$.Name}}Data = s
	}
}
