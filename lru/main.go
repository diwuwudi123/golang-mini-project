package main

import "log"

type Node struct {
	Key, Val  int
	next, pre *Node
}
type DoubleList struct {
	head, tail *Node
	len        int
}

var (
	allMap    = make(map[int]*Node)
	cap       int
	listcache = newDoubleList()
)

func main() {
	cap = 2
	put(1, 1)
	put(2, 2)
	put(3, 3)
	put(4, 4)
	put(5, 5)
	res1 := get(4)
	res2 := get(5)
	res3 := get(3)
	log.Println(res1, res2, res3)
}
func newDoubleList() *DoubleList {
	cache := &DoubleList{}
	cache.head = &Node{}
	cache.tail = &Node{}
	cache.head.next = cache.tail
	cache.tail.pre = cache.head
	cache.len = 0
	return cache
}
func (list *DoubleList) AddFirst(node *Node) {
	node.next = list.head.next
	node.pre = list.head
	list.head.next.pre = node
	list.head.next = node
	list.len++
}

func (list *DoubleList) AddLast(node *Node) {
	node.pre = list.tail.pre
	node.next = list.tail
	list.tail.next = node
	list.tail.pre = node
	list.len++
}

func (list *DoubleList) RemoveX(node *Node) {
	node.pre.next = node.next
	node.next.pre = node.pre
	list.len--
}

func (list *DoubleList) RemoveLast() *Node {
	if list.head.next == list.tail {
		return nil
	}
	last := list.tail.pre
	list.RemoveX(last)
	return last
}
func (list *DoubleList) Len() int {
	return list.len
}

func get(key int) interface{} {
	val, ok := allMap[key]
	if !ok {
		return nil
	}
	//把数据移动到列头 并且返回数据
	put(key, val.Val)
	return allMap[key].Val

}
func put(key, val int) {
	node := &Node{Key: key, Val: val}
	//已经存在的 直接替换到表头
	if val, ok := allMap[key]; ok {
		//删掉这个节点 插入到头节点
		listcache.RemoveX(val)
		listcache.AddFirst(node)
	} else {
		//不存在的 直接插入到表头
		if cap == listcache.Len() {
			//删除最后一个元素
			last := listcache.RemoveLast()
			if last != nil {
				delete(allMap, last.Key)
			}
		}
		//把当前元素追加到表头
		listcache.AddFirst(node)
		allMap[key] = node
	}
}
