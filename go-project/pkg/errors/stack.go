package errors

import (
	"fmt"
	"path"
	"runtime"
	"strings"
)

// Frame 代表一个调用栈帧
type Frame struct {
	pc       uintptr
	function string
	file     string
	line     int
}

// 创建 Frame 时直接解析文件名,行号,函数名,避免多次 runtime.FuncForPC
func newFrame(pc uintptr) Frame {
	pc--
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return Frame{pc: pc, function: "unknown", file: "unknown", line: 0}
	}
	file, line := fn.FileLine(pc)
	return Frame{pc: pc, function: fn.Name(), file: file, line: line}
}

// 格式化输出 Frame
func (f Frame) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		if s.Flag('+') {
			fmt.Fprintf(s, "%s\n\t%s", f.function, f.file)
		} else {
			fmt.Fprint(s, path.Base(f.file))
		}
	case 'd':
		fmt.Fprint(s, f.line)
	case 'n':
		fmt.Fprint(s, funcname(f.function))
	case 'v':
		f.Format(s, 's')
		fmt.Fprint(s, ":")
		f.Format(s, 'd')
	}
}

// 实现 `MarshalText`
func (f Frame) MarshalText() ([]byte, error) {
	if f.function == "unknown" {
		return []byte(f.function), nil
	}
	return []byte(fmt.Sprintf("%s %s:%d", f.function, f.file, f.line)), nil
}

// stack 代表调用栈
type stack []Frame

// 格式化 stack
func (s *stack) Format(st fmt.State, verb rune) {
	if verb == 'v' && st.Flag('+') {
		for _, f := range *s {
			fmt.Fprintf(st, "%+v\n", f)
		}
	}
}

// 获取调用栈
func callers() *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:]) // 2 以获取调用 callers() 之上的栈帧

	st := make(stack, 0, n)
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		name := fn.Name()

		// 过滤掉 runtime 和 asm 相关的帧
		if !strings.Contains(name, "runtime.") && !strings.Contains(name, "asm_") {
			st = append(st, newFrame(pc))
		}
	}

	return &st
}

// 解析函数名,去除路径前缀
func funcname(name string) string {
	if i := strings.LastIndex(name, "/"); i != -1 {
		name = name[i+1:]
	}
	if i := strings.Index(name, "."); i != -1 {
		name = name[i+1:]
	}
	return name
}
