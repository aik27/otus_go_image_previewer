package lrucache

import "sync"

type Key string

type cacheItem struct {
	key   Key
	value interface{}
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mu       sync.Mutex
	capacity int               // Количество сохраняемых в кэше элементов.
	queue    List              // Очередь на основе двусвязного списка.
	items    map[Key]*ListItem // Словарь, отображающий ключ (строка) на элемент очереди.
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

// Set Добавить значение в кэш по ключу.
func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if ok {
		item.Value.(*cacheItem).value = value
		c.queue.MoveToFront(item)
		return true
	}

	newItem := &cacheItem{key: key, value: value}
	c.items[key] = c.queue.PushFront(newItem)

	if c.queue.Len() > c.capacity {
		backItem := c.queue.Back()
		delete(c.items, backItem.Value.(*cacheItem).key)
		c.queue.Remove(backItem)
	}

	return false
}

// Get Получить значение из кэша по ключу.
func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(item)
		return item.Value.(*cacheItem).value, true
	}
	return nil, false
}

// Clear Очистить кэш.
func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
