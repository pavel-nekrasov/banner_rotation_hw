package cache

import "sync"

type lruEntry[K comparable, D interface{}] struct {
	Key  K
	Data D
}

type lruCache[K comparable, D interface{}] struct {
	mut      sync.RWMutex
	capacity int
	queue    LinkedList[lruEntry[K, D]]
	items    map[K]*ListItem[lruEntry[K, D]]
}

func NewLRUCache[K comparable, D interface{}](capacity int) Cache[K, D] {
	return &lruCache[K, D]{
		capacity: capacity,
		queue:    NewList[lruEntry[K, D]](),
		items:    make(map[K]*ListItem[lruEntry[K, D]], capacity),
	}
}

func (c *lruCache[K, D]) Set(key K, value D) bool {
	c.mut.Lock()
	defer c.mut.Unlock()

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
	c.mut.Lock()
	defer c.mut.Unlock()

	if item, ok := c.items[key]; ok {
		c.queue.MoveToFront(item)
		return item.Value.Data, ok
	}

	var result D
	return result, false
}

func (c *lruCache[K, D]) Len() int {
	c.mut.RLock()
	defer c.mut.RUnlock()

	return len(c.items)
}

func (c *lruCache[K, D]) Empty() bool {
	c.mut.RLock()
	defer c.mut.RUnlock()

	return len(c.items) == 0
}

func (c *lruCache[K, D]) Remove(key K) {
	c.mut.Lock()
	defer c.mut.Unlock()

	if item, ok := c.items[key]; ok {
		c.removeItem(item)
	}
}

func (c *lruCache[K, D]) Range(handler func(key K, data D)) {
	c.mut.Lock()
	defer c.mut.Unlock()

	for key, listItem := range c.items {
		handler(key, listItem.Value.Data)
	}
}

func (c *lruCache[K, D]) Clear() {
	c.mut.Lock()
	defer c.mut.Unlock()

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
