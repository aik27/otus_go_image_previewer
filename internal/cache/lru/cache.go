package lrucache

import (
	"sync"
)

type Key string

type CacheItem struct {
	Key   Key
	Value interface{}
}

type ImageItem struct {
	FilePath    string
	Width       int
	Height      int
	OriginalURL string
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mu        sync.Mutex
	capacity  int                   // Количество сохраняемых в кэше элементов.
	queue     List                  // Очередь на основе двусвязного списка.
	items     map[Key]*ListItem     // Словарь, отображающий ключ (строка) на элемент очереди.
	onEvicted func(item *CacheItem) // Функция обратного вызова при удалении элемента из кэша.
}

func NewCache(capacity int, onEvicted func(item *CacheItem)) Cache {
	return &lruCache{
		capacity:  capacity,
		queue:     NewList(),
		items:     make(map[Key]*ListItem, capacity),
		onEvicted: onEvicted,
	}
}

// Set Добавить значение в кэш по ключу.
func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if ok {
		item.Value.(*CacheItem).Value = value
		c.queue.MoveToFront(item)
		return true
	}

	newItem := &CacheItem{Key: key, Value: value}
	c.items[key] = c.queue.PushFront(newItem)

	if c.queue.Len() > c.capacity {
		backItem := c.queue.Back()
		delete(c.items, backItem.Value.(*CacheItem).Key)
		c.queue.Remove(backItem)

		if c.onEvicted != nil {
			c.onEvicted(backItem.Value.(*CacheItem))
		}
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
		return item.Value.(*CacheItem).Value, true
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
