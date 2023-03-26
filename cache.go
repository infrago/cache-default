package default_cache

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	. "github.com/infrago/base"
	"github.com/infrago/cache"
)

var (
	errInvalidCacheConnection = errors.New("Invalid cache connection.")
	errInvalidCacheData       = errors.New("Invalid cache data.")
)

type (
	defaultDriver  struct{}
	defaultConnect struct {
		mutex sync.RWMutex

		instance *cache.Instance
		setting  defaultSetting
		caches   sync.Map
	}
	defaultSetting struct {
	}
	defaultValue struct {
		Value  []byte
		Expiry time.Time
	}
)

// 连接
func (driver *defaultDriver) Connect(inst *cache.Instance) (cache.Connect, error) {
	setting := defaultSetting{}

	return &defaultConnect{
		instance: inst, setting: setting,
		caches: sync.Map{},
	}, nil
}

// 打开连接
func (this *defaultConnect) Open() error {
	return nil
}

// 关闭连接
func (this *defaultConnect) Close() error {
	return nil
}

// 查询缓存，
func (this *defaultConnect) Read(key string) ([]byte, error) {
	if value, ok := this.caches.Load(key); ok {
		if vv, ok := value.(defaultValue); ok {
			if vv.Expiry.Unix() > time.Now().Unix() {
				return vv.Value, nil
			} else {
				//过期了就删除
				this.Delete(key)
			}
		}
	}
	return nil, errInvalidCacheData
}

// 更新缓存
func (this *defaultConnect) Write(key string, data []byte, expiry time.Duration) error {
	now := time.Now()

	value := defaultValue{
		Value: data, Expiry: now.Add(expiry),
	}

	this.caches.Store(key, value)

	return nil
}

// 查询缓存，
func (this *defaultConnect) Exists(key string) (bool, error) {
	if _, ok := this.caches.Load(key); ok {
		return ok, nil
	}
	return false, errors.New("缓存读取失败")
}

// 删除缓存
func (this *defaultConnect) Delete(key string) error {
	this.caches.Delete(key)
	return nil
}

func (this *defaultConnect) Sequence(key string, start, step int64, expiry time.Duration) (int64, error) {
	value := start

	if data, err := this.Read(key); err == nil {
		num, err := strconv.ParseInt(string(data), 10, 64)
		if err == nil {
			value = num
		}
	}

	value += step

	//写入值
	data := []byte(fmt.Sprintf("%v", value))
	err := this.Write(key, data, expiry)
	if err != nil {
		return int64(0), err
	}

	return value, nil
}

func (this *defaultConnect) Keys(prefix string) ([]string, error) {
	keys := []string{}

	this.caches.Range(func(k, _ Any) bool {
		key := fmt.Sprintf("%v", k)

		if strings.HasPrefix(key, prefix) {
			keys = append(keys, key)
		}
		return true
	})
	return keys, nil
}
func (this *defaultConnect) Clear(prefix string) error {
	if keys, err := this.Keys(prefix); err == nil {
		for _, key := range keys {
			this.caches.Delete(key)
		}
		return nil
	} else {
		return err
	}
}
