package errors

import (
	"net/http"

	grpccode "google.golang.org/grpc/codes"
)

var (
// unknownCoder Coder = newCoder(0, http.StatusInternalServerError, grpccode.Internal, "已被占用的code为0的错误码,请误使用该错误码")
)

// Coder暴露出一个error code必须要的接口
type Coder interface {
	// 返回error code 映射的http code
	HTTPCode() int

	// 返回给用户的不敏感的信息
	Message() string

	// 返回error code
	ErrorCode() int

	//返回映射的RpcCode
	GrpcCode() grpccode.Code
}

type defaultCoder struct {
	// 对应的Code
	code int

	// 对应的http code
	httpCode int

	//对应的rpc code
	grpcCode grpccode.Code

	// 给用户的不敏感的信息
	message string
}

func (coder defaultCoder) ErrorCode() int {
	return coder.code

}

func (coder defaultCoder) Message() string {
	return coder.message
}

func (coder defaultCoder) HTTPCode() int {
	if coder.httpCode == 0 {
		return http.StatusInternalServerError
	}

	return coder.httpCode
}

func (coder defaultCoder) GrpcCode() grpccode.Code {
	return coder.grpcCode
}

// 一组记录code的map元数据
var codeMap = map[int]defaultCoder{}

func mustNewCoder(code int, httpCode int, rpcCode grpccode.Code, msg string) Coder {
	if code < 1000000 || httpCode < 200 || rpcCode < grpccode.OK {
		panic(Errorf("错误的code参数,其中code为:%d,httpCode为:%d,grpcCode为:%d", code, httpCode, rpcCode))
	}
	coder := defaultCoder{code, httpCode, rpcCode, msg}

	codeMap[code] = coder
	return coder
}

func ExtractCoderFromError(err error) Coder {
	if err == nil {
		return nil
	}

	if v, ok := err.(*withCode); ok {
		return v.code
	}
	return nil
}

func HasCode(err error, code int) bool {
	if v, ok := err.(*withCode); ok {
		if v.code.ErrorCode() == code {
			return true
		}

		if v.cause != nil {
			return HasCode(v.cause, code)
		}

		return false
	}

	return false
}
