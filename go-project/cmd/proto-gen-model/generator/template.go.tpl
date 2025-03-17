//{{$.Name}}服务中的Contoller层,用于对外暴露grpc接口
type {{$.Name}}Server struct {
	Service *{{$.Name | lower}}logic.{{$.Name}}Service
	Logger  *otelzap.Logger
	proto.Unimplemented{{$.Name}}Server
}

var ProtoJson = protojson.MarshalOptions{
	EmitUnpopulated: true,
}

func MethodInfoRecord(data gproto.Message) string {
	r, err := ProtoJson.Marshal(data)
	if err != nil {
		return ""
	}
	return string(r)
}

//---------------------------------------------------------------------

{{- range $k, $v := $.Methods}}
{{ $v.MethodComment }}
func (s *{{$.Name}}Server) {{$v.HandlerName}}(ctx context.Context, in {{if eq $v.RequestType "Empty"}}*emptypb.Empty{{else}}*proto.{{$v.RequestType}}{{end}}) ({{if eq $v.ReplyType "Empty"}}*emptypb.Empty{{else}}*proto.{{$v.ReplyType}}{{end}}, error) {
    info := MethodInfoRecord(in)
	s.Logger.Sugar().Infof("正在进行一次{{$v.HandlerName}}调用,调用信息为: %s", info)
	res, err := s.Service.{{$v.HandlerName}}Logic(ctx, in)
	if err != nil {
		s.Logger.Sugar().Errorf("调用{{$v.HandlerName}}失败,具体信息为: %s", err.Error())
		return nil, err
	}
	return res, nil
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
    return s.{{$.Name}}Data.{{$v.HandlerName}}DB(ctx, in)
}
{{- end}}

//---------------------------------------------------------------------

//{{$.Name}}DataService是提供{{$.Name}}底层相关数据操作的接口
type {{$.Name}}DataService interface{
{{- range $k, $v := $.Methods}}
{{ $v.MethodComment }}
	{{$v.HandlerName}}DB(ctx context.Context, in {{if eq $v.RequestType "Empty"}}*emptypb.Empty{{else}}*proto.{{$v.RequestType}}{{end}}) ({{if eq $v.ReplyType "Empty"}}*emptypb.Empty{{else}}*proto.{{$v.ReplyType}}{{end}}, error)
{{- end}}
}

//{{$.Name}}服务中的Data层,是数据操作的具体逻辑
type Gorm{{$.Name}}Data struct {
	DB *gorm.DB
}

var _ {{$.Name}}DataService = (*Gorm{{$.Name}}Data)(nil)

func MustNewGorm{{$.Name}}Data(db *gorm.DB) *Gorm{{$.Name}}Data {
	return &Gorm{{$.Name}}Data{DB: db}
}


{{- range $k, $v := $.Methods}}
{{ $v.MethodComment }}
func (s *Gorm{{$.Name}}Data) {{$v.HandlerName}}DB(ctx context.Context, in {{if eq $v.RequestType "Empty"}}*emptypb.Empty{{else}}*proto.{{$v.RequestType}}{{end}}) ({{if eq $v.ReplyType "Empty"}}*emptypb.Empty{{else}}*proto.{{$v.ReplyType}}{{end}}, error) {
    return nil, errors.New("this method is not implemented")
}
{{- end}}

