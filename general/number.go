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
