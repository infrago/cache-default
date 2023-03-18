package default_cache

import (
	"github.com/infrago/cache"
)

func Driver() cache.Driver {
	return &defaultDriver{}
}

func init() {
	cache.Register("default", Driver())
}
