package list

import (
	"fmt"
	"sync"

	"github.com/kennycch/gotools/general"
)

// NewList 创建列表对象
func NewList[T listType]() (listStruct *List[T]) {
	listStruct = &List[T]{
		lists: make([]T, 0),
		lock:  &sync.RWMutex{},
	}
	return
}

// Insert 新增对象
func (l *List[T]) Insert(element T) (index int) {
	l.lock.Lock()
	defer l.lock.Unlock()
	// 压入元素
	l.push(element)
	return index
}

// Remove 删除对象
func (l *List[T]) Remove(index int) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.lists = general.DeleteValueByKey(l.lists, index)
}

// GetElements 获取指定范围的对象
func (l *List[T]) GetElements(start, num int) (elements []T) {
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
func (l *List[T]) GetLen() int {
	return len(l.lists)
}

// PopFront 弹出链表首部对象
func (l *List[T]) PopFront() (element T, err error) {
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
func (l *List[T]) PopBack() (element T, err error) {
	if len(l.lists) == 0 {
		err = fmt.Errorf("list is empty")
		return
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	lastIndex := len(l.lists) - 1
	element = l.lists[lastIndex]
	l.lists = general.DeleteValueByKey(l.lists, lastIndex)
	return
}

// ClearList 清空链表
func (l *List[T]) ClearList() {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.lists = l.lists[:0]
}

// push 压入元素
func (l *List[T]) push(element T) (index int) {
	// 根据分数找到对象下标
	index = l.getIndex(element)
	// 增加元素
	l.addList(index, element)
	return
}

// addList 增加元素
func (l *List[T]) addList(index int, element T) {
	// 裁剪切片
	before := general.ArrayCopy(l.lists[:index])
	after := general.ArrayCopy(l.lists[index:])
	// 重新组成切片
	l.lists = append(before, element)
	l.lists = append(l.lists, after...)
}

// getIndex 根据分数安排对象下标
func (l *List[T]) getIndex(element T) int {
	index := -1
	if len(l.lists) != 0 {
		// 从中间开始对比，减少遍历数量
		start, end := 0, len(l.lists)
		for {
			if end-start <= 5 {
				break
			}
			// 取中间值
			middle := (end-start)/2 + start
			middleElement := l.lists[middle]
			if middleElement.GetScore() >= element.GetScore() {
				end = middle
			} else {
				start = middle
			}
		}
		// 遍历列表
		for i := start; i < end; i++ {
			oldElement := l.lists[i]
			if element.GetScore() <= oldElement.GetScore() {
				index = i
				break
			}
		}
		if index == -1 {
			index = end
		}
	} else {
		index = 0
	}
	return index
}
