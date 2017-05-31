package app

import (
	"errors"
	"runtime"
	"sync"
	"time"
)

type item interface {
	getExpired() int64
}

type simpleItem struct {
	object  interface{}
	expired int64
}

type listItem struct {
	listObject []interface{}
	expired    int64
}

type dictItem struct {
	dictObject map[string]interface{}
	expired    int64
}

func (si simpleItem) getExpired() int64 {
	return si.expired
}

func (li listItem) getExpired() int64 {
	return li.expired
}

func (di dictItem) getExpired() int64 {
	return di.expired
}

type cache struct {
	items                 map[string]item
	mu                    sync.RWMutex
	janitor               *janitor
	ExpiredTimeMultiplier time.Duration
}

func NewCache(interval time.Duration) *cache {
	if interval == 0 {
		interval = time.Duration(10 * time.Millisecond)
	}

	c := &cache{
		items: make(map[string]item),
		ExpiredTimeMultiplier: time.Second,
	}
	runJanitor(c, interval)
	runtime.SetFinalizer(c, stopJanitor)

	return c
}

func (c *cache) set(key string, value interface{}, duration int) bool {
	var e int64
	if duration > 0 {
		e = time.Now().Add(time.Duration(duration) * c.ExpiredTimeMultiplier).UnixNano()
	}

	c.mu.Lock()
	c.items[key] = simpleItem{
		object:  value,
		expired: e,
	}
	c.mu.Unlock()

	return true
}

func (c *cache) get(key string) (interface{}, error) {
	c.mu.RLock()
	item, found := c.items[key]

	if !found {
		c.mu.RUnlock()
		return nil, errors.New("not found")
	}

	if item.getExpired() > 0 {
		if time.Now().UnixNano() > item.getExpired() {
			c.mu.RUnlock()
			return nil, errors.New("not found")
		}
	}

	si, ok := item.(simpleItem)
	if !ok {
		c.mu.RUnlock()
		return nil, errors.New("wrong object type")
	}

	c.mu.RUnlock()
	return si.object, nil
}

func (c *cache) keys() []string {
	c.mu.RLock()

	var keys []string

	for k := range c.items {
		keys = append(keys, k)
	}

	c.mu.RUnlock()

	return keys
}

func (c *cache) deleteItem(key string) {
	delete(c.items, key)
}

func (c *cache) rpush(key string, value interface{}, duration int) (bool, error) {
	c.mu.Lock()
	item, found := c.items[key]
	var e int64
	if duration > 0 {
		e = time.Now().Add(time.Duration(duration) * c.ExpiredTimeMultiplier).UnixNano()
	}

	if !found {
		var object []interface{}
		object = append(object, value)

		li := listItem{
			expired:    e,
			listObject: object,
		}

		c.items[key] = li

		c.mu.Unlock()

		return true, nil
	}

	li, ok := item.(listItem)

	if !ok {
		return false, errors.New("invalid object type")
	}

	li.listObject = append(li.listObject, value)
	if duration > 0 {
		li.expired = e
	}

	c.items[key] = li

	c.mu.Unlock()

	return true, nil
}

func (c *cache) lgetall(key string) ([]interface{}, error) {
	c.mu.RLock()

	item, found := c.items[key]

	if !found {
		c.mu.RUnlock()
		return nil, errors.New("not found")
	}

	if item.getExpired() > 0 {
		if time.Now().UnixNano() > item.getExpired() {
			c.mu.RUnlock()
			return nil, errors.New("not found")
		}
	}

	li, ok := item.(listItem)

	if !ok {
		c.mu.RUnlock()
		return nil, errors.New("wrong type")
	}
	c.mu.RUnlock()

	return li.listObject, nil
}

func (c *cache) lget(key string, id int) (interface{}, error) {
	c.mu.RLock()
	item, found := c.items[key]

	if !found {
		c.mu.RUnlock()
		return nil, errors.New("not found")
	}

	li, ok := item.(listItem)

	if !ok {
		c.mu.RUnlock()
		return nil, errors.New("wrong type")
	}

	if len(li.listObject) < id+1 {
		c.mu.RUnlock()
		return nil, errors.New("not found")
	}

	value := li.listObject[id]

	c.mu.RUnlock()

	return value, nil

}

func (c *cache) pop(key string) (interface{}, error) {
	c.mu.Lock()
	item, found := c.items[key]

	if !found {
		c.mu.Unlock()
		return nil, errors.New("not found")
	}

	li, ok := item.(listItem)

	if !ok {
		c.mu.Unlock()
		return nil, errors.New("wrong type")
	}

	var object interface{}
	object, li.listObject = li.listObject[len(li.listObject)-1], li.listObject[:len(li.listObject)-1]

	if len(li.listObject) == 0 {
		c.mu.Unlock()
		delete(c.items, key)
		return object, nil
	}

	c.items[key] = li

	c.mu.Unlock()

	return object, nil
}

func (c *cache) hset(key string, value map[string]interface{}, duration int) error {
	c.mu.Lock()
	item, found := c.items[key]

	var e int64
	if duration > 0 {
		e = time.Now().Add(time.Duration(duration) * c.ExpiredTimeMultiplier).UnixNano()
	}

	if !found {
		object := map[string]interface{}{}
		for k, v := range value {
			object[k] = v
		}

		di := dictItem{
			expired:    e,
			dictObject: object,
		}
		c.items[key] = di

		c.mu.Unlock()

		return nil
	}

	di, ok := item.(dictItem)

	if !ok {
		return errors.New("invalid object type")
	}

	for k, v := range value {
		di.dictObject[k] = v
	}

	if duration > 0 {
		di.expired = e
	}

	c.items[key] = di

	c.mu.Unlock()

	return nil
}

func (c *cache) hgetall(key string) (map[string]interface{}, error) {
	c.mu.RLock()
	item, found := c.items[key]

	if !found {
		c.mu.RUnlock()
		return nil, errors.New("not found")
	}

	if item.getExpired() > 0 {
		if time.Now().UnixNano() > item.getExpired() {
			c.mu.RUnlock()
			return nil, errors.New("not found")
		}
	}

	di, ok := item.(dictItem)

	if !ok {
		c.mu.RUnlock()
		return nil, errors.New("wrong type")
	}
	c.mu.RUnlock()

	return di.dictObject, nil
}

func (c *cache) hget(key string, dictKey string) (interface{}, error) {
	c.mu.RLock()
	item, found := c.items[key]

	if !found {
		c.mu.RUnlock()
		return nil, errors.New("not found")
	}

	if item.getExpired() > 0 {
		if time.Now().UnixNano() > item.getExpired() {
			c.mu.RUnlock()
			return nil, errors.New("not found")
		}
	}

	di, ok := item.(dictItem)

	if !ok {
		c.mu.RUnlock()
		return nil, errors.New("wrong type")
	}

	value, ok := di.dictObject[dictKey]

	if !ok {
		c.mu.RUnlock()
		return nil, errors.New("not found")
	}

	c.mu.RUnlock()

	return value, nil
}

func (c *cache) DeleteExpired() {
	now := time.Now().UnixNano()
	c.mu.Lock()

	for k, v := range c.items {
		if v.getExpired() > 0 && now > v.getExpired() {
			delete(c.items, k)
		}
	}

	c.mu.Unlock()
}
