package sort

import "github.com/kennycch/gotools/general"

/*
希尔排序
array：原数组
newArray：排序后数组
*/
func Shell[T general.Number](array []T, sortType SortType) (newArray []T) {
	newArray = make([]T, len(array))
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
				if (newArray[j] > newArray[j+step] && sortType == ASC) ||
					(newArray[j] < newArray[j+step] && sortType == DESC) {
					newArray[j], newArray[j+step] = newArray[j+step], newArray[j]
				} else {
					break
				}
			}
		}
	}
	return
}
