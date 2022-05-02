package core

import (
	"fmt"
	"project/pkg/logger"
	"reflect"
	"runtime"
	"strings"
)

func FuncName(f any) string {
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	if i := strings.LastIndex(name, "."); i > -1 {
		name = name[i+1:]
	}
	if i := strings.Index(name, "-"); i > 0 {
		name = name[:i]
	}
	return name
}

func Recover(ctx *Context) {
	if err := recover(); err != nil {
		buf := make([]byte, 2<<10)
		n := runtime.Stack(buf, false)
		logger.FromContext(ctx).Fatal("recover", err, buf[:n])
	}
}

func RecoverE(ctx *Context, e *error) {
	if err := recover(); err != nil {
		buf := make([]byte, 2<<10)
		n := runtime.Stack(buf, false)
		logger.FromContext(ctx).Fatal("recover", err, buf[:n])
		*e = fmt.Errorf("PANIC! %v", err)
	}
}
