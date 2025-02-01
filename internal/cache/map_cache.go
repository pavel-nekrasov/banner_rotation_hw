package cache

import "sync"

type mapCache[K comparable, D any] struct {
	mu    sync.RWMutex
	items map[K]D
}

func NewMapCache[K comparable, D any]() Cache[K, D] {
	return &mapCache[K, D]{
		items: make(map[K]D),
	}
}

func (c *mapCache[K, D]) Get(key K) (D, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if item, ok := c.items[key]; ok {
		return item, ok
	}
	var result D
	return result, false
}

func (c *mapCache[K, D]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[K]D)
}

func (c *mapCache[K, D]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}

func (c *mapCache[K, D]) Empty() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items) == 0
}

func (c *mapCache[K, D]) Range(handler func(key K, data D)) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for key, data := range c.items {
		handler(key, data)
	}
}

func (c *mapCache[K, D]) Remove(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

func (c *mapCache[K, D]) Set(key K, value D) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = value
	return false
}
