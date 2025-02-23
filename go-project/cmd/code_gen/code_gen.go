package main // import "golang.org/x/tools/cmd/stringer"

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"os"
	"strconv"
	"strings"
)

var (
	outPath   = flag.String("code_out", "./", "生成的文件位置")
	outName   = flag.String("out_name", "code_out.go", "生成的文件名称")
	goPackage = flag.String("go_package", "", "生成文件的package名称,默认使用传入文件的package")
	inPath    = flag.String("path", "", "需要传入的error.go文件地址")
	fileName  = ""
)

//go:embed code.tmpl
var tpl string

func Usage() {
	fmt.Println("CodeGen发生错误,请遵循以下守则")
	s := `
	CodeGen具有以下几个flag:
		1. -code_out  可选的指定生成的文件位置 默认为当前路径./ 
		2. -out_name   可选的生成的文件名称 默认名称为code_out.go
		3. -go_package 可选的生成文件的package名称 默认使用传入文件的package
		4. -path	   必须的传入的error.go文件的地址
	`
	fmt.Println(s)
}

type Code struct {
	CodeName string
	HttpCode string
	GrpcCode string
	Message  string
	CodeNum  string
}

func main() {
	flag.Usage = Usage
	flag.Parse()
	if *inPath == "" {
		*inPath = "./"
	}
	args := flag.Args()
	if len(args) < 1 {
		panic("缺少指定的error.go文件")
	}
	fileName = args[0]
	file, err := parser.ParseFile(token.NewFileSet(), *inPath+fileName, nil, parser.ParseComments)
	if err != nil {
		flag.Usage()
		panic(err)
	}
	*goPackage = file.Name.Name
	codes := genDecl(file.Decls)
	GenFile(codes)
}

func genDecl(decls []ast.Decl) [][]Code {
	var res [][]Code = make([][]Code, 0)
	for _, decl := range decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			flag.Usage()
			panic(errors.New("序列化代码失败"))
		}
		if genDecl.Tok != token.CONST {
			continue
		}
		var codes []Code = make([]Code, 0)
		var IotaValue string = ""
		for _, spec := range genDecl.Specs {
			var code Code = Code{}
			switch tv := spec.(type) {
			case *ast.TypeSpec:
				continue
			case *ast.ValueSpec:
				if tv.Doc != nil && tv.Doc.Text() != "" {
					code.Message = tv.Doc.Text()
				} else if tv.Comment != nil && tv.Comment.Text() != "" {
					code.Message = tv.Comment.Text()
				}
				mtdata := strings.Split(code.Message, ":")
				code.Message = strings.TrimSpace(mtdata[2])
				//如果传入的不是数字就是codes包下的错误码
				if mtdata[0] == "" {
					code.GrpcCode = "codes.Internal"
				} else if _, err := strconv.Atoi(mtdata[0]); err != nil {
					code.GrpcCode = "codes." + mtdata[0]
				}
				//这个一定能转化成功(之前排除了TypeSpec以及非const类型)
				if tv.Values != nil && tv.Values[0] != nil {
					switch t := tv.Values[0].(type) {
					case *ast.BasicLit:
						IotaValue = t.Value
					case *ast.BinaryExpr:
						IotaValue = t.X.(*ast.BasicLit).Value
					}
				} else {
					num, _ := strconv.Atoi(IotaValue)
					IotaValue = strconv.Itoa(num + 1)
				}
				code.HttpCode = mtdata[1]
				code.CodeName = tv.Names[0].Name

				code.CodeNum = IotaValue

				codes = append(codes, code)
			}
		}
		res = append(res, codes)
	}
	return res
}

func GenFile(codesSet [][]Code) {
	file, err := os.OpenFile(*outPath+*outName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		flag.Usage()
		panic(err)
	}
	defer file.Close()
	if err := os.Rename(*outPath+*outName, *outName); err != nil {
		panic(err)
	}
	tmpl, err := template.New("http").Parse(strings.TrimSpace(tpl))
	if err != nil {
		panic(err)
	}

	if err := tmpl.Execute(file, codesSet); err != nil {
		panic(err)
	}
}
