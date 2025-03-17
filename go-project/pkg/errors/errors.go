package errors

import (
	"fmt"
	"io"
)

// fundamental 记录基本错误信息
type fundamental struct {
	msg string
	*stack
}

// withStack 记录错误的堆栈信息
type withStack struct {
	msg   string
	cause error
	*stack
}

// withCode 记录错误码、消息及堆栈
type withCode struct {
	msg   string
	code  Coder
	cause error
	*stack
}

// causer 接口用于获取根本错误
type causer interface {
	Cause() error
}

// New 创建基础错误
func New(message string) error {
	return &fundamental{
		msg:   message,
		stack: callers(),
	}
}

// Errorf 创建带格式化信息的错误
func Errorf(format string, args ...interface{}) error {
	return &fundamental{
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
	}
}

// Wrap 包装错误,并添加消息
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return &withStack{
		msg:   message,
		cause: err,
		stack: callers(),
	}
}

// Wrapf 包装错误,并格式化消息
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &withStack{
		msg:   fmt.Sprintf(format, args...),
		cause: err,
		stack: callers(),
	}
}

// WithStack 仅添加堆栈信息
func WithStack(err error) error {
	if err == nil {
		return nil
	}
	return &withStack{
		cause: err,
		stack: callers(),
	}
}

// WithCoder 记录错误码、消息及堆栈
func WithCoder(err error, coder Coder, message string) error {
	if err == nil {
		return nil
	}
	return &withCode{
		msg:   message,
		code:  coder,
		cause: err,
		stack: callers(),
	}
}

// Cause 返回最原始的错误
func Cause(err error) error {
	for err != nil {
		if cause, ok := err.(causer); ok {
			err = cause.Cause()
		} else {
			break
		}
	}
	return err
}

// Message 获取错误信息
func Message(err error) string {
	switch e := err.(type) {
	case *fundamental:
		return e.msg
	case *withStack:
		return e.msg
	case *withCode:
		return e.Message()
	default:
		return err.Error()
	}
}

// Error 实现 fundamental 的错误消息
func (f *fundamental) Error() string {
	return f.msg
}

// Format 实现 fundamental 的格式化输出
func (f *fundamental) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintln(s, f.msg)
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

// 实现 withStack
func (w *withStack) Error() string {
	return fmt.Sprintf("%s: %v", w.msg, w.cause)
}

func (w *withStack) Unwrap() error {
	return w.cause
}

func (w *withStack) Cause() error {
	return w.cause
}

func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintln(s, w.Error())
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

// withCode 实现
func (w *withCode) Error() string {
	return fmt.Sprintf("%s: %v", w.msg, w.cause)
}

func (w *withCode) Message() string {
	if w.msg == "" {
		return w.code.Message()
	}
	return w.msg
}

func (w *withCode) Cause() error {
	return w.cause
}

func (w *withCode) Unwrap() error {
	return w.cause
}

func (w *withCode) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintln(s, w.Error())
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
