package errors

import (
	"errors"
	"fmt"
)

// MessageCountMap 包含每个错误消息的出现次数,
type MessageCountMap map[string]int

// ErrorGroup 表示一个包含多个错误的对象,但这些错误不一定具有单一的语义意义,
// 可以使用 errors.Is 来检查是否存在特定类型的错误,
// 不支持 errors.As,因为多个错误并不总是具有单一的类型,不支持将这些错误转换为特定类型的指针,
type ErrorGroup interface {
	error
	Errors() []error
	Is(error) bool
}

// NewErrorGroup 将一组错误转换为一个ErrorGroup接口,
// 该接口本身实现了error接口,如果输入的错误列表为空,它会返回 nil,
// 它还会检查输入的错误列表是否包含nil,以避免在调用Error()时发生nil指针panic,
func NewErrorGroup(errlist []error) ErrorGroup {
	if len(errlist) == 0 {
		return nil
	}
	// 如果输入的错误列表中包含nil
	var errs []error
	for _, e := range errlist {
		if e != nil {
			errs = append(errs, e)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errorGroup(errs)
}

type errorGroup []error

func (agg errorGroup) Error() string {
	// 保险,但一般情况下不可能满足条件
	if len(agg) == 0 {
		return ""
	}
	if len(agg) == 1 {
		return agg[0].Error()
	}
	seenerrs := NewString()
	result := ""
	agg.visit(func(err error) bool {
		msg := err.Error()
		if seenerrs.Has(msg) {
			return false
		}
		seenerrs.Insert(msg)
		if len(seenerrs) > 1 {
			result += ", "
		}
		result += msg
		return false
	})
	if len(seenerrs) == 1 {
		return result
	}
	return "[" + result + "]"
}

// Is 是 ErrorGroup 接口的一部分,
func (agg errorGroup) Is(target error) bool {
	return agg.visit(func(err error) bool {
		return errors.Is(err, target)
	})
}

func (agg errorGroup) visit(f func(err error) bool) bool {
	for _, err := range agg {
		switch err := err.(type) {
		case errorGroup:
			if match := err.visit(f); match {
				return match
			}
		case ErrorGroup:
			for _, nestedErr := range err.Errors() {
				if match := f(nestedErr); match {
					return match
				}
			}
		default:
			if match := f(err); match {
				return match
			}
		}
	}

	return false
}

// Errors 是 ErrorGroup 接口的一部分,
func (agg errorGroup) Errors() []error {
	return []error(agg)
}

// Matcher 用于匹配错误,返回 true 表示错误匹配,
type Matcher func(error) bool

// FilterOut 从输入错误中移除所有匹配任何匹配器的错误,
// 如果输入的是单一错误,则只测试该错误,如果输入实现了 ErrorGroup 接口,
// 则会递归处理错误列表,
//
// 例如,它可以用来移除已知的无效错误（如 io.EOF 或 os.PathNotFound）,
func FilterOut(err error, fns ...Matcher) error {
	if err == nil {
		return nil
	}
	if agg, ok := err.(ErrorGroup); ok {
		return NewErrorGroup(filterErrors(agg.Errors(), fns...))
	}
	if !matchesError(err, fns...) {
		return err
	}
	return nil
}

// matchesError 如果任何 Matcher 返回 true,则返回 true
func matchesError(err error, fns ...Matcher) bool {
	for _, fn := range fns {
		if fn(err) {
			return true
		}
	}
	return false
}

// filterErrors 返回所有 fns 返回 false 的错误（或嵌套错误,
// 如果列表中包含嵌套的 Errors）,如果没有错误剩余,则返回 nil 列表,
// 该函数的副作用是扁平化所有嵌套的错误列表,
func filterErrors(list []error, fns ...Matcher) []error {
	result := []error{}
	for _, err := range list {
		r := FilterOut(err, fns...)
		if r != nil {
			result = append(result, r)
		}
	}
	return result
}

// Flatten 接受一个 ErrorGroup,该 ErrorGroup 可能包含其他嵌套的 Aggregates,
// 并将它们全部扁平化成一个单一的 ErrorGroup,递归地处理,
func Flatten(agg ErrorGroup) ErrorGroup {
	result := []error{}
	if agg == nil {
		return nil
	}
	for _, err := range agg.Errors() {
		if a, ok := err.(ErrorGroup); ok {
			r := Flatten(a)
			if r != nil {
				result = append(result, r.Errors()...)
			}
		} else {
			if err != nil {
				result = append(result, err)
			}
		}
	}
	return NewErrorGroup(result)
}

// CreateAggregateFromMessageCountMap 将 MessageCountMap 转换为 ErrorGroup
func CreateAggregateFromMessageCountMap(m MessageCountMap) ErrorGroup {
	if m == nil {
		return nil
	}
	result := make([]error, 0, len(m))
	for errStr, count := range m {
		var countStr string
		if count > 1 {
			countStr = fmt.Sprintf(" (repeated %v times)", count)
		}
		result = append(result, fmt.Errorf("%v%v", errStr, countStr))
	}
	return NewErrorGroup(result)
}

// Reduce 将返回 err,或者,如果 err 是一个 ErrorGroup 且只有一个元素,
// 则返回 ErrorGroup 中的第一个元素,
func Reduce(err error) error {
	if agg, ok := err.(ErrorGroup); ok && err != nil {
		switch len(agg.Errors()) {
		case 1:
			return agg.Errors()[0]
		case 0:
			return nil
		}
	}
	return err
}

// AggregateGoroutines 并行运行提供的函数,将所有非 nil 错误收集到返回的 ErrorGroup 中,
// 如果所有函数成功执行,则返回 nil,
func AggregateGoroutines(funcs ...func() error) ErrorGroup {
	errChan := make(chan error, len(funcs))
	for _, f := range funcs {
		go func(f func() error) { errChan <- f() }(f)
	}
	errs := make([]error, 0)
	for i := 0; i < cap(errChan); i++ {
		if err := <-errChan; err != nil {
			errs = append(errs, err)
		}
	}
	return NewErrorGroup(errs)
}

// ErrPreconditionViolated 当违反前置条件时返回的错误
var ErrPreconditionViolated = errors.New("precondition is violated")
