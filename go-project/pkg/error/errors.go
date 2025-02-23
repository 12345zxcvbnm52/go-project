package errors

import (
	"fmt"
	"io"
)

// 提供记录异常栈和只记录最外层错误信息的错误,类似于单一cause
type fundamental struct {
	msg string
	*stack
}

// 提供记录异常栈和多层次错误信息(cause)功能的错误,类似于链条cause
type withStack struct {
	msg   string
	bferr error
	*stack
}

// 提供记录异常栈,单一cause,链条cause功能,自定义code功能的错误,该Error为最高级,不允许被其它error覆盖
type withCode struct {
	bferr error
	msg   string
	code  Coder
	cause error
	*stack
}

type causer interface {
	Cause() error
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

// 类似于WithStack,不过Wrap还可以往新一层错误栈中添加message
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*withCode); ok {
		return &withCode{
			bferr: e,
			msg:   message,
			code:  e.code,
			cause: e.cause,
			stack: callers(),
		}
	}

	return &withStack{
		msg:   message,
		bferr: err,
		stack: callers(),
	}
}

func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*withCode); ok {
		return &withCode{
			bferr: e,
			msg:   fmt.Sprintf(format, args...),
			code:  e.code,
			cause: e.cause,
			stack: callers(),
		}
	}

	return &withStack{
		bferr: err,
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
	}
}

func WithCoder(err error, coder Coder, message string) error {
	if err == nil {
		return nil
	}
	return &withCode{
		msg:   message,
		bferr: err,
		code:  coder,
		stack: callers(),
		cause: err,
	}
}

// Cause 返回记录了根本原因(即内部的敏感的信息)的error
// 如果一个错误值有根本原因,它需要实现以下接口
//
//	type causer interface {
//	       Cause() error
//	}
func Cause(err error) error {
	if err == nil {
		return err
	}

	if e, ok := err.(*withCode); ok {
		return e.cause
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		e := cause.Cause()
		err = e
	}
	return err
}

func Message(err error) string {
	if cerr, ok := err.(Coder); ok {
		return cerr.Message()
	}
	switch cerr := err.(type) {
	case *fundamental:
		return cerr.msg
	case *withStack:
		return cerr.msg
	default:
		return err.Error()
	}
}

// fundamental的Error函数默认是不会打印stack信息的
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
func WithStack(err error, message string) error {
	if err == nil {
		return nil
	}

	// //单独处理codeError,其为最高级Error
	// if e, ok := err.(*withCode); ok {
	// 	return &withCode{
	// 		msg:   message,
	// 		bferr: e,
	// 		code:  e.code,
	// 		cause: e.cause,
	// 		stack: callers(),
	// 	}
	// }

	return &withStack{
		msg:   message,
		bferr: err,
		stack: callers(),
	}
}

func (w *withStack) Cause() error { return w.bferr }

// WithStack的Error函数默认也不会打印栈信息和各个栈内错误信息
func (w *withStack) Error() string { return w.msg }

// 用于兼容stderr中的Unwarp
// 记得处理WithStack的Format函数,保证每一层的error信息都能被打印出来
func (w *withStack) Unwrap() error { return w.bferr }

func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintln(s, w.Error())
			Cause(w).(*fundamental).stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}

// 返回不敏感的error message
func (w *withCode) Message() string {
	if w.msg == "" {
		return w.code.Message()
	}
	return w.msg
}

// 返回用户不敏感的error message
func (w *withCode) Error() string { return w.Cause().Error() }

// 返回原始cause error
func (w *withCode) Cause() error {
	//如果传入的codeError没有最基础的err就默认其本身是原始error
	// if w.cause == nil {
	// 	return w
	// }
	// //如果传入的codeError有最基础的err就递归找到最原始的error
	// causeErr, ok := w.cause.(causer)
	// if !ok {
	// 	return w.cause
	// }
	// return causeErr.Cause()
	return w.cause
}

func (w *withCode) Unwrap() error {
	return w.bferr
}

func (w *withCode) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintln(s, w.Error())
			w.cause.(*fundamental).stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}
