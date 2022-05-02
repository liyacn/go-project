package worker

import "project/pkg/core"

type HandlerFunc func(*core.Context)

func Wrap(h HandlerFunc) func() {
	return func() {
		ctx := core.NewContext("", "", "", "")
		defer core.Recover(ctx)
		h(ctx)
	}
}
