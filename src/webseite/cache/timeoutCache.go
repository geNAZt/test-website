package cache

import (
	"container/list"
	"errors"
	"github.com/astaxie/beego"
	"sync"
	"time"
)

type TimeoutCache struct {
	evictAfter int64
	evictList  *list.List
	items      map[interface{}]interface{}
	lock       sync.RWMutex
}

// entry is used to hold a value in the evictList
type entry struct {
	key   interface{}
	value interface{}
	evict int64
}

func NewTimeoutCache(timeout int64) (*TimeoutCache, error) {
	beego.BeeLogger.Info("New timeoutCache")

	if timeout <= 0 {
		return nil, errors.New("Must provide a positive timeout (in seconds)")
	}

	c := &TimeoutCache{
		evictAfter: timeout,
		evictList:  list.New(),
		items:      make(map[interface{}]interface{}),
	}

	ticker := time.NewTicker(time.Millisecond * 1000)
	go func() {
		for {
			select {
			case <-ticker.C:
				c.autoEvict()
			}
		}
	}()

	return c, nil
}

// Get looks up a key's value from the cache.
func (c *TimeoutCache) Get(key interface{}) (value interface{}, ok bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if ent, ok := c.items[key]; ok {
		return ent, true
	}

	return nil, false
}

// Remove all timeout Entries
func (c *TimeoutCache) autoEvict() {
	c.lock.Lock()
	defer c.lock.Unlock()

	time := time.Now().Unix()

	for e := c.evictList.Back(); e != nil; e = e.Prev() {
		kv := e.Value.(*entry)
		if time > kv.evict {
			c.removeElement(e)
		}
	}
}

// Add adds a value to the cache.
func (c *TimeoutCache) Add(key, value interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Check for existing item
	if ent, ok := c.items[key]; ok && ent != nil {
		c.items[key] = value
		return
	}

	// Add new item
	ent := &entry{key, value, time.Now().Add(time.Duration(c.evictAfter * int64(time.Second))).Unix()}
	c.evictList.PushFront(ent)
	c.items[key] = value
}

// removeElement is used to remove a given list element from the cache
func (c *TimeoutCache) removeElement(e *list.Element) {
	c.evictList.Remove(e)
	kv := e.Value.(*entry)
	delete(c.items, kv.key)
}

// Search for an element with a key
func (c *TimeoutCache) searchElement(key interface{}) *list.Element {
	for e := c.evictList.Back(); e != nil; e = e.Prev() {
		kv := e.Value.(*entry)
		if kv.key == key {
			return e
		}
	}

	return nil
}

// Remove removes the provided key from the cache.
func (c *TimeoutCache) Remove(key interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if ele := c.searchElement(key); ele != nil {
		c.removeElement(ele)
	}
}

// Len returns the number of items in the cache.
func (c *TimeoutCache) Len() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.evictList.Len()
}
