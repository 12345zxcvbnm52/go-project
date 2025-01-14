package main // import "golang.org/x/tools/cmd/stringer"

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strconv"
	"strings"
)

var (
	outPath   = flag.String("code_out", "./", "生成的文件位置")
	outName   = flag.String("out_name", "code_out.go", "生成的文件名称")
	goPackage = flag.String("go_package", "", "生成文件的package名称,默认使用传入文件的package")
	inPath    = flag.String("path", "", "需要传入的error.go模板文件")
	fileName  = ""
)

func Usage() {
	fmt.Println("CodeGen发生错误,请遵循以下守则")
	s := `
	CodeGen具有以下几个flag:
		1. -code_out  可选的指定生成的文件位置 默认为当前路径./ 
		2. -out_name   可选的生成的文件名称 默认名称为code_out.go
		3. -format     可选的用于解析注释中格式化http_code和rpc_code的形式,语句内msg代表错误信息,http代表http_code,rpc代表rpc_code 默认为{msg}:{http}:{rpc}
		4. -go_package 可选的生成文件的package名称 默认使用传入文件的package
		5. -path	   必须的传入的error.go模板文件的地址
	`
	fmt.Println(s)
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
	genDecl(file.Decls)

}

func genDecl(decls []ast.Decl) {
	buf, err := os.OpenFile(*outPath+*outName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		flag.Usage()
		panic(err)
	}
	if err := os.Rename(*outPath+*outName, *outName); err != nil {
		panic(err)
	}
	for _, decl := range decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			flag.Usage()
			panic(errors.New("序列化代码失败"))
		}
		if genDecl.Tok != token.CONST {
			continue
		}
		fmt.Fprintf(buf, "package %s\n", *goPackage)
		fmt.Fprintf(buf, "import \"google.golang.org/grpc/codes\"\n\n")
		for _, spec := range genDecl.Specs {
			switch tv := spec.(type) {
			case *ast.TypeSpec:

			case *ast.ValueSpec:
				comment := ""
				if tv.Doc != nil && tv.Doc.Text() != "" {
					comment = tv.Doc.Text()
				} else if tv.Comment != nil && tv.Comment.Text() != "" {
					comment = tv.Comment.Text()
				}
				mtdata := strings.Split(comment, ":")
				mtdata[2] = strings.TrimSpace(mtdata[2])
				if mtdata[0] == "" {
					mtdata[0] = "Internal"
				}
				//如果传入的不是数字就是codes包下的错误码
				if _, err := strconv.Atoi(mtdata[0]); err != nil {
					mtdata[0] = "codes." + mtdata[0]
				}
				value, ok := tv.Values[0].(*ast.BasicLit)
				if !ok {
					flag.Usage()
					panic(errors.New("传入的文件中存在一个非const的error实例"))
				}

				fmt.Fprintf(buf,
					"var Code%s=newCoder(%s,%s,%s,\"%s\",\"%s\")\n\n",
					tv.Names[0].Name[3:],
					value.Value,
					mtdata[1],
					mtdata[0],
					mtdata[2],
					tv.Names[0].Name,
				)

				fmt.Fprintf(buf,
					"func New%s(err error, format string, args ...any) error {\n",
					tv.Names[0].Name,
				)
				fmt.Fprintf(buf,
					"	return WrapCoder(err, Code%s, format, args...)\n",
					tv.Names[0].Name[3:],
				)
				fmt.Fprintf(buf, "}\n\n")

			}
		}

	}
}
