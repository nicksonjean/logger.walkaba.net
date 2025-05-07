package logger

import (
	"fmt"
	"runtime"
)

func GetStacktrace() string {
	pc := make([]uintptr, 10)
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])

	frame, _ := frames.Next()
	return fmt.Sprintf("%s:%d", frame.File, frame.Line)
}
