package sort

import (
	"github.com/kennycch/gotools/general"
)

/*
计数排序
1、本函数仅适用于整形元素
2、函数运行速度取决于数组最大值与最小值的值差
3、本函数比较适用于数组中存在多个重复元素情况
array：原数组
newArray：排序后数组
*/
func Conting[T general.Integer](array []T, sortType SortType) (newArray []T) {
	newArray = make([]T, len(array))
	copy(newArray, array)
	// 长度少于2直接返回
	if len(newArray) < 2 {
		return
	}
	// 创建映射存放数组，并获取数组中的最小和最大值
	min, max := T(0), T(0)
	mapping := make(map[T]int, 0)
	// 将数组元素置于映射并计算次数
	for _, value := range newArray {
		mapping[value] += 1
		if value < min {
			min = value
		} else if value > max {
			max = value
		}
	}
	// 清空原来数组
	newArray = make([]T, 0)
	// 根据排序类型将映射放回数组
	if sortType == ASC { // 升序
		for i := min; i <= max; i++ {
			// 不存在的元素不作处理
			if num, ok := mapping[i]; ok {
				for j := 0; j < num; j++ {
					newArray = append(newArray, i)
				}
			}
		}
	} else { // 降序
		for i := max; i >= min; i-- {
			// 不存在的元素不作处理
			if num, ok := mapping[i]; ok {
				for j := 0; j < num; j++ {
					newArray = append(newArray, i)
				}
			}
		}
	}
	return
}
