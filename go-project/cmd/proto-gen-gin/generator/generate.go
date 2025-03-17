package generator

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	reflect "reflect"
	"strings"
	"text/template"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

const (
	queryStr = `	if err := c.ShouldBindQuery(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}`
	urlStr = `	if err := c.ShouldBindUri(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}`
	headerStr = `	if err := c.ShouldBindHeader(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}`
	bodyStr = `	if err := c.ShouldBind(u); err != nil {
		httputil.WriteValidateError(c, s.Server.Validator.Trans, err, s.Server.UseAbort)
		return
	}`
)

const (
	ginPkg = protogen.GoImportPath("github.com/gin-gonic/gin")
	//stringsPkg = protogen.GoImportPath("strings")
	smptyPkg   = protogen.GoImportPath("google.golang.org/protobuf/types/known/emptypb")
	httpsrvPkg = protogen.GoImportPath("kenshop/goken/server/httpserver")
	logPkg     = protogen.GoImportPath("kenshop/pkg/log")
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
	//g.QualifiedGoIdent(stringsPkg.Ident(""))
	g.QualifiedGoIdent(httpsrvPkg.Ident(""))
	g.QualifiedGoIdent(logPkg.Ident(""))

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
		AllRequestForm: make(map[string]map[string]*RequestParam),
		AllRequestUsed: make(map[string]int),
	}

	// if sd.ServiceComment != "" && sd.ServiceComment[len(sd.ServiceComment)-1] == '\n' {
	// 	sd.ServiceComment = sd.ServiceComment[:len(sd.ServiceComment)-1]
	// }

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
		if m == nil {
			continue
		}
		if m.MethodComment != "" && m.MethodComment[len(m.MethodComment)-1] == '\n' {
			m.MethodComment = m.MethodComment[:len(m.MethodComment)-1]
		}
		if m.RequestFormName == "Empty" {
			sd.Methods = append(sd.Methods, m)
			continue
		}

		for _, v := range m.RequestParams {
			if v.SubRequestParams != nil {
				for reqName, subInfo := range v.SubRequestParams {
					subFormName := reqTurnFormName(reqName)
					if _, ok := sd.AllRequestForm[subFormName]; !ok {
						sd.AllRequestForm[subFormName] = subInfo
					}
				}

			}
		}
		//如果还没有这个Form类型就记录它
		if v, ok := sd.AllRequestForm[m.RequestFormName]; !ok {
			sd.AllRequestForm[m.RequestFormName] = m.RequestParams
			//如果有了这个Form类型,把该form的使用次数加一,并且给该form的命名加上后缀_i,i是第几次出现,以便分别
		} else {
			sd.AllRequestUsed[m.RequestFormName]++
			i := sd.AllRequestUsed[m.RequestFormName]
			r := fmt.Sprintf("%s_%d", m.RequestFormName, i)
			if i == 1 {
				//记得i==1时要把原初的form加上_0
				sd.AllRequestForm[fmt.Sprintf("%s_%d", m.RequestFormName, 0)] = v
				for _, v := range sd.Methods {
					if v.RequestFormName == m.RequestFormName {
						v.RequestFormName = fmt.Sprintf("%s_%d", m.RequestFormName, 0)
					}
				}
			}
			sd.AllRequestForm[r] = m.RequestParams
			//最后要删除原始form
			delete(sd.AllRequestForm, m.RequestFormName)
			m.RequestFormName = fmt.Sprintf("%s_%d", m.RequestFormName, i)
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

	//不导出的方法过滤
	unimport := m.GoName[:1]
	if unimport == strings.ToLower(unimport) {
		return nil
	}

	rule, ok := proto.GetExtension(m.Desc.Options(), annotations.E_Http).(*annotations.HttpRule)
	methods = buildHTTPRule(m, rule, ok)
	buildForm(m, methods)
	s := ""
	for k, v := range methods.RequestParams {
		s += fmt.Sprintf("%s:%+v\n", k, *v)
	}
	os.Stderr.Write([]byte(s))
	buildSwagger(m, methods)

	return methods
}

func buildHTTPRule(m *protogen.Method, rule *annotations.HttpRule, conv bool) *method {
	var path, method string
	if rule != nil && conv {
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
	} else {
		method = "GET"
		path = "/default/" + m.GoName
	}

	md := buildMethodDesc(m, method, path)
	return md
}

func reqTurnFormName(s string) string {
	if strings.HasSuffix(s, "Req") {
		return strings.TrimSuffix(s, "Req") + "Form"
	}
	return s
}

func buildMethodDesc(m *protogen.Method, httpMethod string, path string) *method {

	md := &method{
		MethodComment:   m.Comments.Leading.String(),
		HandlerName:     m.GoName,
		RequestType:     m.Input.GoIdent.GoName,
		ReplyType:       m.Output.GoIdent.GoName,
		RequestFormName: reqTurnFormName(m.Input.GoIdent.GoName),
		Path:            path,
		Path2Http:       path,
		Path2Swagger:    path,
		Method:          httpMethod,
		SwaggoInfo:      make(map[string]string),
		RequestParams:   make(map[string]*RequestParam),
		SwaggoParams:    make([]string, 0),
		SwaggoHeaders:   make([]string, 0),
		SwaggoFailures:  make([]string, 0),
	}
	md.pathParams2Http()
	md.pathParams2Swagger()
	return md
}

func buildSwagger(m *protogen.Method, methods *method) {
	//先把Router定义下来
	if methods.Method+methods.Path != "" {
		methods.SwaggoInfo["Router"] = fmt.Sprintf("%s [%s]", methods.Path2Swagger, methods.Method)
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
			//排除不可导出字段可不可赋值字段
			if fieldType.PkgPath != "" || !val.Field(i).CanSet() {
				continue
			}

			//如果是普通的swagger字段(记录一条string)就记录
			fieldValue, ok := val.Field(i).Interface().(string)
			if ok {
				//如果swagger注解为空就不记录
				if fieldValue == "" {
					continue
				}
				fieldName := val.Type().Field(i).Name

				methods.SwaggoInfo[fieldName] = fieldValue
				//如果是swagger的Params字段(记录多个string,即[]string)
			} else {
				fieldValue, ok := val.Field(i).Interface().([]string)
				if ok {
					fieldName := val.Type().Field(i).Name
					switch fieldName {
					case "Params":
						//把'号替换为"
						f := func(inputs ...string) []string {
							res := make([]string, 0)
							for _, v := range inputs {
								res = append(res, strings.ReplaceAll(v, "'", "\""))
							}
							return res
						}
						methods.SwaggoParams = append(methods.SwaggoParams, f(fieldValue...)...)
					case "Headers":
						methods.SwaggoHeaders = append(methods.SwaggoHeaders, fieldValue...)
					case "Failures":
						methods.SwaggoFailures = append(methods.SwaggoFailures, fieldValue...)
					}
				}
			}
		}
		for _, v := range methods.SwaggoParams {
			words := strings.Fields(v)
			im, ok := methods.RequestParams[GoExportedCamelCase(words[0])]
			if !ok {
				methods.RequestParams[GoExportedCamelCase(words[0])] = &RequestParam{}
				im = methods.RequestParams[GoExportedCamelCase(words[0])]
			}
			//查看Params的第四个参数是否是必须的
			switch words[3] {
			case "true":
				if words[1] == "body" {
					im.Required = "required"
					for _, v := range im.SubRequestParams {
						for _, v2 := range v {
							v2.Required = "required"
						}
					}
				}
			default:
				im.Required = ""
			}

			switch words[1] {
			case "path":
				im.Url = SnakeCase(words[0])
				for _, v := range im.SubRequestParams {
					for _, v2 := range v {
						v2.Url = v2.Snack
					}
				}
				methods.UrlStr = urlStr
			case "header":
				im.Header = SnakeCase(words[0])
				for _, v := range im.SubRequestParams {
					for _, v2 := range v {
						v2.Header = v2.Snack
					}
				}
				methods.HeaderStr = headerStr
			case "query":
				methods.QueryStr = queryStr
			case "body":
				methods.BodyStr = bodyStr
			case "formData":
				methods.BodyStr = bodyStr
			default:
			}

			im.Json = SnakeCase(words[0])
			for _, v := range im.SubRequestParams {
				for _, v2 := range v {
					v2.Json = v2.Snack
				}
			}
			im.Form = SnakeCase(words[0])
			for _, v := range im.SubRequestParams {
				for _, v2 := range v {
					v2.Form = v2.Snack
				}
			}
		}
		// if v, ok := methods.SwaggoInfo["Description"]; !ok || v == "" {
		// 	methods.SwaggoInfo["Description"] = methods.MethodComment
		// }
	}
}

func buildForm(m *protogen.Method, methods *method) {
	for _, v := range m.Input.Fields {
		im, ok := methods.RequestParams[GoExportedCamelCase(v.GoName)]
		if !ok {
			methods.RequestParams[GoExportedCamelCase(v.GoName)] = &RequestParam{}
			im = methods.RequestParams[GoExportedCamelCase(v.GoName)]
		}
		if v.Oneof == nil && v.GoName != "" {
			im.Camel = GoExportedCamelCase(v.GoName)
			im.Snack = SnakeCase(v.GoName)
			if v.Desc.Kind().String() == "message" {
				im.SubRequestParams = make(map[string]map[string]*RequestParam)
				im.SubRequestParams[GoExportedCamelCase(v.GoName)] = make(map[string]*RequestParam)
				for _, iv := range v.Message.Fields {
					if iv.Oneof == nil && iv.GoName != "" {
						if iv.Desc.Kind().String() == "message" {
							panic(errors.New("protoc-gen-gin not allow three message enbed"))
						}
						im.SubRequestParams[GoExportedCamelCase(v.GoName)][GoExportedCamelCase(iv.GoName)] = &RequestParam{}
						cim := im.SubRequestParams[GoExportedCamelCase(v.GoName)][GoExportedCamelCase(iv.GoName)]
						cim.Camel = GoExportedCamelCase(iv.GoName)
						cim.Snack = SnakeCase(iv.GoName)
						ctp := iv.Desc.Kind().String()
						if ctp == "float" {
							ctp = "float32"
						} else if ctp == "double" {
							ctp = "float64"
						} else if iv.Desc.IsList() {
							ctp = "[]*" + ctp
							cim.IsSlice = true
						}
						cim.Type = ctp
					}
				}
				if v.Desc.IsList() {
					im.Type = "[]" + reqTurnFormName(v.GoName)
					if v.Desc.IsList() {
						im.Type = "[]*" + reqTurnFormName(v.GoName)
					}
					im.IsSlice = true
				} else {
					im.Type = reqTurnFormName(v.GoName)
				}
			} else {
				tp := v.Desc.Kind().String()
				if tp == "float" {
					tp = "float32"
				} else if tp == "double" {
					tp = "float64"
				} else if v.Desc.IsList() {
					tp = "[]" + tp
					im.IsSlice = true
				}
				im.Type = tp
			}
		}
	}
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
