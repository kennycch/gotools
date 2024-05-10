package sort

import "github.com/kennycch/gotools/general"

// Merge 归并排序
func Merge[V any, T general.Number](array []V, sortType SortType, sortValue func(value V) T) (newArray []V) {
	newArray = make([]V, len(array))
	copy(newArray, array)
	// 长度少于2直接返回
	if len(newArray) < 2 {
		return
	}
	m := &mergeSort[V, T]{
		array:    newArray,
		sortType: sortType,
	}
	m.splitAction(0, len(newArray)-1)
	newArray = m.array
	return
}

type mergeSort[V any, T general.Number] struct {
	array     []V      // 要排序的数组
	sortType  SortType // 排序类型
	sortValue func(V) T
}

func (m *mergeSort[V, T]) splitAction(start, end int) {
	if start >= end {
		return
	}
	// 取中间值
	middle := (start + end) / 2
	m.splitAction(start, middle)
	m.splitAction(middle+1, end)
	// 合并处理
	m.mergeAcion(start, middle, end)
}

func (m *mergeSort[V, T]) mergeAcion(start, middle, end int) {
	// 创建下标与临时数组
	leftIndex, rightIndex, tempIndex := start, middle+1, 0
	temp := make([]V, 1+end-start)
	// 取左右下标中较小/大的放入临时数组
	for leftIndex <= middle && rightIndex <= end {
		if (m.sortValue(m.array[leftIndex]) <= m.sortValue(m.array[rightIndex]) && m.sortType == ASC) ||
			(m.sortValue(m.array[leftIndex]) >= m.sortValue(m.array[rightIndex]) && m.sortType == DESC) {
			temp[tempIndex] = m.array[leftIndex]
			leftIndex++
			tempIndex++
		} else {
			temp[tempIndex] = m.array[rightIndex]
			rightIndex++
			tempIndex++
		}
	}
	// 判断哪边序列还有剩余元素
	appendStart, appendEnd := 0, 0
	if leftIndex > middle {
		appendStart = rightIndex
		appendEnd = end
	} else {
		appendStart = leftIndex
		appendEnd = middle
	}
	// 将剩余元素放入临时数组
	for appendStart <= appendEnd {
		temp[tempIndex] = m.array[appendStart]
		tempIndex++
		appendStart++
	}
	// 将临时数组放入原数组
	for k, v := range temp {
		m.array[start+k] = v
	}
}
