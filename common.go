package openai

import (
	"os"
	"strconv"
	"sync"
)

var (
	_debug     bool
	_debugOnce sync.Once
)

func Debug() bool {
	_debugOnce.Do(func() {
		_debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	})
	return _debug
}
