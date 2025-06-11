package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

func StackTrace(skip int) string {
	const maxDepth = 32
	pcs := make([]uintptr, maxDepth)
	n := runtime.Callers(skip, pcs)
	frames := runtime.CallersFrames(pcs[:n])

	var sb strings.Builder
	for {
		frame, more := frames.Next()
		file := filepath.Base(frame.File)
		sb.WriteString(fmt.Sprintf("%s:%d\t%s\n", file, frame.Line, frame.Function))
		if !more {
			break
		}
	}
	return sb.String()
}
