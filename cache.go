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

	nearlyExpire struct {
		TimeToCompute time.Duration
		Beta          float64
	}

	// callback trigger by expired key
	expireCallBack CallBackFunc
}

func NewCache(size int) *Cache {
	c := Cache{}
	c.limitSize = size
	c.internalMap = make(map[interface{}]*Node)
	c.list = NewDoubleLinkedList()
	c.nearlyExpire.Beta = 2
	c.nearlyExpire.TimeToCompute = 100 * time.Millisecond
	return &c
}

func (c *Cache) SetNearlyExpireCallBack(f CallBackFunc) {
	c.nearlyExpireCallBack = f
}

func (c *Cache) SetExpireCallBack(f CallBackFunc) {
	c.nearlyExpireCallBack = f
}

func (c *Cache) SetNearlyExpireBeta(beta float64) {
	c.nearlyExpire.Beta = beta
}
func (c *Cache) SetNearlyExpireTimeToCompute(d time.Duration) {
	c.nearlyExpire.TimeToCompute = d
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
expire + ( timeToCompute * beta * log(rand()) )
currentTime 是当前时间戳。
timeToCompute 是重新计算缓存值所花费的时间。
beta 是一个大于 0 的非负数，默认值为 1，是可配置的。
rand()是一个返回 0 到 1 之间随机数的函数。
expiry 是缓存值未来被设置为过期的时间戳。
*/
func CalcuNearlyExpire(expire time.Duration, timeToCompute time.Duration, beta float64) time.Duration {
	return expire - time.Duration(float64(timeToCompute)*beta*math.Log(rand.Float64()/2+0.1))
}
