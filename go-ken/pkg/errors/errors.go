// errors 包提供了简单的错误处理工具
// errors.Wrap函数会返回一个新错误,通过在调用Wrap时记录堆栈追踪,
// 并添加指定的消息,为原始错误增加上下文例如:
//
//	_, err := ioutil.ReadAll(r)
//	if err != nil {
//	        return errors.Wrap(err, "读取失败")
//	}
//
// 如果需要更细粒度的控制,可以使用errors.WithStack和errors.WithMessage函数,
// 将errors.Wrap分解为两个基本操作:为错误添加堆栈追踪和附加消息
//
// # 获取错误的根本原因
// 使用errors.Wrap会构造一个错误堆栈,为之前的错误添加上下文
// 根据错误的性质,可能需要反转errors.Wrap的操作以提取原始错误进行检查
// 任何实现了以下接口的错误值:
//
//	type causer interface {
//	        Cause() error
//	}
//
// 都可以通过 errors.Cause 进行检查errors.Cause 会递归地检索
// 顶层不实现causer接口的错误,并假设其为原始错误例如:
//
//	switch err := errors.Cause(err).(type) {
//	case *MyError:
//	        // 特定处理
//	default:
//	        // 未知错误
//	}
//
// 虽然 causer 接口并未被此包导出,但它被视为其稳定公共接口的一部分
//
// # 格式化打印错误
//
// 此包返回的所有错误值都实现了 fmt.Formatter 接口,可以通过 fmt 包进行格式化
// 支持以下格式化符:
//
//	%s    打印错误如果错误具有 Cause,将递归打印
//	%v    等同于 %s
//	%+v   扩展格式错误堆栈追踪的每一帧都会详细打印
//
// # 获取错误或包装器的堆栈追踪
//
// New、Errorf、Wrap 和 Wrapf 会在调用时记录堆栈追踪
// 此信息可以通过以下接口获取:
//
//	type stackTracer interface {
//	        StackTrace() errors.StackTrace
//	}
//
// 返回的 errors.StackTrace 类型定义为:
//
//	type StackTrace []Frame
//
// Frame 类型表示堆栈追踪中的一个调用点
// Frame 支持 fmt.Formatter 接口,可以用于打印有关此错误堆栈追踪的信息例如:
//
//	if err, ok := err.(stackTracer); ok {
//	        for _, f := range err.StackTrace() {
//	                fmt.Printf("%+s:%d\n", f, f)
//	        }
//	}
//
// 虽然 stackTracer 接口未被此包导出,但它被视为其稳定公共接口的一部分
//
// 有关 Frame.Format 的更多细节,请参阅其文档

package errors

import (
	"fmt"
	"io"
)

// 只提供记录异常栈和基本的msg记录功能的错误
type fundamental struct {
	msg string
	*stack
}

// 提供记录异常栈和cause功能的错误
type withStack struct {
	error
	*stack
}

// 提供记录异常栈,cause功能,自定义code功能的错误
type withCode struct {
	err   error
	code  Coder
	cause error
	*stack
}

type withMessage struct {
	cause error
	msg   string
}

// New 返回一个基于message的错误,同时这个错误会记录出错的函数栈点信息
func New(message string) error {
	return &fundamental{
		msg:   message,
		stack: callers(),
	}
}

// 在new的基础上格式化了信息,注意这里待测式多层error的嵌套
func Errorf(format string, args ...interface{}) error {
	return &fundamental{
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
	}
}

func (f *fundamental) Error() string { return f.msg }

func (f *fundamental) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, f.msg)
			f.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, f.msg)
	case 'q':
		fmt.Fprintf(s, "%q", f.msg)
	}
}

// WithStack在调用WithStack的位置为错误err添加一层堆栈追踪
func WithStack(err error) error {
	if err == nil {
		return nil
	}

	//嵌套一层error
	if e, ok := err.(*withCode); ok {
		return &withCode{
			err:   e.err,
			code:  e.code,
			cause: err,
			stack: callers(),
		}
	}

	return &withStack{
		err,
		callers(),
	}
}

func (w *withStack) Cause() error { return w.error }

// 用于兼容stderr中的Unwarp
func (w *withStack) Unwrap() error {
	if e, ok := w.error.(interface{ Unwrap() error }); ok {
		return e.Unwrap()
	}

	return w.error
}

func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", w.Cause())
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}

// 类似于WithStack,不过Wrap还可以往新一层错误栈中添加message
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*withCode); ok {
		return &withCode{
			err:   fmt.Errorf(message),
			code:  e.code,
			cause: err,
			stack: callers(),
		}
	}

	err = &withMessage{
		cause: err,
		msg:   message,
	}
	return &withStack{
		err,
		callers(),
	}
}

func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*withCode); ok {
		return &withCode{
			err:   fmt.Errorf(format, args...),
			code:  e.code,
			cause: err,
			stack: callers(),
		}
	}

	err = &withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
	return &withStack{
		err,
		callers(),
	}
}

// 将对应的错误封装一层并添加一段message
func WithMessage(err error, message string) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		cause: err,
		msg:   message,
	}
}

func WithMessagef(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
}

func (w *withMessage) Error() string { return w.msg }
func (w *withMessage) Cause() error  { return w.cause }
func (w *withMessage) Unwrap() error { return w.cause }

func (w *withMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", w.Cause())
			io.WriteString(s, w.msg)
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, w.Error())
	}
}

func WithCode(code int, format string, args ...interface{}) error {
	return &withCode{
		err:   fmt.Errorf(format, args...),
		code:  GetCoder(code),
		stack: callers(),
	}
}

func WithCoder(coder Coder, format string, args ...interface{}) error {
	return &withCode{
		err:   fmt.Errorf(format, args...),
		code:  coder,
		stack: callers(),
	}
}

func WrapCode(err error, code int, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return &withCode{
		err:   fmt.Errorf(format, args...),
		code:  GetCoder(code),
		cause: err,
		stack: callers(),
	}
}

func WrapCoder(err error, coder Coder, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	return &withCode{
		err:   fmt.Errorf(format, args...),
		code:  coder,
		cause: err,
		stack: callers(),
	}
}

// // 返回不敏感的error message
// func (w *withCode) Error() string { return fmt.Sprintf("%v", w) }

// 返回不敏感的error message
func (w *withCode) Error() string { return w.code.Message() }

type causer interface {
	Cause() error
}

// 返回原始cause error
func (w *withCode) Cause() error {
	//如果传入的codeError没有最基础的err就默认其本身是原始error
	if w.cause == nil {
		return w
	}
	//如果传入的codeError有最基础的err就递归找到最原始的error
	causeErr, ok := w.cause.(causer)
	if !ok {
		return w.cause
	}
	return causeErr.Cause()

}

func (w *withCode) Unwrap() error {
	return w.cause
}

// Cause 返回记录了根本原因(即内部的敏感的信息)的error
// 如果一个错误值有根本原因,它需要实现以下接口
//
//	type causer interface {
//	       Cause() error
//	}
//
// 如果错误没有实现 Cause 方法,则返回原始错误如果错误为 nil,
// 则直接返回 nil,而不进行进一步的检查
func Cause(err error) error {
	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}

		if cause.Cause() == nil {
			break
		}

		err = cause.Cause()
	}
	return err
}
