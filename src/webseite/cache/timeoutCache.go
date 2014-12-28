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
	items      map[interface{}]*list.Element
	lock       sync.Mutex
}

// entry is used to hold a value in the evictList
type entry struct {
	key   interface{}
	value interface{}
	added int64
}

func NewTimeoutCache(timeout int64) (*TimeoutCache, error) {
	beego.BeeLogger.Info("New timeoutCache")

	if timeout <= 0 {
		return nil, errors.New("Must provide a positive timeout (in seconds)")
	}

	c := &TimeoutCache{
		evictAfter: timeout,
		evictList:  list.New(),
		items:      make(map[interface{}]*list.Element),
	}

	ticker := time.NewTicker(time.Millisecond * 100)
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
	c.lock.Lock()
	defer c.lock.Unlock()

	if ent, ok := c.items[key]; ok {
		c.evictList.MoveToFront(ent)
		return ent.Value.(*entry).value, true
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
		if time-kv.added > c.evictAfter {
			c.removeElement(e)
		}
	}
}

// Add adds a value to the cache.
func (c *TimeoutCache) Add(key, value interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Check for existing item
	if ent, ok := c.items[key]; ok {
		c.evictList.MoveToFront(ent)
		ent.Value.(*entry).value = value
		return
	}

	// Add new item
	ent := &entry{key, value, time.Now().Unix()}
	entry := c.evictList.PushFront(ent)
	c.items[key] = entry
}

// removeElement is used to remove a given list element from the cache
func (c *TimeoutCache) removeElement(e *list.Element) {
	c.evictList.Remove(e)
	kv := e.Value.(*entry)
	delete(c.items, kv.key)
}

// Remove removes the provided key from the cache.
func (c *TimeoutCache) Remove(key interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if ent, ok := c.items[key]; ok {
		c.removeElement(ent)
	}
}

// Len returns the number of items in the cache.
func (c *TimeoutCache) Len() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.evictList.Len()
}
