package general

/*
Map深度拷贝
mapping：要拷贝的map
newMapping：拷贝出来的map
*/
func MapCopy[K Ordered, V any](mapping map[K]V) (newMapping map[K]V) {
	newMapping = make(map[K]V)
	for k, v := range mapping {
		newMapping[k] = v
	}
	return newMapping
}

/*
获取Map键
mapping：要拷贝的map
keys：map的键
*/
func MapKeys[K Ordered, V any](mapping map[K]V) (keys []K) {
	keys = make([]K, 0)
	for k, _ := range mapping {
		keys = append(keys, k)
	}
	return keys
}

/*
获取Map值
mapping：要拷贝的map
values：map的键
*/
func MapValues[K Ordered, V any](mapping map[K]V) (values []V) {
	values = make([]V, 0)
	for _, v := range mapping {
		values = append(values, v)
	}
	return values
}