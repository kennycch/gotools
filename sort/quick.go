package sort

import "gotools/general"

/*
快速排序
array：原数组
newArray：排序后数组
*/
func Quick[T general.Number](array []T, sortType SortType) (newArray []T) {
	newArray = make([]T, len(array))
	copy(newArray, array)
	// 长度少于2直接返回
	if len(newArray) < 2 {
		return
	}
	// 调用递归函数
	newArray = quickAction[T](newArray, sortType, 0, len(newArray)-1)
	return
}

// 递归函数
func quickAction[T general.Number](array []T, sortType SortType, left, right int) []T {
	l := left
	r := right
	// 获取中轴
	pivot := array[(l+r)/2]
	// 遍历数组，将比中轴小/大的放于左侧，反则放于右侧
	for l < r {
		// 找到中轴两侧不符合条件的元素
		if sortType == ASC { // 升序
			for array[l] < pivot {
				l++
			}
			for array[r] > pivot {
				r--
			}
		} else { // 降序
			for array[l] > pivot {
				l++
			}
			for array[r] < pivot {
				r--
			}
		}
		// 如果左下标比右下标大，证明已完成本次任务
		if l >= r {
			break
		}
		// 对调元素
		array[l], array[r] = array[r], array[l]
		// 优化算法，减少重复遍历
		if array[l] == pivot {
			r--
		}
		if array[r] == pivot {
			l++
		}
	}
	// 如果 l == r, 再移动一下
	if l == r {
		l++
		r--
	}
	// 开始左递归
	if left < r {
		array = quickAction[T](array, sortType, left, r)
	}
	// 开始右递归
	if right > l {
		array = quickAction[T](array, sortType, l, right)
	}
	return array
}
