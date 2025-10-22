package list

import (
	"fmt"
	"sync"

	"github.com/kennycch/gotools/general"
)

// NewList 创建列表对象
func NewList[T any, V general.Ordered](f func(element T) V) (listStruct *List[T, V]) {
	listStruct = &List[T, V]{
		lists:    make([]T, 0),
		lock:     &sync.RWMutex{},
		sortFunc: f,
	}
	return
}

// Insert 新增对象
func (l *List[T, V]) Insert(element T) (index int) {
	l.lock.Lock()
	defer l.lock.Unlock()
	// 压入元素
	l.push(element)
	return index
}

// Remove 删除对象
func (l *List[T, V]) Remove(index int) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.lists = general.DeleteValueByKey(l.lists, index)
}

// GetElements 获取指定范围的对象
func (l *List[T, V]) GetElements(start, num int) (elements []T) {
	elements = make([]T, 0)
	lenght := len(l.lists)
	if lenght == 0 || start < 0 || start >= lenght {
		return
	}
	l.lock.RLock()
	defer l.lock.RUnlock()
	if num < 0 {
		num = lenght
	}
	for i := start; i < len(l.lists); i++ {
		elements = append(elements, l.lists[i])
		num--
		if num <= 0 {
			break
		}
	}
	return
}

// GetLen 获取链表长度
func (l *List[T, V]) GetLen() int {
	return len(l.lists)
}

// PopFront 弹出链表首部对象
func (l *List[T, V]) PopFront() (element T, err error) {
	if len(l.lists) == 0 {
		err = fmt.Errorf("list is empty")
		return
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	element = l.lists[0]
	l.lists = general.DeleteValueByKey(l.lists, 0)
	return
}

// PopBack 弹出链表尾部对象
func (l *List[T, V]) PopBack() (element T, err error) {
	if l.GetLen() == 0 {
		err = fmt.Errorf("list is empty")
		return
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	lastIndex := l.GetLen() - 1
	element = l.lists[lastIndex]
	l.lists = general.DeleteValueByKey(l.lists, lastIndex)
	return
}

// ClearList 清空链表
func (l *List[T, V]) ClearList() {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.lists = l.lists[:0]
}

// Iterator 元素迭代器
func (l *List[T, V]) Iterator(f func(index int, element T) bool) {
	l.lock.Lock()
	defer l.lock.Unlock()
	for k, element := range l.lists {
		if !f(k, element) {
			break
		}
	}
}

// DeleteBefore 删除指定下标前的元素
func (l *List[T, V]) DeleteBefore(index int) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.GetLen() == 0 {
		return
	}
	// 限定边界
	if index > l.GetLen() {
		index = l.GetLen()
	}
	l.lists = l.lists[index:]
}

// DeleteAfter 除指定下标后的元素（包括指定下标也会被删除）
func (l *List[T, V]) DeleteAfter(index int) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.GetLen() == 0 {
		return
	}
	// 限定边界
	if index > l.GetLen() {
		index = l.GetLen()
	}
	l.lists = l.lists[:index]
}

// push 压入元素
func (l *List[T, V]) push(element T) (index int) {
	// 根据分数找到对象下标
	index = l.getIndex(element)
	// 增加元素
	l.addList(index, element)
	return
}

// addList 增加元素
func (l *List[T, V]) addList(index int, element T) {
	// 裁剪切片
	before := general.ArrayCopy(l.lists[:index])
	after := general.ArrayCopy(l.lists[index:])
	// 重新组成切片
	l.lists = append(before, element)
	l.lists = append(l.lists, after...)
}

// getIndex 根据分数安排对象下标
func (l *List[T, V]) getIndex(element T) int {
	if l.GetLen() == 0 {
		return 0
	}
	elementScore := l.sortFunc(element)
	// 使用标准二分查找
	left, right := 0, l.GetLen()-1
	for left <= right {
		mid := left + (right-left)/2
		midScore := l.sortFunc(l.lists[mid])
		if midScore >= elementScore {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}
	return left
}
