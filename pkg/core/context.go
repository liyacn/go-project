package core

import (
	"context"
	"os"
	"project/pkg/random"
	"sync"
)

type Context struct {
	mu sync.RWMutex
	kv map[string]any
	context.Context
}

func (c *Context) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.kv[key] = value
}

func (c *Context) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists := c.kv[key]
	return value, exists
}

func (c *Context) Value(key any) any {
	if ks, ok := key.(string); ok {
		if value, exists := c.Get(ks); exists {
			return value
		}
	}
	return c.Context.Value(key)
}

func NewContext(v0, v1, v2, v3 string) *Context {
	c := &Context{
		kv:      make(map[string]any),
		Context: context.Background(),
	}
	if v0 == "" {
		v0 = random.UUID()
	}
	if v1 == "" && len(os.Args) > 1 {
		v1 = os.Args[1]
	}
	c.Set("v0", v0)
	c.Set("v1", v1)
	c.Set("v2", v2)
	c.Set("v3", v3)
	return c
}
