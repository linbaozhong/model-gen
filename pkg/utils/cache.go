// @Title 应用内的cache库
// @Description
// @Author 蔺保仲 2020/04/20
// @Update 蔺保仲 2020/04/20
package utils

import (
	"sync"
	"time"
)

type Cache struct {
	//table   cmap.ConcurrentMap
	table   sync.Map
	expired time.Duration
}
type cacheItem struct {
	val     interface{}
	expired time.Time
}

var caches = make([]*Cache, 0, 10)

//func init() {
//	ticker := time.NewTicker(10 * time.Second)
//
//	go func() {
//		for {
//			t := <-ticker.C
//			for i := 0; i < len(caches); i++ {
//				//if caches[i].table.Count() == 0 {
//				//	continue
//				//}
//				//caches[i].table.IterCb(func(k string, v interface{}) {
//				//	if val, ok := v.(cacheItem); ok {
//				//		if val.expired.Before(t) {
//				//			caches[i].table.Remove(k)
//				//		}
//				//	}
//				//})
//				caches[i].table.Range(func(k, v interface{}) bool {
//					if val, ok := v.(cacheItem); ok {
//						if val.expired.Before(t) {
//							caches[i].table.Delete(k)
//						}
//					}
//					return true
//				})
//			}
//		}
//	}()
//}

// 实例化cache，生命期缺省60秒
func NewCache(expired ...time.Duration) *Cache {
	exp := 60 * time.Second

	if len(expired) > 0 {
		exp = expired[0]
	}
	cache := &Cache{
		expired: exp,
		table:   sync.Map{},
		//table:   cmap.New(),
	}
	caches = append(caches, cache)
	return cache
}

func (c *Cache) Set(key string, value interface{}, expired ...time.Duration) {
	exp := c.expired
	if len(expired) > 0 {
		exp = expired[0]
	}
	//c.table.Set(key, cacheItem{val: value, expired: time.Now().Add(exp)})
	c.table.Store(key, cacheItem{val: value, expired: time.Now().Add(exp)})
}

func (c *Cache) Get(key string) (interface{}, bool) {
	if val, ok := c.table.Load(key); ok {
		if val, ok := val.(cacheItem); ok {
			if val.expired.After(time.Now()) {
				return val.val, true
			}
			c.Delete(key)
		}
	}

	return nil, false
}

func (c *Cache) Delete(key string) {
	c.table.Delete(key)
}

//func (c *Cache) Items() map[string]interface{} {
//	return c.table.Items()
//}

// 清空cache
func (c *Cache) Empty() {
	c.table = sync.Map{}
}
