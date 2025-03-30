package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	isset := false
	if v, ok := c.items[key]; ok {
		isset = true
		v.Value = value
		c.queue.MoveToFront(v)
	}

	if !isset {
		if c.queue.Len() == c.capacity {
			c.queue.Remove(c.queue.Back())
		}

		c.items[key] = c.queue.PushFront(value)
	}

	return isset
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	if v, ok := c.items[key]; ok {
		c.mu.Lock()
		defer c.mu.Unlock()

		c.queue.MoveToFront(v)
		return v.Value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
