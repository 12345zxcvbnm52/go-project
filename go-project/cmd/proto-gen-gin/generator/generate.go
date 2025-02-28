package generator

import (
	"fmt"
	reflect "reflect"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

const (
	ginPkg     = protogen.GoImportPath("github.com/gin-gonic/gin")
	stringsPkg = protogen.GoImportPath("strings")
	smptyPkg   = protogen.GoImportPath("google.golang.org/protobuf/types/known/emptypb")
)

func GenerateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	if len(file.Services) == 0 {
		return nil
	}

	//设置生成的文件名,文件名会被protoc使用,生成的文件会被放在响应的目录下
	filename := file.GeneratedFilenamePrefix + "_gin.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	//该注释会被go的ide识别到, 表示该文件是自动生成的,尽量不要修改
	g.P("package ", file.GoPackageName)

	//该函数是注册全局的packge的内容,但是此时不会写入
	g.QualifiedGoIdent(ginPkg.Ident(""))
	g.QualifiedGoIdent(stringsPkg.Ident(""))
	data := ""
	for _, service := range file.Services {
		data += genService(file, g, service)
	}

	g.P(data)
	return g
}

func genService(_ *protogen.File, _ *protogen.GeneratedFile, s *protogen.Service) string {
	sd := &service{
		Name:           s.GoName,
		FullName:       string(s.Desc.FullName()),
		ServiceComment: s.Comments.Leading.String(),
		SwaggerApi:     make(map[string]string),
	}

	if sd.ServiceComment[len(sd.ServiceComment)-1] == '\n' {
		sd.ServiceComment = sd.ServiceComment[:len(sd.ServiceComment)-1]
	}

	serviceRule, ok := proto.GetExtension(s.Desc.Options(), E_ServiceOpt).(*ServiceOption)
	if serviceRule != nil && ok {
		val := reflect.ValueOf(serviceRule)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		tp := val.Type()
		for i := 0; i < val.NumField(); i++ {
			fieldType := tp.Field(i)
			if fieldType.PkgPath != "" || !val.Field(i).CanSet() {
				continue
			}
			fieldValue := val.Field(i).Interface().(string)
			if fieldValue == "" {
				continue
			}
			fieldName := val.Type().Field(i).Name
			sd.SwaggerApi[fieldName] = fieldValue
		}
	}

	for _, method := range s.Methods {
		m := genMethod(method)
		if m.MethodComment[len(m.MethodComment)-1] == '\n' {
			m.MethodComment = m.MethodComment[:len(m.MethodComment)-1]
		}
		sd.Methods = append(sd.Methods, m)
	}

	// for k, v := range sd.Methods[0].SwaggerApi {
	// 	fmt.Println(k, v)
	// }

	return sd.execute()
}

func genMethod(m *protogen.Method) *method {
	var methods *method

	rule, ok := proto.GetExtension(m.Desc.Options(), annotations.E_Http).(*annotations.HttpRule)
	if rule != nil && ok {
		methods = buildHTTPRule(m, rule)
	} else {
		methods = defaultMethod(m)
	}

	if methods.Method+methods.Path != "" {
		methods.SwaggerApi["Router"] = fmt.Sprintf("%s [%s]", methods.Path, methods.Method)
	}

	methodRule, ok := proto.GetExtension(m.Desc.Options(), E_MethodOpt).(*MethodOption)
	if methodRule != nil && ok {
		val := reflect.ValueOf(methodRule)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		tp := val.Type()

		// 遍历结构体的所有字段
		for i := 0; i < val.NumField(); i++ {
			fieldType := tp.Field(i)
			if fieldType.PkgPath != "" || !val.Field(i).CanSet() {
				continue
			}
			fieldValue := val.Field(i).Interface().(string)
			if fieldValue == "" {
				continue
			}
			fieldName := val.Type().Field(i).Name

			methods.SwaggerApi[fieldName] = fieldValue
		}
	}
	return methods
}

func defaultMethod(m *protogen.Method) *method {
	// TODO path
	// $prefix + / + ${package}.${service} + / + ${method}
	// /api/demo.v0.Demo/GetName
	md := buildMethodDesc(m, "GET", "/"+"default/"+m.GoName)
	md.Body = "*"
	return md
}

func buildHTTPRule(m *protogen.Method, rule *annotations.HttpRule) *method {
	var path, method string
	switch pattern := rule.Pattern.(type) {
	case *annotations.HttpRule_Get:
		path = pattern.Get
		method = "GET"
	case *annotations.HttpRule_Put:
		path = pattern.Put
		method = "PUT"
	case *annotations.HttpRule_Post:
		path = pattern.Post
		method = "POST"
	case *annotations.HttpRule_Delete:
		path = pattern.Delete
		method = "DELETE"
	case *annotations.HttpRule_Patch:
		path = pattern.Patch
		method = "PATCH"
	case *annotations.HttpRule_Custom:
		path = pattern.Custom.Path
		method = pattern.Custom.Kind
	}

	md := buildMethodDesc(m, method, path)
	return md
}

func buildMethodDesc(m *protogen.Method, httpMethod string, path string) *method {
	md := &method{
		MethodComment: m.Comments.Leading.String(),
		HandlerName:   m.GoName,
		RequestType:   m.Input.GoIdent.GoName,
		ReplyType:     m.Output.GoIdent.GoName,
		Path:          path,
		Method:        httpMethod,
		SwaggerApi:    make(map[string]string),
	}

	md.initPathParams()
	return md
}
