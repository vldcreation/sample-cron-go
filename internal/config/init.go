package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
)

func ENV() interface {
	Require(required ...string) (err error)
	Get(key string) string
	Coalesce(key string, def string) string
} {
	return _env
}

type env struct {
	requireHasCalled bool

	r map[string]struct{}
	x map[string]string
}

var (
	onceEnv = new(sync.Once)
	_env    = new(env)
)

// Require only once.
func (e *env) Require(required ...string) (err error) {
	if len(required) < 1 {
		return nil
	} else if e.requireHasCalled {
		return errors.New("env: require has called")
	}

	onceEnv.Do(func() {
		e.r = make(map[string]struct{})
		for _, v := range required {
			v = strings.TrimSpace(v)
			if v != "" {
				e.r[v] = struct{}{}
			}
		}
		e.requireHasCalled = true
	})

	return nil
}

func (e *env) Coalesce(key string, def string) string {
	if out := e.Get(key); out != "" {
		return out
	}

	return def
}

func (e *env) Get(key string) string {
	if len(e.x) > 0 {
		return e.x[key]
	} else if !e.requireHasCalled {
		return os.Getenv(key)
	}

	e.x = make(map[string]string)

	for _, env := range os.Environ() {
		kv := strings.SplitN(env, "=", 2)

		if _, ok := e.r[kv[0]]; ok {
			delete(e.r, kv[0])
		}

		e.x[kv[0]] = kv[1]
	}

	for k := range e.r {
		panic(fmt.Sprintf("env: %q is not found.", k))
	}

	return e.Get(key)
}
