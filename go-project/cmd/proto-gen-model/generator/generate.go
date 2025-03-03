package generator

import (
	"google.golang.org/protobuf/compiler/protogen"
)

const (
	stringsPkg = protogen.GoImportPath("strings")
	emptyPkg   = protogen.GoImportPath("google.golang.org/protobuf/types/known/emptypb")
)

func GenerateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	if len(file.Services) == 0 {
		return nil
	}

	//设置生成的文件名,文件名会被protoc使用,生成的文件会被放在响应的目录下
	filename := file.GeneratedFilenamePrefix + "_model.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	//该注释会被go的ide识别到, 表示该文件是自动生成的,尽量不要修改
	g.P("package ", file.GoPackageName)

	//该函数是注册全局的packge的内容,但是此时不会写入
	g.QualifiedGoIdent(stringsPkg.Ident(""))

	for _, v1 := range file.Services {
		for _, v2 := range v1.Methods {
			if v2.Input.GoIdent.GoName == "Empty" || v2.Output.GoIdent.GoName == "Empty" {
				g.QualifiedGoIdent(emptyPkg.Ident(""))
				break
			}
		}
	}

	data := ""
	for _, service := range file.Services {
		data += genService(file, g, service)
	}

	g.P(data)
	return g
}

func genService(_ *protogen.File, _ *protogen.GeneratedFile, s *protogen.Service) string {
	// HTTP Server
	sd := &service{
		Name:           s.GoName,
		FullName:       string(s.Desc.FullName()),
		ServiceComment: s.Comments.Leading.String(),
	}

	for _, method := range s.Methods {
		s := genMethod(method)
		// if s.MethodComment[len(s.MethodComment)-1] == '\n' {
		// 	s.MethodComment = s.MethodComment[:len(s.MethodComment)-1]
		// }
		sd.Methods = append(sd.Methods, s)
	}

	return sd.execute()
}

func genMethod(m *protogen.Method) *method {
	return buildMethodDesc(m)
}

func buildMethodDesc(m *protogen.Method) *method {
	md := &method{
		MethodComment: m.Comments.Leading.String(),
		HandlerName:   m.GoName,
		RequestType:   m.Input.GoIdent.GoName,
		ReplyType:     m.Output.GoIdent.GoName,
	}
	return md
}
