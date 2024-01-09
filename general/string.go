package general

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

/*
SHA1加密（返回小写字母）
str：要加密的字符串
sha1Str：加密后的字符串
*/
func Sha1(str string) (sha1Str string) {
	sha1Byte := sha1.Sum([]byte(str))
	sha1Str = fmt.Sprintf("%x", sha1Byte)
	return
}

/*
MD5加密（返回32位小写）
str：要加密的字符串
md5Str：加密后的字符串
*/
func Md5(str string) (md5Str string) {
	md5Byte := md5.Sum([]byte(str))
	md5Str = fmt.Sprintf("%x", md5Byte)
	return
}

/*
MD5加密（返回16位小写）
str：要加密的字符串
short：加密后的字符串
*/
func Md5Short(str string) (short string) {
	short = Md5(str)[8:24]
	return
}

/*
随机字符串
length 字符串长度
randomStr 生成的字符串
*/
func RandomStr(length uint) (randomStr string) {
	charArr := strings.Split(CHAR, "")
	charlen := len(charArr)
	for i := 1; i <= int(length); i++ {
		randomStr += charArr[Rand.Intn(charlen)]
	}
	return
}

/*
获取唯一ID
id：唯一ID
*/
func GetUniqueId() (id string) {
	id = strings.ReplaceAll(uuid.New().String(), "-", "")
	return
}

/*
Json序列化（返回字符串）
obj：要序列化的数据
encode：序列化后的字符串
*/
func JsonEncode(obj interface{}) (encode string) {
	jsonS, _ := json.Marshal(obj)
	encode = string(jsonS)
	return
}
