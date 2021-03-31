package cache

import (
	"time"
)

type CallBackFunc func(key interface{})

type Cache struct {
	interalMap map[interface{}]*Node
	list       *DoubleLinkedList
	limitSize  int

	// callback
	// 即将过期触发的callback
	nearlyExpireCallBack CallBackFunc
	// 已经过期触发的callback
	expireCallBack CallBackFunc
}

func NewCache(size int) *Cache {
	c := Cache{}
	c.limitSize = size
	c.interalMap = make(map[interface{}]*Node)
	c.list = NewDoubleLinkedList()
	return &c
}

func (c *Cache) Size() int {
	return len(c.interalMap)
}

func (c *Cache) Set(key, value interface{}, expire time.Duration) {

	var node *Node
	var ok bool
	if node, ok = c.interalMap[key]; ok {
		c.list.Remove(node)

	}
	node = c.list.PushBack(key, value, time.Now().Add(expire))
	c.interalMap[key] = node

	if c.Size() > c.limitSize {
		node = c.list.RemoveFront()
		delete(c.interalMap, node.key)
		// 需要同时清除，list和map
	}
}

func (c *Cache) Get(key interface{}) (interface{}, bool) {
	node, ok := c.interalMap[key]
	nowTime := time.Now()
	if ok && nowTime.Before(node.expireTime) {
		return node.Data(), true
	}
	return nil, false
}
