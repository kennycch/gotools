package general

import (
	"strconv"
	"strings"
)

/*
数字拆分法
给定一个数字,拆分成若干份,限制每份长度范围,求余数最少拆分份数最少结果
若最佳方案依然包含余数，返回数组最后一个元素是余数
num：要拆分的数字
min：每份最少长度
max: 每份最大长度
results：算法结果
*/
func SplitNumberWithRange(num, min, max int) (results []int) {
	if num <= 0 || min <= 0 || max <= 0 || min > max {
		return
	}
	// 如果给定数字少于最少长度，直接返回
	if num < min {
		results = append(results, num)
		return
	}
	// 列举符合范围的基数
	rangeNums := make([]int, 0)
	for i := min; i <= max; i++ {
		rangeNums = append(rangeNums, i)
	}
	// 穷举所有基数组合，取余数最少，结果最短的返回
	minRemainder := num
	minLenght := num
	for _, rangeNum := range rangeNums {
		// 如果基数比要拆分的数字大，直接终止穷举
		if rangeNum > num {
			break
		}
		// 整除向下取整，生成数组
		arrayLenght := num / rangeNum
		temp := make([]int, 0)
		for i := 0; i < arrayLenght; i++ {
			temp = append(temp, rangeNum)
		}
		// 将余数均到每个元素中，直至余数为0或元素超出范围
		remainder := num % rangeNum
		index := 0
		for remainder > 0 {
			if index >= len(temp)-1 {
				index = 0
			}
			// 如果元素+1超过最大长度，停止均摊并将余数置于数组末端
			if temp[index]+1 > max {
				temp = append(temp, remainder)
				// 如果余数在长度范围内，最后一个元素不算余数，清空余数
				if remainder >= min {
					remainder = 0
				}
				break
			}
			// 将余数均到每个元素上
			temp[index] += 1
			remainder -= 1
			index += 1
		}
		// 如果有余数，尝试将前面元素均到余数上，消除余数
		if remainder > 0 {
			// 计算最大可均数
			maxReduce := 0
			for _, i := range temp {
				// 符合范围的加到可均数中
				if i >= min && i <= max {
					maxReduce += i - min
				}
			}
			// 如果余数加上最大可均数达到每份最少长度，进行处理
			if remainder+maxReduce >= min {
				// 消除余数
				remainder = 0
				// 从后往前均
				k := len(temp) - 2
				for temp[len(temp)-1] < min {
					if k < 0 {
						k = len(temp) - 2
					}
					temp[k] -= 1
					temp[len(temp)-1] += 1
					k -= 1
				}
			}
		}
		// 如果数组结果最短，且余数最少，现有方案暂定为最佳方案
		if len(temp) <= minLenght && remainder <= minRemainder {
			results = temp
			minLenght = len(temp)
			minRemainder = remainder
		}
	}
	return
}

/*
版本对比
返回：-1：现在版本新，0：版本一样，1：要比较的本版新
nowVersion：现在的版本
newVersion：要比较的版本
result：对比结果
err：错误信息
*/
func VersionComparison(nowVersion, newVersion string) (result int, err error) {
	// 分隔字符串并放入映射
	nowString := strings.Split(nowVersion, ".")
	nowMap := make(map[int]int, 0)
	for k, str := range nowString {
		num := 0
		num, err = strconv.Atoi(str)
		if err != nil {
			return
		}
		nowMap[k] = num
	}
	newString := strings.Split(newVersion, ".")
	newMap := make(map[int]int, 0)
	for k, str := range newString {
		num := 0
		num, err = strconv.Atoi(str)
		if err != nil {
			return
		}
		newMap[k] = num
	}
	// 选长度较长的作为循环次数
	maxLenght := len(nowMap)
	if len(newMap) > maxLenght {
		maxLenght = len(newMap)
	}
	for i := 0; i < maxLenght; i++ {
		nowNum, ok := nowMap[i]
		newNum, ok1 := newMap[i]
		if ok && ok1 {
			// 对比数值
			if nowNum > newNum {
				result = -1
				break
			} else if nowNum < newNum {
				result = 1
				break
			}
		} else if !ok {
			result = 1
			break
		} else {
			result = -1
			break
		}
	}
	err = nil
	return
}
