package gcache

import (
	"sync"
	"time"
)

//CacheDriver
type CacheDriver interface {
	Set(key string, val interface{}, ttl time.Duration) error
	Get(key string) (interface{}, bool, error)
}

//CacheParam
type CacheParam struct {
	Key string
	ttl time.Duration
}

type memoryPair struct {
	value    interface{}
	expireAt time.Time
}
type memoryDriver struct {
	container map[string]memoryPair
	mutex     sync.RWMutex
}

func NewCacheParam(key string, ttl time.Duration) (string, CacheParam) {
	return CacheOption, CacheParam{
		Key: key,
		ttl: ttl,
	}
}

//NewMemoryDriver 未完善,不建议使用
func NewMemoryDriver() *memoryDriver {
	driver := &memoryDriver{
		container: make(map[string]memoryPair),
	}
	return driver
}

//Set
func (m *memoryDriver) Set(key string, val interface{}, ttl time.Duration) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.container[key] = memoryPair{
		value:    val,
		expireAt: time.Now().Add(ttl),
	}
	return nil
}

//Get
func (m *memoryDriver) Get(key string) (interface{}, bool, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	pair, exists := m.container[key]
	if !exists || pair.expireAt.Before(time.Now()) {
		return nil, false, nil
	}

	return pair.value, true, nil
}
