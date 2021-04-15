package cache

import (
	"math"
	"math/rand"
	"sync"
	"time"
)

type CallBackFunc func(key interface{})

type Cache struct {
	sync.RWMutex
	internalMap map[interface{}]*Node
	list        *DoubleLinkedList
	limitSize   int

	// callback
	// callback trigger by nearly-expire key
	nearlyExpireCallBack CallBackFunc
	// A parameter that affects the calculation of nearExpire
	// nearlyExpire = expire * beta * log(rand())
	// default 0.9
	nearlyExpireBeta float64
	// callback trigger by expired key
	expireCallBack CallBackFunc
}

func NewCache(size int) *Cache {
	c := Cache{}
	c.limitSize = size
	c.internalMap = make(map[interface{}]*Node)
	c.list = NewDoubleLinkedList()
	c.nearlyExpireBeta = 0.9
	return &c
}

func (c *Cache) SetNearlyExpireCallBack(f CallBackFunc) {
	c.nearlyExpireCallBack = f
}

func (c *Cache) SetExpireCallBack(f CallBackFunc) {
	c.nearlyExpireCallBack = f
}

func (c *Cache) SetNearlyExpireBeta(beta float64) {
	c.nearlyExpireBeta = beta
}

func (c *Cache) Size() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.internalMap)
}

func (c *Cache) Set(key, value interface{}, expire time.Duration) {
	c.Lock()
	defer c.Unlock()

	var node *Node
	var ok bool
	if node, ok = c.internalMap[key]; ok {
		c.list.Remove(node)

	}
	node = c.list.PushBack(key, value, time.Now().Add(expire))
	c.internalMap[key] = node

	if c.Size() > c.limitSize {
		node = c.list.RemoveFront()
		delete(c.internalMap, node.key)
		//clean both list and map
	}
}

func (c *Cache) Get(key interface{}) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()

	node, ok := c.internalMap[key]
	nowTime := time.Now()
	if ok && nowTime.Before(node.expireTime) {
		return node.Data(), true
	}
	return nil, false
}

/*
	expire * beta * log(rand())
	beta = 0.9
*/
func CalcuNearlyExpire(expire time.Duration, beta float64) time.Duration {
	return time.Duration(float64(expire) * beta * math.Log(rand.Float64()))
}
