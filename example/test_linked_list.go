package main

import (
	"fmt"
	cache "github.com/vearne/localcache"
	"time"
)

func main() {
	list := cache.NewDoubleLinkedList()
	nowTime := time.Now()
	exTime := nowTime.Add(time.Hour)
	list.PushBack(1, 1, exTime)
	fmt.Println("----1----", list, list.Head.Data(), list.Tail.Data())

	for i := 2; i < 10; i++ {
		list.PushBack(i, i, exTime)
	}
	fmt.Println("----2----", list, list.Head.Data(), list.Tail.Data())
	list.TraversalPrint()

	var x *cache.Node
	target := 3
	p := list.Head
	for target != p.Data() {
		p = p.Next
	}
	x = p
	list.Remove(x)
	//fmt.Println("----3----")
	//list.TraversalPrint()
	//fmt.Println("----4----")
	//list.RemoveFront()
	//list.TraversalPrint()
	//for list.Size() > 0 {
	//	fmt.Println("--------")
	//	list.RemoveFront()
	//	list.TraversalPrint()
	//}
	for list.Size() > 0 {
		fmt.Println("--------")
		list.Remove(list.Tail)
		list.TraversalPrint()
	}
}
