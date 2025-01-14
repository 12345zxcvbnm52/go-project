package errors

import (
	"fmt"
	"net/http"
	"sync"

	rpcCode "google.golang.org/grpc/codes"
)

var (
	unknownCoder defaultCoder = defaultCoder{0,
		http.StatusInternalServerError, rpcCode.Internal,
		"An internal server error occurred", "http://imooc/mxshop/pkg/errors/README.md"}
)

// Coder暴露出一个error code必须要的接口
type Coder interface {
	// 返回error code 映射的http code
	HTTPCode() int
	// 返回给用户的不敏感的信息
	Message() string

	// 返回文档给用户查看
	Reference() string

	// 返回error code
	ErrorCode() int

	//返回映射的RpcCode
}

type defaultCoder struct {
	// 对应的Code
	Code int

	// 对应的http code
	HttpCode int

	//对应的rpc code
	RpcCode rpcCode.Code

	// 给用户的不敏感的信息
	Msg string

	// 给用户的文档
	Refer string
}

func (coder defaultCoder) ErrorCode() int {
	return coder.Code

}

func (coder defaultCoder) Message() string {
	return coder.Msg
}

func (coder defaultCoder) HTTPCode() int {
	if coder.HttpCode == 0 {
		return http.StatusInternalServerError
	}

	return coder.HttpCode
}

func (coder defaultCoder) Reference() string {
	return coder.Refer
}

// 一组记录code的map元数据
var codeMap = map[int]Coder{}
var codeMtx = &sync.Mutex{}

func newCoder(code int, httpCode int, rpcCode rpcCode.Code, msg string, ref string) Coder {
	coder := defaultCoder{code, httpCode, rpcCode, msg, ref}
	Register(coder)
	return coder
}

func Register(coder Coder) {
	if coder.ErrorCode() == 0 {
		panic("错误码0已被ErrUnknown占用")
	}
	codeMtx.Lock()
	defer codeMtx.Unlock()

	codeMap[coder.ErrorCode()] = coder
}

func MustRegister(coder Coder) {
	if coder.ErrorCode() == 0 {
		panic("错误码0已被ErrUnknown占用")
	}

	codeMtx.Lock()
	defer codeMtx.Unlock()

	if _, ok := codeMap[coder.ErrorCode()]; ok {
		panic(fmt.Sprintf("code: %d 已被占用,占用的Err为: %s", coder.ErrorCode(), codeMap[coder.ErrorCode()]))
	}
	codeMap[coder.ErrorCode()] = coder
}

func ParseToCoder(err error) Coder {
	if err == nil {
		return nil
	}

	if v, ok := err.(*withCode); ok {
		if coder, ok := codeMap[v.code.ErrorCode()]; ok {
			return coder
		}
	}

	return unknownCoder
}

func IsCode(err error, code int) bool {
	if v, ok := err.(*withCode); ok {
		if v.code.ErrorCode() == code {
			return true
		}

		if v.cause != nil {
			return IsCode(v.cause, code)
		}

		return false
	}

	return false
}

func GetCoder(code int) Coder {
	if code < 100101 {
		return unknownCoder
	}
	coder, ok := codeMap[code]
	if !ok {
		return unknownCoder
	}
	return coder
}

func init() {
	codeMap[unknownCoder.ErrorCode()] = unknownCoder
}
