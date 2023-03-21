package default_cache

import (
	"github.com/infrago/cache"
	"github.com/infrago/infra"
)

func Driver() cache.Driver {
	return &defaultDriver{}
}

func init() {
	infra.Register("default", Driver())
}
