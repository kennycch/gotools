package sort

import "github.com/kennycch/gotools/general"

// Select 选择排序
func Select[V any, T general.Number](array []V, sortType SortType, sortValue func(value V) T) (newArray []V) {
	newArray = make([]V, len(array))
	copy(newArray, array)
	// 长度少于2直接返回
	if len(newArray) < 2 {
		return
	}
	// 开始两两对比
	for i := 0; i < len(newArray)-1; i++ {
		minIndex := i
		// 找出最后面比自己小/大的元素下标，与自己对换
		for j := i + 1; j < len(newArray); j++ {
			if (sortValue(newArray[minIndex]) > sortValue(newArray[j]) && sortType == ASC) ||
				(sortValue(newArray[minIndex]) < sortValue(newArray[j]) && sortType == DESC) {
				minIndex = j
			}
		}
		newArray[i], newArray[minIndex] = newArray[minIndex], newArray[i]
	}
	return
}
