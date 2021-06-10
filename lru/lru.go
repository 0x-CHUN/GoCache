package lru

import "container/list"

type Cache struct {
	maxBytes  int64
	numBytes  int64
	List      *list.List
	cache     map[string]*list.Element
	onEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		List:      list.New(),
		cache:     make(map[string]*list.Element),
		onEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (Value, bool) {
	if elem, ok := c.cache[key]; ok {
		c.List.MoveToFront(elem)
		kv := elem.Value.(*entry)
		return kv.value, true
	}
	return nil, false
}

func (c *Cache) RemoveOldest() {
	elem := c.List.Back()
	if elem != nil {
		c.List.Remove(elem)
		kv := elem.Value.(*entry)
		delete(c.cache, kv.key)
		c.numBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.onEvicted != nil {
			c.onEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	if elem, ok := c.cache[key]; ok {
		c.List.MoveToFront(elem)
		kv := elem.Value.(*entry)
		c.numBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		elem := c.List.PushFront(&entry{key: key, value: value})
		c.cache[key] = elem
		c.numBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.numBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.List.Len()
}
