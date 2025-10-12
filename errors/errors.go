package myerrors

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

var (
	colorReset = "\033[0m"
	colorRed   = "\033[31m" // only red for errors
)

var goInternalPaths = []string{
	"go/pkg/",
	"go/src/",
	"rest/",
	"/middleware",
}

type Error struct {
	err   error
	frame string
}

func New(text string) error {
	return &Error{
		err:   errors.New(text),
		frame: fileWithLineNumber(),
	}
}

func shouldIgnore(path string) bool {
	for _, p := range goInternalPaths {
		if strings.Contains(path, p) {
			return true
		}
	}
	return false
}

func ErrorfOneLine(format string, a ...any) error {
	LogError(3, 1, format, a...)

	return &Error{
		err:   fmt.Errorf(format, a...),
		frame: fileWithLineNumber(),
	}
}

func Errorf(format string, a ...any) error {
	LogError(3, 20, format, a...)

	return &Error{
		err:   fmt.Errorf(format, a...),
		frame: fileWithLineNumber(),
	}
}

func Wrap(err error, text string) error {
	if err == nil {
		return nil
	}
	return &Error{
		err:   fmt.Errorf("%s: %w", text, err),
		frame: fileWithLineNumber(),
	}
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	return e.err.Error()
}

func (e *Error) ErrorWithFrame() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("%s %s", e.frame, e.err.Error())
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}

// fileWithLineNumber returns the file and line number of the caller's caller.
// For example: "/path/to/file.go:123"
func fileWithLineNumber() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	// return string(strconv.AppendInt(append([]byte(frame.File), ':'), int64(frame.Line), 10))
	// return fmt.Sprintf("%s:%d", file, line)
	return string(strconv.AppendInt(append([]byte(file), ':'), int64(line), 10))
}

func LogErrorOneLine(skip int, msg string, args ...interface{}) {
	// 1 = caller of this function, 2 = caller’s caller, etc.
	_, file, line, ok := runtime.Caller(skip)
	if ok {
		fmt.Printf("%s:%d: %s\n", file, line, fmt.Sprintf(msg, args...))
	} else {
		fmt.Printf("%s\n", fmt.Sprintf(msg, args...))
	}
}

func LogError(skip, depth int, msg string, args ...interface{}) {
	pc := make([]uintptr, 32) // capture stack
	n := runtime.Callers(skip, pc)
	frames := runtime.CallersFrames(pc[:n])

	count := 0
	for {
		frame, more := frames.Next()
		if !shouldIgnore(frame.File) {
			if count == 0 {
				// first caller: include error message here
				fmt.Printf("%s%s:%d %s%s\n",
					colorRed, frame.File, frame.Line,
					fmt.Sprintf(msg, args...), colorReset,
				)
			} else {
				// subsequent callers: just show file:line + func
				fmt.Printf("%s%s:%d %s\n",
					colorRed, frame.File, frame.Line, colorReset,
				)
			}
			count++
			if count >= depth {
				break
			}
		}
		if !more {
			break
		}
	}
}
