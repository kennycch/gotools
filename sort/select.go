package sort

import "github.com/kennycch/gotools/general"

/*
选择排序
array：原数组
newArray：排序后数组
*/
func Select[T general.Number](array []T, sortType SortType) (newArray []T) {
	newArray = make([]T, len(array))
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
			if (newArray[minIndex] > newArray[j] && sortType == ASC) ||
				(newArray[minIndex] < newArray[j] && sortType == DESC) {
				minIndex = j
			}
		}
		newArray[i], newArray[minIndex] = newArray[minIndex], newArray[i]
	}
	return
}
