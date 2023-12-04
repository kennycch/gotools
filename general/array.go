package general

/*
数组穷举
穷举数组拆分的所有可能（如：[]int{1,2,3}，拆分长度2，最终拆分结果为[][]int{{1,2},{1,3},{2,3}}）
array：需要拆分的数组
minSize：最少长度
maxSize：最大长度
result：拆分结果
*/
func LimitedCombinations[T any](array []T, minSize, maxSize uint) (result [][]T) {
	currentCombination := make([]T, 0)
	var backtrack func(startIndex, currentSize uint)
	n := len(array)
	if n == 0 || minSize <= 0 || maxSize <= 0 || minSize > maxSize {
		return result
	}
	backtrack = func(startIndex, currentSize uint) {
		if currentSize >= minSize && currentSize <= maxSize {
			combination := make([]T, len(currentCombination))
			copy(combination, currentCombination)
			result = append(result, combination)
		}
		for i := int(startIndex); i < n; i++ {
			currentCombination = append(currentCombination, array[i])
			backtrack(uint(i+1), currentSize+1)
			currentCombination = currentCombination[:len(currentCombination)-1]
		}
	}
	backtrack(0, 0)
	return
}
