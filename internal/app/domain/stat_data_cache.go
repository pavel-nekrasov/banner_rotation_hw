package appdomain

import (
	"sync"

	"github.com/pavel-nekrasov/banner_rotation_hw/internal/cache"
)

type StatDataCache struct {
	items   cache.Cache[GroupSlotKey, *StatData]
	mutexes sync.Map
}

func NewStatDataCache(capacity int) *StatDataCache {
	return &StatDataCache{
		items: cache.NewLRUCache[GroupSlotKey, *StatData](capacity),
	}
}

func (c *StatDataCache) Lock(key GroupSlotKey) {
	value, _ := c.mutexes.LoadOrStore(key, &sync.Mutex{})
	keyMutex, ok := value.(*sync.Mutex)
	if !ok {
		panic("stored value not a mutex ptr")
	}

	keyMutex.Lock()
}

func (c *StatDataCache) Unlock(key GroupSlotKey) {
	value, ok := c.mutexes.Load(key)
	if !ok {
		panic("key mutex was not found to unlock")
	}

	keyMutex, ok := value.(*sync.Mutex)
	if !ok {
		panic("stored value not a mutex ptr")
	}

	keyMutex.Unlock()
}

func (c *StatDataCache) Get(key GroupSlotKey) (*StatData, bool) {
	if item, ok := c.items.Get(key); ok {
		return item, ok
	}
	return nil, false
}

func (c *StatDataCache) Remove(key GroupSlotKey) {
	if _, ok := c.items.Get(key); ok {
		c.items.Remove(key)
	}
}

func (c *StatDataCache) Set(key GroupSlotKey, value *StatData) bool {
	c.items.Set(key, value)
	return false
}

func (c *StatDataCache) Range(handler func(key GroupSlotKey, data *StatData)) {
	c.items.Range(handler)
}
