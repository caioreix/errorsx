package errorsx

import (
	"fmt"
	"runtime"
	"strings"
)

type StackFrame struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

type Stack []*StackFrame

func (sf StackFrame) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s\n\t%s:%d\n", sf.Function, sf.File, sf.Line))
	return sb.String()
}

func (s Stack) String() string {
	var sb strings.Builder
	for _, sf := range s {
		sb.WriteString(sf.String())
	}
	return sb.String()
}

func getStack(skip int) Stack {
	buf := make([]uintptr, 32)
	n := runtime.Callers(skip, buf[:])
	stack := make([]uintptr, n)
	copy(stack, buf[:n])

	s := make(Stack, 0, n)
	frames := runtime.CallersFrames(stack)
	for {
		frame, more := frames.Next()
		s = append(s, &StackFrame{
			Function: frame.Function,
			File:     frame.File,
			Line:     frame.Line,
		})
		if !more {
			break
		}
	}
	return s
}
