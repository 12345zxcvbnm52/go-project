package errors
import "google.golang.org/grpc/codes"

var CodeCanceled=newCoder(100101,499,codes.Canceled,"客户端关闭请求或连接超时","ErrCanceled")

func NewErrCanceled(err error, format string, args ...any) error {
	return WrapCoder(err, CodeCanceled, format, args...)
}

var CodeRequestTimeout=newCoder(100102,408,codes.Unavailable,"客户端请求超时,请稍后再试","ErrRequestTimeout")

func NewErrRequestTimeout(err error, format string, args ...any) error {
	return WrapCoder(err, CodeRequestTimeout, format, args...)
}

