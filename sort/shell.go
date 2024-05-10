package sort

import "github.com/kennycch/gotools/general"

// Shell 希尔排序
func Shell[V any, T general.Number](array []V, sortType SortType, sortValue func(value V) T) (newArray []V) {
	newArray = make([]V, len(array))
	copy(newArray, array)
	// 长度少于2直接返回
	if len(newArray) < 2 {
		return
	}
	// 基数为数组长度一半, 逐次递减基数，直至基数为1
	for step := len(newArray) / 2; step >= 1; step /= 2 {
		// 对比基数及基数倍数下标的元素
		for i := step; i < len(newArray); i += step {
			for j := i - step; j >= 0; j -= step {
				if (sortValue(newArray[j]) > sortValue(newArray[j+step]) && sortType == ASC) ||
					(sortValue(newArray[j]) < sortValue(newArray[j+step]) && sortType == DESC) {
					newArray[j], newArray[j+step] = newArray[j+step], newArray[j]
				} else {
					break
				}
			}
		}
	}
	return
}
