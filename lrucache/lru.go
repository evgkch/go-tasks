//go:build !solution

package lrucache

import "container/list"

type entry struct {
	key   int
	value int
}

type lruCache struct {
	cap   int
	cache map[int]*list.Element
	list  *list.List
}

func New(cap int) Cache {
	return &lruCache{
		cap:   cap,
		cache: make(map[int]*list.Element, cap),
		list:  list.New(),
	}
}

func (c *lruCache) Get(key int) (int, bool) {
	if v, ok := c.cache[key]; ok {
		c.list.MoveToFront(v)
		return v.Value.(*entry).value, true
	}
	return 0, false
}

func (c *lruCache) Set(key, value int) {
	if v, ok := c.cache[key]; ok {
		v.Value.(*entry).value = value
		c.list.MoveToFront(v)
		return
	}

	v := c.list.PushFront(&entry{key, value})
	c.cache[key] = v

	if c.list.Len() > c.cap {
		oldest := c.list.Back()
		if oldest != nil {
			c.list.Remove(oldest)
			delete(c.cache, oldest.Value.(*entry).key)
		}
	}
}

func (c *lruCache) Range(f func(key, value int) bool) {
	for elem := c.list.Back(); elem != nil; elem = elem.Prev() {
		entry := elem.Value.(*entry)
		if !f(entry.key, entry.value) {
			break
		}
	}
}

func (c *lruCache) Clear() {
	clear(c.cache)
	c.list = list.New()
}
