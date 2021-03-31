package cache

import (
	"fmt"
	"time"
)

type Node struct {
	key   interface{}
	value interface{}
	// 过期时间 纳秒
	expireTime time.Time
	// ---------
	// 前向指针
	Prev *Node
	// 后向指针
	Next *Node
}

func (n *Node) Data() interface{} {
	if n == nil {
		return nil
	}
	return n.value
}

func (n *Node) PrevData() interface{} {
	if n.Prev == nil {
		return nil
	}
	return n.Prev.value
}

func (n *Node) NextData() interface{} {
	if n.Next == nil {
		return nil
	}
	return n.Next.value
}

// 双向链表
// 尾部的Node更新
// 如果进行LRU删除，删除头部的数据
type DoubleLinkedList struct {
	count int
	Head  *Node
	Tail  *Node
}

func NewDoubleLinkedList() *DoubleLinkedList {
	return &DoubleLinkedList{count: 0}
}

func (l *DoubleLinkedList) IsEmpty() bool {
	if l.count > 0 {
		return false
	}
	return true
}

func (l *DoubleLinkedList) Size() int {
	return l.count
}

func (l *DoubleLinkedList) RemoveFront() *Node {
	result := l.Head
	if l.Size() == 1 {
		l.Head = nil
		l.Tail = nil
	} else {
		if l.Size() > 0 {
			l.Head = l.Head.Next
			l.Head.Prev = nil
		}
	}
	l.count--
	return result
}

func (l *DoubleLinkedList) Remove(n *Node) *Node {
	if n == l.Head {
		l.RemoveFront()
	} else {
		preNode := n.Prev
		nextNode := n.Next
		preNode.Next = nextNode
		if nextNode != nil {
			nextNode.Prev = preNode
		}
		if l.Tail == n {
			l.Tail = preNode
		}
		l.count--
	}
	return n
}

func (l *DoubleLinkedList) PushBack(key, value interface{}, expire time.Time) *Node {
	node := &Node{key: key, value: value, expireTime: expire}
	if l.count <= 0 {
		l.Head = node
		l.Tail = node
	} else {
		l.Tail.Next = node
		node.Prev = l.Tail
		l.Tail = node
	}
	l.count++
	return node
}

func (l *DoubleLinkedList) TraversalPrint() {
	curr := l.Head
	fmt.Println("count:", l.count, "head:", l.Head.Data(), "tail:", l.Tail.Data())
	for curr != nil {
		fmt.Println("current", curr.Data(), "prev", curr.PrevData(), "next", curr.NextData())
		curr = curr.Next
	}
}
