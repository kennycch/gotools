package general

import "math/rand"

// 科学计算法计算类型
const (
	Add DecimalType = iota
	Sub
	Multiply
	Divide
)

type DecimalType int

// 随机字符串数据源
const (
	CHAR = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// Ordered 可排序类型
type Ordered interface {
	Number | ~string
}

// Number 数字类型
type Number interface {
	Integer | Float
}

// SignedNumber 有符号数字类型
type SignedNumber interface {
	Signed | Float
}

// Integer 整数类型
type Integer interface {
	Signed | Unsigned
}

// Signed 有符号整数类型
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned 无符号整数类型
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Float 浮点类型
type Float interface {
	~float32 | ~float64
}

var (
	// 随机种子
	Rand *rand.Rand
)
