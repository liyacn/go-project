package core

import (
	"fmt"
	"project/pkg/logger"
	"runtime"
)

func Recover(ctx *Context) {
	if err := recover(); err != nil {
		buf := make([]byte, logger.MaxBodyBytes)
		n := runtime.Stack(buf, false)
		logger.FromContext(ctx).Fatal("recover", err, buf[:n])
	}
}

func RecoverE(ctx *Context, e *error) {
	if err := recover(); err != nil {
		buf := make([]byte, logger.MaxBodyBytes)
		n := runtime.Stack(buf, false)
		logger.FromContext(ctx).Fatal("recover", err, buf[:n])
		*e = fmt.Errorf("PANIC! %v", err)
	}
}
