package cache

import "sync"

type LruCache[K comparable, D interface{}] interface {
	Set(key K, value D) bool
	Get(key K) (D, bool)
	Remove(key K)
	Range(handler func(key K, data D))
	Clear()
}

type lruEntry[K comparable, D interface{}] struct {
	Key  K
	Data D
}

type lruCache[K comparable, D interface{}] struct {
	mu       sync.Mutex
	capacity int
	queue    LinkedList[lruEntry[K, D]]
	items    map[K]*ListItem[lruEntry[K, D]]
}

func NewCache[K comparable, D interface{}](capacity int) LruCache[K, D] {
	return &lruCache[K, D]{
		capacity: capacity,
		queue:    NewList[lruEntry[K, D]](),
		items:    make(map[K]*ListItem[lruEntry[K, D]], capacity),
	}
}

func (c *lruCache[K, D]) Set(key K, value D) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	data := lruEntry[K, D]{Key: key, Data: value}

	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		item.Value = data
		return true
	}

	c.storeItem(key, c.queue.PushFront(data))
	if c.queue.Len() > c.capacity {
		c.removeItem(c.queue.Back())
	}

	return false
}

func (c *lruCache[K, D]) Get(key K) (D, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		return item.Value.Data, ok
	}

	var result D
	return result, false
}

func (c *lruCache[K, D]) Remove(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.items[key]; ok {
		c.removeItem(item)
	}
}

func (c *lruCache[K, D]) Range(handler func(key K, data D)) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, listItem := range c.items {
		handler(key, listItem.Value.Data)
	}
}

func (c *lruCache[K, D]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.queue = NewList[lruEntry[K, D]]()
	c.items = make(map[K]*ListItem[lruEntry[K, D]], c.capacity)
}

func (c *lruCache[K, D]) storeItem(key K, item *ListItem[lruEntry[K, D]]) {
	c.items[key] = item
}

func (c *lruCache[K, D]) removeItem(item *ListItem[lruEntry[K, D]]) {
	data := item.Value

	delete(c.items, data.Key)
	c.queue.Remove(item)
}
