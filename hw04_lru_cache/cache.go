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
	capacity int
	queue    List
	items    map[Key]*ListItem
	mu       *sync.RWMutex
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		mu:       new(sync.RWMutex),
	}
}

type pear struct {
	key Key
	val interface{}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, exist := c.items[key]
	if exist {
		v.Value = pear{key: key, val: value}
		c.queue.MoveToFront(v)
		return true
	}

	if c.queue.Len() >= c.capacity {
		deletedKey := c.queue.Back().Value.(pear).key
		c.queue.Remove(c.queue.Back())
		delete(c.items, deletedKey)
	}
	list := c.queue.PushFront(pear{key: key, val: value})
	c.items[key] = list
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, exist := c.items[key]
	if !exist {
		return nil, false
	}

	saved := v.Value.(pear)
	c.queue.MoveToFront(v)
	return saved.val, true
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
	c.mu.Unlock()
}
