package general

import (
	"math"
	"math/rand"

	"github.com/shopspring/decimal"
)

/*
科学计算法
num1：前运算值
num2：后运算值
decType：运算类型
places：保留小数点后位数（不传默认4位，传多个仅首个生效）
result：运算结果
*/
func Decimal(num1, num2 float64, decType DecimalType, places ...uint) (result float64) {
	switch decType {
	case Add:
		result, _ = decimal.NewFromFloat(num1).Add(decimal.NewFromFloat(float64(num2))).Float64()
	case Sub:
		result, _ = decimal.NewFromFloat(num1).Sub(decimal.NewFromFloat(float64(num2))).Float64()
	case Multiply:
		result, _ = decimal.NewFromFloat(num1).Mul(decimal.NewFromFloat(float64(num2))).Float64()
	case Divide:
		result, _ = decimal.NewFromFloat(num1).Div(decimal.NewFromFloat(float64(num2))).Float64()
	}
	// 处理小数点后位数
	if len(places) == 0 {
		result = FloatPlaces(result, 4)
	} else {
		result = FloatPlaces(result, places[0])
	}
	return
}

/*
小数保留小数点后位数
num：原小数
places：保留小数点后位数
result：运算结果
*/
func FloatPlaces(num float64, places uint) (result float64) {
	result, _ = decimal.NewFromFloat(num).RoundFloor(int32(places)).Float64()
	return
}

/*
生成随机数，指定范围
mix：最小值
max：最大值
result：随机结果
*/
func RandomInt(mix int, max int) (result int) {
	if mix > max {
		return
	} else if mix == max {
		result = mix
		return
	}
	diff := max - mix + 1
	result = rand.Intn(diff) + mix
	return
}

/*
获取绝对值
num：要处理的值
result：处理结果
*/
func GetAbs[T Number](num T) (result T) {
	return T(math.Abs(float64(num)))
}

/*
数字合并
由两个int32整形合并成一个int64整形
high：高位数
low：低位数
mergedNumber：合并数
*/
func MergeNumbers(high int32, low int32) (mergedNumber int64) {
	mergedNumber = (int64(high) << 32) | (int64(low) & 0xFFFFFFFF)
	return
}

/*
数字分解
由一个int64整形分解成两个int32整形
mergedNumber：合并数
high：高位数
low：低位数
*/
func SplitNumbers(mergedNumber int64) (high, low int32) {
	high = int32(mergedNumber >> 32)
	low = int32(mergedNumber & 0xFFFFFFFF)
	return
}

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