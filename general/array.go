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

/*
数组转映射
array：要转换的数组
mapping：转换后的映射
*/
func ArrayToMapping[T Ordered](array []T) (mapping map[T]bool) {
	mapping = make(map[T]bool, 0)
	for _, value := range array {
		mapping[value] = false
	}
	return
}

/*
随机打乱数组
array：原始数组
result：打乱后的数组
*/
func Shuffle[T any](array []T) (result []T) {
	result = make([]T, len(array))
	copy(result, array)
	if len(result) <= 1 {
		return result
	}
	Rand.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})
	return
}

/*
是否在数组里
array：要搜索的数组
value：要对比的元素
flag：对比结果
*/
func InArray[T Ordered](array []T, value T) (flag bool) {
	for _, v := range array {
		if v == value {
			flag = true
			break
		}
	}
	return
}

/*
获取合集
注：重复元素会被自动去重
arrayA：参照数组
arrayB：对比数组
sameArray：数组合集
*/
func SameArray[T Ordered](arrayA []T, arrayB []T) (sameArray []T) {
	sameArray = make([]T, 0)
	mappingA := ArrayToMapping[T](arrayA)
	mappingB := ArrayToMapping[T](arrayB)
	for key := range mappingA {
		if _, ok := mappingB[key]; ok {
			sameArray = append(sameArray, key)
		}
	}

	return
}

/*
获取差集
注：返回结果为存在于参照数组且不存在于对比数组的元素，重复元素会被自动去重
arrayA：参照数组
arrayB：对比数组
diffArray：数组差集
*/
func DiffArray[T Ordered](arrayA []T, arrayB []T) (diffArray []T) {
	diffArray = make([]T, 0)
	mappingA := ArrayToMapping[T](arrayA)
	mappingB := ArrayToMapping[T](arrayB)
	for key := range mappingA {
		if _, ok := mappingB[key]; !ok {
			diffArray = append(diffArray, key)
		}
	}
	return
}

/*
数组去重
array：要去重的数组
uniqueArray：去重后结果
*/
func UniqueArray[T Ordered](array []T) (uniqueArray []T) {
	uniqueArray = make([]T, 0)
	mapping := ArrayToMapping[T](array)
	for key := range mapping {
		uniqueArray = append(uniqueArray, key)
	}
	return
}

/*
删除数组指定下标元素
array：要删除元素的数组
deletionArray：删除元素后的数组
*/
func DeleteValueByKey[T any](array []T, key int) (deletionArray []T) {
	// 创建临时数组
	deletionArray = ArrayCopy(array)
	// 如果下标不符合要求，直接返回
	if key < 0 || key >= len(array) {
		return
	}
	deletionArray = append(deletionArray[:key], deletionArray[key+1:]...)
	return
}

/*
删除数组指定值元素
array：要处理的数组
value：要删除的元素值
repeat：是否重复删除（如果否，删除一次该值后将直接返回）
*/
func DeleteValueByValue[T Ordered](array []T, value T, repeat ...bool) (deletionArray []T) {
	flag := false
	// 创建临时数组
	deletionArray = make([]T, 0)
	// 将无须删除的放入临时数组
	for _, v := range array {
		if v != value || (flag && (len(repeat) == 0 || !repeat[0])) {
			deletionArray = append(deletionArray, v)
		} else {
			flag = true
		}
	}
	return
}

/*
转数组元素类型
F：旧数组元素类型
T：新数组元素类型
oldArray：要转换的数组
newArray：转换后的数组
*/

func ArrayChangeType[F Number, T Number](oldArray []F) (newArray []T) {
	newArray = make([]T, 0)
	for _, value := range oldArray {
		newArray = append(newArray, T(value))
	}
	return
}

/*
复制数组
array：要复制的数组
newArray：复制后的数组
*/
func ArrayCopy[T any](array []T) (newArray []T) {
	newArray = make([]T, len(array))
	copy(newArray, array)
	return
}
