package sort

import "github.com/kennycch/gotools/general"

// Quick 快速排序
func Quick[V any, T general.Number](array []V, sortType SortType, sortValue func(value V) T) (newArray []V) {
	newArray = make([]V, len(array))
	copy(newArray, array)
	// 长度少于2直接返回
	if len(newArray) < 2 {
		return
	}
	// 调用递归函数
	newArray = quickAction[V, T](newArray, sortType, 0, len(newArray)-1, sortValue)
	return
}

func quickAction[V any, T general.Number](array []V, sortType SortType, left, right int, sortValue func(value V) T) []V {
	l := left
	r := right
	// 获取中轴
	pivot := array[(l+r)/2]
	// 遍历数组，将比中轴小/大的放于左侧，反则放于右侧
	for l < r {
		// 找到中轴两侧不符合条件的元素
		if sortType == ASC { // 升序
			for sortValue(array[l]) < sortValue(pivot) {
				l++
			}
			for sortValue(array[r]) > sortValue(pivot) {
				r--
			}
		} else { // 降序
			for sortValue(array[l]) > sortValue(pivot) {
				l++
			}
			for sortValue(array[r]) < sortValue(pivot) {
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
		if sortValue(array[l]) == sortValue(pivot) {
			r--
		}
		if sortValue(array[r]) == sortValue(pivot) {
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
		array = quickAction[V, T](array, sortType, left, r, sortValue)
	}
	// 开始右递归
	if right > l {
		array = quickAction[V, T](array, sortType, l, right, sortValue)
	}
	return array
}
