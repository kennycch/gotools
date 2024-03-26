package sort

import "github.com/kennycch/gotools/general"

/*
冒泡排序
array：原数组
newArray：排序后数组
*/
func Bubble[V any, T general.Number](array []V, sortType SortType, sortValue func(vlaue V) T) (newArray []V) {
	newArray = make([]V, len(array))
	copy(newArray, array)
	// 长度少于2直接返回
	if len(newArray) < 2 {
		return
	}
	// 开始两两对比
	for i := 0; i < len(newArray)-1; i++ {
		// 数组变更标识
		isChanged := false
		for j := 0; j < len(newArray)-1-i; j++ {
			// 根据排序类型判断对比条件
			if (sortValue(newArray[j]) > sortValue(newArray[j+1]) && sortType == ASC) ||
				(sortValue(newArray[j]) < sortValue(newArray[j+1]) && sortType == DESC) {
				newArray[j], newArray[j+1] = newArray[j+1], newArray[j]
				isChanged = true
			}
		}
		// 如果数组没有变更，终止后续对比
		if !isChanged {
			break
		}
	}
	return
}
