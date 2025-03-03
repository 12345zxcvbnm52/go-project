# 基于github.com/pkg/errors包
增加对error code的支持,完全兼容github.com/pkg/errors,性能跟github.com/pkg/errors基本持平
该errors包匹配的错误码设计请参考:(https:github.com/marmotedu/sample-code/blob/master/README.md)

# errors包提供了简单的错误处理工具
errors.Wrap函数会返回一个新错误,通过在调用Wrap时记录堆栈追踪,
并添加指定的消息,为原始错误增加上下文例如:
	_, err := ioutil.ReadAll(r)
	if err != nil {
	        return errors.Wrap(err, "读取失败")
	}
如果需要更细粒度的控制,可以使用errors.WithStack和errors.WithMessage函数,
将errors.Wrap分解为两个基本操作:为错误添加堆栈追踪和附加消息

# 获取错误的根本原因(或者说获得最初始的错误)
 使用errors.Wrap会构造一个错误堆栈,为之前的错误添加上下文
 根据错误的性质,可能需要反转errors.Wrap的操作以提取原始错误进行检查
 任何实现了以下接口的错误值:
	type causer interface {
	        Cause() error
	}
 都可以通过errors.Cause进行检查,errors.Cause会递归地检索内层未实现causer接口的错误,并假设其为原始错误
	switch err := errors.Cause(err).(type) {
	case *MyError:
	         特定处理
	default:
	         未知错误
	}

# 格式化打印错误
此包返回的所有错误值都实现了fmt.Formatter接口,可以通过fmt包进行格式化
支持以下格式化符:
	%s    打印错误如果错误具有 Cause,将递归打印
	%v    等同于 %s
	%+v   扩展格式错误堆栈追踪的每一帧都会详细打印

# 获取错误或包装器的堆栈追踪

 New、Errorf、Wrap 和 Wrapf 会在调用时记录堆栈追踪
 此信息可以通过以下接口获取:

	type stackTracer interface {
	        StackTrace() errors.StackTrace
	}

 返回的 errors.StackTrace 类型定义为:

	type StackTrace []Frame

 Frame 类型表示堆栈追踪中的一个调用点
 Frame 支持 fmt.Formatter 接口,可以用于打印有关此错误堆栈追踪的信息例如:

	if err, ok := err.(stackTracer); ok {
	        for _, f := range err.StackTrace() {
	                fmt.Printf("%+s:%d\n", f, f)
	        }
	}

虽然 stackTracer 接口未被此包导出,但它被视为其稳定公共接口的一部分
有关 Frame.Format 的更多细节,请参阅其文档
