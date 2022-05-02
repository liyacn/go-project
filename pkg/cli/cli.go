package cli

import "project/pkg/core"

type HandlerFunc func(*core.Context)

func Wrap(h HandlerFunc) func() {
	v2 := core.FuncName(h)
	return func() {
		ctx := core.NewContext("", "", v2, "")
		defer core.Recover(ctx)
		h(ctx)
	}
}
