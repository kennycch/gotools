package general

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func init() {
	Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
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
