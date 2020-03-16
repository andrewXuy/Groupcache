package lru

import "container/list"
// FIFO LRU, LFU
// FIFO will drop the oldest data
// LFU least frequently used 好处：命中高，缺点：维护队列，并且收到历史影响比较大
// LRU Least recent used 如果最近访问过，放到队尾，淘汰队首 map + double linked list

type Cache struct {
	// Max allowed memory size
	maxBytes int64
	curBytes int64
	// using library double linked list
	ll *list.List
	// <String, The pointer in the linked list>
	cache map[string]*list.Element
	// execute when an entry is purged
	Remove func(key string, value Value)

}

// entry has interface Value, is one of its implementation
type entry struct {
	key string
	value Value
}

// Frequency of a Key
type Value interface{
	Len() int
}

// Constructor of Cache
func New(maxBytes int64, Remove func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		curBytes: 0,
		ll:       list.New(),
		cache:    make(map[string]*list.Element),
		Remove:   Remove,
	}
}

// Look up a key's value
func (c *Cache) Get(key string) (value Value, ok bool) {
	if element, ok := c.cache[key]; ok{
		// Move the node to end
		c.ll.MoveToFront(element)
		// Explicitly convert
		keyV := element.Value.(*entry)
		return keyV.value,true
	}
	return
}

// Remove the least frequent used element
func (c *Cache) RemoveOldest(){
	// The First Element is the least frequent Element
	element := c.ll.Back()
	if element != nil {
		// Remove this element
		c.ll.Remove(element)
		keyV := element.Value.(*entry)
		delete(c.cache, keyV.key)

		// Remove the size of <key,count> pair
		c.curBytes -= int64(len(keyV.key)) + int64(keyV.value.Len())
		if c.Remove != nil {
			c.Remove(keyV.key,keyV.value)
		}

	}
}

func (c *Cache) Add(key string, value Value) {
	// key already in map
	if element, ok := c.cache[key]; ok {
		// Move the element to t
		c.ll.MoveToFront(element)
		keyV := element.Value.(*entry)
		// update
		c.curBytes += int64(value.Len()) - int64(keyV.value.Len())
		keyV.value = value
	} else {
		// create new key after the last pointer
		element := c.ll.PushFront(&entry{key,value})
		c.cache[key] = element
		c.curBytes += int64(len(key)) + int64(value.Len())
	}
	// if it exceed the size, keep removing util curBytes <= maxBytes
	// 0 means unlimited size
	for c.maxBytes!=0 && c.curBytes > c.maxBytes {
		c.RemoveOldest()
	}
}

// Number of cache entry
func (c *Cache) Len() int {
	return c.ll.Len()
}
