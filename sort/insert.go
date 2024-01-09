package sort

import "github.com/kennycch/gotools/general"

/*
插入排序
array：原数组
newArray：排序后数组
*/
func Insert[T general.Number](array []T, sortType SortType) (newArray []T) {
	newArray = make([]T, len(array))
	copy(newArray, array)
	// 长度少于2直接返回
	if len(newArray) < 2 {
		return
	}
	// 创建有序数组
	temp := make([]T, 0)
	// 开始插入元素
	for i := 0; i < len(newArray); i++ {
		// 如果是第一个元素，直接插入
		if i == 0 {
			temp = append(temp, newArray[i])
			continue
		}
		// 寻找比自己小/大的元素下标
		target := -1
		for j := 0; j < len(temp); j++ {
			if (newArray[i] > temp[j] && sortType == ASC) ||
				(newArray[i] < temp[j] && sortType == DESC) {
				target = j
			}
		}
		// 插入元素
		if target == -1 { // 如果没有比自己小/大元素，直接插入最前
			temp = append([]T{newArray[i]}, temp...)
		} else { // 插入目标元素后面
			newTemp := make([]T, 0)
			for k, v := range temp {
				newTemp = append(newTemp, v)
				if k == target {
					newTemp = append(newTemp, array[i])
				}
			}
			temp = newTemp
		}
	}
	newArray = temp
	return
}
