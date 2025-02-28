//{{$.Name}}服务中的Contoller层,用于对外暴露grpc接口
type {{$.Name}}Server struct {
	Service *{{$.Name | lower}}logic.{{$.Name}}Service
	Logger  *otelzap.Logger
	proto.Unimplemented{{$.Name}}Server
}

{{- range $k, $v := $.Methods}}
{{ $v.MethodComment }}
func (s *{{$.Name}}Server) {{$v.HandlerName}}(ctx context.Context, in {{if eq $v.RequestType "Empty"}}*emptypb.Empty{{else}}*proto.{{$v.RequestType}}{{end}}) ({{if eq $v.ReplyType "Empty"}}*emptypb.Empty{{else}}*proto.{{$v.ReplyType}}{{end}}, error) {
    return nil, errors.New("this method is not implemented")
}
{{- end}}

//---------------------------------------------------------------------

//{{$.Name}}服务中的Service层,编写具体的服务逻辑
type {{$.Name}}Service struct {
	{{$.Name}}Data {{$.Name | lower}}data.{{$.Name}}DataService
}

{{- range $k, $v := $.Methods}}
{{ $v.MethodComment }}
func (s *{{$.Name}}Service) {{$v.HandlerName}}Logic(ctx context.Context, in {{if eq $v.RequestType "Empty"}}*emptypb.Empty{{else}}*proto.{{$v.RequestType}}{{end}}) ({{if eq $v.ReplyType "Empty"}}*emptypb.Empty{{else}}*proto.{{$v.ReplyType}}{{end}}, error) {
    return nil, errors.New("this method is not implemented")
}
{{- end}}

//---------------------------------------------------------------------

//{{$.Name}}DataService是提供{{$.Name}}底层相关数据操作的接口
type UserDataService interface{
{{- range $k, $v := $.Methods}}
{{ $v.MethodComment }}
	{{$v.HandlerName}}DB(ctx context.Context, in {{if eq $v.RequestType "Empty"}}*emptypb.Empty{{else}}*proto.{{$v.RequestType}}{{end}}) ({{if eq $v.ReplyType "Empty"}}*emptypb.Empty{{else}}*proto.{{$v.ReplyType}}{{end}}, error)
{{- end}}
}

//{{$.Name}}服务中的Data层,是数据操作的具体逻辑
type {{$.Name}}Data struct {}

{{- range $k, $v := $.Methods}}
{{ $v.MethodComment }}
func (s *{{$.Name}}Data) {{$v.HandlerName}}DB(ctx context.Context, in {{if eq $v.RequestType "Empty"}}*emptypb.Empty{{else}}*proto.{{$v.RequestType}}{{end}}) ({{if eq $v.ReplyType "Empty"}}*emptypb.Empty{{else}}*proto.{{$v.ReplyType}}{{end}}, error) {
    return nil, errors.New("this method is not implemented")
}
{{- end}}

