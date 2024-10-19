package game_config

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/kennycch/gotools/general"
)

// NewGameConfig 配置转换工具对象
func NewGameConfig(options ...option) *GameConfig {
	gameConfig := &GameConfig{
		blackList:     make([]string, 0),
		groupList:     make(map[string]Group),
		primaryKeyMap: make(map[string]string),
		keyTypeMap:    make(map[string]KeyType),
	}
	for _, op := range options {
		op(gameConfig)
	}
	return gameConfig
}

// AddOptions 添加配置
func (g *GameConfig) AddOptions(options ...option) {
	for _, op := range options {
		op(g)
	}
}

// CreateStructByDir 生成目标目录的结构体
func (g *GameConfig) CreateStructByDir() []error {
	errs := []error{}
	if g.dirPath == "" || g.targetPath == "" {
		errs = append(errs, fmt.Errorf("dir path or target path can not empty"))
		return errs
	}
	// 读取目录
	files, err := os.ReadDir(g.dirPath)
	if err != nil {
		errs = append(errs, err)
		return errs
	}
	// 如果目标目录不存在，生成目录
	if _, err := os.ReadDir(g.targetPath); err != nil {
		if err = os.Mkdir(g.targetPath, 0777); err != nil {
			errs = append(errs, err)
			return errs
		}
	}
	// 开始处理
	for _, file := range files {
		// 过滤非json文件
		strArr := strings.Split(file.Name(), ".")
		if len(strArr) != 2 || strArr[1] != "json" {
			continue
		}
		// 不解析文件直接跳过
		if general.InArray(g.blackList, file.Name()) {
			continue
		}
		// 读取Json内容
		content, err := os.ReadFile(fmt.Sprintf("%s/%s", g.dirPath, file.Name()))
		if err != nil {
			errs = append(errs, err)
			continue
		}
		g.createGoFile(file.Name(), strArr[0], content)
	}
	// 生成管理文件
	g.createManagerFile()
	return errs
}

// CreateStructByFile 生成目标文件的结构体
func (g *GameConfig) CreateStructByFile() error {
	if g.filePath == "" || g.targetPath == "" {
		return fmt.Errorf("file path or target path can not empty")
	}
	// 如果目标目录不存在，生成目录
	if _, err := os.ReadDir(g.targetPath); err != nil {
		if err = os.Mkdir(g.targetPath, 0777); err != nil {
			return err
		}
	}
	// 过滤非json文件
	fileName := filepath.Base(g.filePath)
	strArr := strings.Split(fileName, ".")
	if len(strArr) != 2 || strArr[1] != "json" {
		return fmt.Errorf("file is not json")
	}
	// 不解析文件直接跳过
	if general.InArray(g.blackList, fileName) {
		return fmt.Errorf("file is in black list")
	}
	// 读取Json内容
	content, err := os.ReadFile(g.filePath)
	if err != nil {
		return err
	}
	g.createGoFile(fileName, strArr[0], content)
	// 生成管理文件
	g.createManagerFile()
	return nil
}

// createManagerFile 生成管理文件
func (g *GameConfig) createManagerFile() {
	fullName := filepath.Join(g.targetPath, "manager.go")
	// 文件内容
	content := fmt.Sprintf(managerTemplate, filepath.Base(g.targetPath))
	g.createFile(fullName, content)
}

// createGoFile 创建游戏配置文件
func (g *GameConfig) createGoFile(fileFullName, fileName string, content []byte) {
	// 解析Json
	js := &jsonStruct{
		FileName:    fileFullName,
		Upper:       general.HumpFormat(fileName, true),
		Lower:       general.HumpFormat(fileName, false),
		ExtraStruct: make([]string, 0),
		Content:     content,
		Keys:        []key{},
	}
	if keyTypes, ok := g.keyTypeMap[js.FileName]; ok {
		js.KeyTypeSets = keyTypes.KeyTypeSets
	}
	fileContent := g.analysisJson(js)
	if fileContent == "" {
		return
	}
	// 生成go文件内容
	fullName := filepath.Join(g.targetPath, fmt.Sprintf("%s.go", js.Lower))
	g.createFile(fullName, fileContent)
}

// createFile 生成文件
func (g *GameConfig) createFile(fullName, content string) {
	// 如果文件已存在，需要先删除文件
	_, err := os.Stat(fullName)
	if err == nil {
		os.Remove(fullName)
	}
	// 生成配置文件
	goFile, _ := os.OpenFile(fullName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	defer goFile.Close()
	goFile.WriteString(content)
}

// analysisJson 解析Json
func (g *GameConfig) analysisJson(js *jsonStruct) string {
	content := ""
	if len(js.Content) == 0 {
		return content
	}
	js.BaseStruct = fmt.Sprintf("type %s struct {\n", js.Upper)
	js.JsonStruct = fmt.Sprintf("type %sJson struct {\n", js.Upper)
	// 判断是数组还是对象
	firstStr := string(js.Content[0])
	switch firstStr {
	case "[": // 数组
		content = g.analysisArray(js)
	case "{": // 对象
		content = g.analysisObject(js)
	default:
		return content
	}
	return content
}

// analysisArray 数组分析
func (g *GameConfig) analysisArray(js *jsonStruct) string {
	array := []map[string]interface{}{}
	if err := json.Unmarshal(js.Content, &array); err != nil {
		return ""
	}
	// 获取主键
	primaryKey := g.getPrimaryKey(js)
	if len(array) > 0 {
		keys := []string{}
		for k := range array[0] {
			if k == general.HumpFormat(primaryKey, false) {
				continue
			}
			keys = append(keys, k)
		}
		sort.Strings(keys)
		keys = append([]string{general.HumpFormat(primaryKey, false)}, keys...)
		for _, k := range keys {
			v, ok := array[0][k]
			if !ok {
				continue
			}
			upper, lower := general.HumpFormat(k, true), general.HumpFormat(k, false)
			if general.InArray(disableKey, lower) {
				upper += "Dk"
				lower += "Dk"
			}
			jsKey := key{
				Upper: upper,
				Lower: lower,
			}
			// 指定类型
			specifiedType := js.getSpecifiedType(k, "")
			if specifiedType != "" {
				js.BaseStruct += fmt.Sprintf("	%s %s\n", lower, specifiedType)
				js.JsonStruct += fmt.Sprintf("	%s %s `json:\"%s\"`\n", upper, specifiedType, k)
				jsKey.KindType = specifiedType
			} else {
				valueType := reflect.TypeOf(v).Kind()
				switch valueType {
				case reflect.Bool: // 布尔
					js.BaseStruct += fmt.Sprintf("	%s bool\n", lower)
					js.JsonStruct += fmt.Sprintf("	%s bool `json:\"%s\"`\n", upper, k)
					jsKey.KindType = "bool"
				case reflect.Float64: // 数字（所有数字类型都会被识别为小数）
					if g.isIntValue(array, k) {
						js.BaseStruct += fmt.Sprintf("	%s int32\n", lower)
						js.JsonStruct += fmt.Sprintf("	%s int32 `json:\"%s\"`\n", upper, k)
						jsKey.KindType = "int32"
					} else {
						js.BaseStruct += fmt.Sprintf("	%s float64\n", lower)
						js.JsonStruct += fmt.Sprintf("	%s float64 `json:\"%s\"`\n", upper, k)
						jsKey.KindType = "float64"
					}
				case reflect.String: // 字符串
					js.BaseStruct += fmt.Sprintf("	%s string\n", lower)
					js.JsonStruct += fmt.Sprintf("	%s string `json:\"%s\"`\n", upper, k)
					jsKey.KindType = "string"
				case reflect.Array, reflect.Slice: // 数组
					valueType, jsonType := g.handleArray(array, js, k, v)
					if valueType == "" {
						continue
					}
					js.HasArray = true
					js.BaseStruct += fmt.Sprintf("	%s %s\n", lower, "[]"+valueType)
					js.JsonStruct += fmt.Sprintf("	%s %s `json:\"%s\"`\n", upper, "[]"+jsonType, k)
					jsKey.KindType = "[]" + valueType
					jsKey.BaseStruct = valueType
					jsKey.JsonStruct = jsonType
				}
			}
			js.Keys = append(js.Keys, jsKey)
		}
	}
	js.BaseStruct += "}\n"
	js.JsonStruct += "}\n"
	// 格式化结构体
	js.BaseStruct = baseStructFormat(js.BaseStruct)
	js.JsonStruct = jsonStructFormat(js.JsonStruct)
	// 复制方法
	copyStr := g.getCopy(js)
	// 引用部分
	pkg := ""
	if js.HasArray {
		pkg = "\n\n	\"github.com/kennycch/gotools/general\""
	}
	// 获取配置方法
	funcs := g.getFuncs(js)
	// 分组配置
	hasGroup, groupAnalysis := g.getGroupConfig(js)
	// 额外结构体
	exStructs := ""
	for _, exStruct := range js.ExtraStruct {
		exStructs += exStruct
	}
	return fmt.Sprintf(arrayTemplate,
		filepath.Base(g.targetPath), pkg, // 包名、引用部分
		js.Lower, js.Upper, js.Upper, js.Upper, js.Upper, js.Upper, js.Upper, // Config部分
		js.BaseStruct, js.JsonStruct, // 结构体部分
		js.Upper, js.Upper, js.Upper, js.Upper, js.FileName, js.Upper, hasGroup, // 注册、结构体名称、文件名称、是否分组
		js.Upper, js.Lower, js.Lower, js.Lower, // 获取配置
		js.Upper, js.Lower, js.Lower, js.Lower, // 获取组配置
		js.Upper, js.Lower, js.Lower, js.Lower, // 全部配置迭代器
		js.Upper, js.Lower, js.Lower, js.Lower, // 按类迭代
		js.Upper, js.Lower, js.Lower, js.Lower, js.Upper, js.Lower, js.Upper, js.Upper, js.Lower, groupAnalysis, // 解析Json
		js.Upper, primaryKey, // 主键
		copyStr,   // 复制
		funcs,     // 获取配置方法
		exStructs, // 额外结构体
	)
}

// analysisObject 对象分析
func (g *GameConfig) analysisObject(js *jsonStruct) string {
	mapping := map[string]interface{}{}
	if err := json.Unmarshal(js.Content, &mapping); err != nil {
		return ""
	}
	keys := []string{}
	for k := range mapping {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := mapping[k]
		if v == nil {
			continue
		}
		upper, lower := general.HumpFormat(k, true), general.HumpFormat(k, false)
		if general.InArray(disableKey, lower) {
			upper += "Dk"
			lower += "Dk"
		}
		jsKey := key{
			Upper: upper,
			Lower: lower,
		}
		// 指定类型
		specifiedType := js.getSpecifiedType(k, "")
		if specifiedType != "" {
			js.BaseStruct += fmt.Sprintf("	%s %s\n", lower, specifiedType)
			js.JsonStruct += fmt.Sprintf("	%s %s `json:\"%s\"`\n", upper, specifiedType, k)
			jsKey.KindType = specifiedType
		} else {
			valueType := reflect.TypeOf(v).Kind()
			switch valueType {
			case reflect.Bool: // 布尔
				js.BaseStruct += fmt.Sprintf("	%s bool\n", lower)
				js.JsonStruct += fmt.Sprintf("	%s bool `json:\"%s\"`\n", upper, k)
				jsKey.KindType = "bool"
			case reflect.Float64: // 数字（所有数字类型都会被识别为小数）
				num := v.(float64)
				if math.Floor(num) == num {
					js.BaseStruct += fmt.Sprintf("	%s int32\n", lower)
					js.JsonStruct += fmt.Sprintf("	%s int32 `json:\"%s\"`\n", upper, k)
					jsKey.KindType = "int32"
				} else {
					js.BaseStruct += fmt.Sprintf("	%s float64\n", lower)
					js.JsonStruct += fmt.Sprintf("	%s float64 `json:\"%s\"`\n", upper, k)
					jsKey.KindType = "float64"
				}
			case reflect.String: // 字符串
				js.BaseStruct += fmt.Sprintf("	%s string\n", lower)
				js.JsonStruct += fmt.Sprintf("	%s string `json:\"%s\"`\n", upper, k)
				jsKey.KindType = "string"
			case reflect.Array, reflect.Slice: // 数组
				valueType, jsonType := g.handleArray([]map[string]interface{}{mapping}, js, k, v)
				if valueType == "" {
					continue
				}
				js.HasArray = true
				js.BaseStruct += fmt.Sprintf("	%s %s\n", lower, "[]"+valueType)
				js.JsonStruct += fmt.Sprintf("	%s %s `json:\"%s\"`\n", upper, "[]"+jsonType, k)
				jsKey.KindType = "[]" + valueType
				jsKey.BaseStruct = valueType
				jsKey.JsonStruct = jsonType
			}
		}
		js.Keys = append(js.Keys, jsKey)
	}
	js.BaseStruct += "}\n"
	js.JsonStruct += "}\n"
	// 格式化结构体
	js.BaseStruct = baseStructFormat(js.BaseStruct)
	js.JsonStruct = jsonStructFormat(js.JsonStruct)
	// 复制方法
	copyStr := g.getCopy(js)
	// 引用部分
	pkg := ""
	if js.HasArray {
		pkg = "\n\n	\"github.com/kennycch/gotools/general\""
	}
	// 获取配置方法
	funcs := g.getFuncs(js)
	// 额外结构体
	exStructs := ""
	for _, exStruct := range js.ExtraStruct {
		exStructs += exStruct
	}
	return fmt.Sprintf(objectTemplate,
		filepath.Base(g.targetPath), pkg, // 包名、引用部分
		js.Lower, js.Upper, js.Upper, js.Upper, js.Upper, // Config部分
		js.BaseStruct, js.JsonStruct, // 结构体部分
		js.Upper, js.Upper, js.Upper, js.Upper, js.FileName, js.Upper, // 注册、结构体名称、文件名称、是否分组
		js.Upper, js.Lower, js.Lower, js.Lower, // 获取配置
		js.Upper,           // 获取组配置
		js.Upper, js.Upper, // 全部配置迭代器、按类迭代
		js.Upper, js.Lower, js.Lower, js.Upper, js.Lower, // 解析Json
		copyStr,   // 复制
		funcs,     // 获取配置方法
		exStructs, // 额外结构体
	)
}

// getCopy 生成复制内容
func (g *GameConfig) getCopy(js *jsonStruct) string {
	copyStr := fmt.Sprintf("func (cj %sJson) copy() %s {\n	c := %s{\n", js.Upper, js.Upper, js.Upper)
	arrayStructs := []key{}
	fields := ""
	for _, k := range js.Keys {
		if string(k.KindType[0]) == "[" {
			if general.InArray(normalArray, k.KindType) { // 一唯数组
				fields += fmt.Sprintf("		%s: general.ArrayCopy(cj.%s),\n", k.Lower, k.Upper)
			} else { // 数组结构体
				arrayStructs = append(arrayStructs, k)
			}
		} else {
			fields += fmt.Sprintf("		%s: cj.%s,\n", k.Lower, k.Upper)
		}
	}
	fields = baseStructFormat(fields)
	fields += "	}\n"
	// 数组结构体
	for _, k := range arrayStructs {
		fields += fmt.Sprintf("	%s := make([]%s, 0)\n	for _, ex := range cj.%s {\n		%s = append(%s, ex.copy())\n	}\n	c.%s = %s\n",
			k.Lower, js.Upper+k.Upper, k.Upper, k.Lower, k.Lower, k.Lower, k.Lower)
	}
	copyStr += fields
	copyStr += "	return c\n}"
	return copyStr
}

// getPrimaryKey 获取主键
func (g *GameConfig) getPrimaryKey(js *jsonStruct) string {
	primaryKey, ok := g.primaryKeyMap[js.FileName]
	if !ok {
		primaryKey = "Id"
	}
	return primaryKey
}

// handleArray 处理数组
func (g *GameConfig) handleArray(array []map[string]interface{}, js *jsonStruct, key string, value interface{}) (string, string) {
	valueType, jsonType := "", ""
	arr := value.([]interface{})
	if len(arr) > 0 {
		v := arr[0]
		kind := reflect.TypeOf(v).Kind()
		switch kind {
		case reflect.Bool: // 布尔
			valueType, jsonType = "bool", "bool"
		case reflect.Float64: // 数字（所有数字类型都会被识别为小数）
			if key == "id" || g.isIntArray(array, key) {
				valueType, jsonType = "int32", "int32"
			} else {
				valueType, jsonType = "float64", "float64"
			}
		case reflect.String: // 字符串
			valueType, jsonType = "string", "string"
		case reflect.Map: // 对象
			valueType, jsonType = g.handleExtraStruct(array, js, key, value)
		}
	}
	return valueType, jsonType
}

// handleExtraStruct 处理额外结构体
func (g *GameConfig) handleExtraStruct(array []map[string]interface{}, js *jsonStruct, k string, v interface{}) (string, string) {
	mappings := v.([]interface{})
	mapping := mappings[0].(map[string]interface{})
	// 获取结构体值方法
	name := general.HumpFormat(k, true)
	exJs := &jsonStruct{
		Upper: general.HumpFormat(js.Upper+name, true),
		Lower: general.HumpFormat(js.Upper+name, false),
	}
	valueType, jsonType := exJs.Upper, exJs.Upper+"Json"
	exJs.BaseStruct = fmt.Sprintf("type %s struct {\n", exJs.Upper)
	exJs.JsonStruct = fmt.Sprintf("type %sJson struct {\n", exJs.Upper)
	keys := []string{}
	for k := range mapping {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, ke := range keys {
		val := mapping[ke]
		if val == nil {
			continue
		}
		upper, lower := general.HumpFormat(ke, true), general.HumpFormat(ke, false)
		if general.InArray(disableKey, lower) {
			upper += "Dk"
			lower += "Dk"
		}
		jsKey := key{
			Upper: upper,
			Lower: lower,
		}
		// 指定类型
		specifiedType := js.getSpecifiedType(k, ke)
		if specifiedType != "" {
			exJs.BaseStruct += fmt.Sprintf("	%s %s\n", lower, specifiedType)
			exJs.JsonStruct += fmt.Sprintf("	%s %s `json:\"%s\"`\n", upper, specifiedType, ke)
			jsKey.KindType = specifiedType
		} else {
			vType := reflect.TypeOf(val).Kind()
			switch vType {
			case reflect.Bool: // 布尔
				exJs.BaseStruct += fmt.Sprintf("	%s bool\n", upper)
				exJs.JsonStruct += fmt.Sprintf("	%s bool `json:\"%s\"`\n", upper, ke)
				jsKey.KindType = "bool"
			case reflect.Float64: // 数字（所有数字类型都会被识别为小数）
				if g.isIntMap(array, k, ke) {
					exJs.BaseStruct += fmt.Sprintf("	%s int32\n", lower)
					exJs.JsonStruct += fmt.Sprintf("	%s int32 `json:\"%s\"`\n", upper, ke)
					jsKey.KindType = "int32"
				} else {
					exJs.BaseStruct += fmt.Sprintf("	%s float64\n", lower)
					exJs.JsonStruct += fmt.Sprintf("	%s float64 `json:\"%s\"`\n", upper, ke)
					jsKey.KindType = "float64"
				}
			case reflect.String: // 字符串
				exJs.BaseStruct += fmt.Sprintf("	%s string\n", lower)
				exJs.JsonStruct += fmt.Sprintf("	%s string `json:\"%s\"`\n", upper, ke)
				jsKey.KindType = "string"
			}
		}
		exJs.Keys = append(exJs.Keys, jsKey)
	}
	exJs.BaseStruct += "}\n"
	exJs.JsonStruct += "}\n"
	// 格式化结构体
	exJs.BaseStruct = baseStructFormat(exJs.BaseStruct)
	exJs.JsonStruct = jsonStructFormat(exJs.JsonStruct)
	// 生成获取值方法
	funcs := g.getFuncs(exJs)
	// 生成复制方法
	copyStr := g.getCopy(exJs)
	exStruct := fmt.Sprintf("\n\n%s\n%s\n%s%s", exJs.BaseStruct, exJs.JsonStruct, copyStr, funcs)
	js.ExtraStruct = append(js.ExtraStruct, exStruct)
	return valueType, jsonType
}

// getFuncs 生成获取值方法
func (g *GameConfig) getFuncs(js *jsonStruct) string {
	funcs := ""
	for _, key := range js.Keys {
		if string(key.KindType[0]) == "[" {
			funcs += fmt.Sprintf("\n\nfunc (c %s) %s() %s {\n	return general.ArrayCopy(c.%s)\n}", js.Upper, key.Upper, key.KindType, key.Lower)
		} else {
			funcs += fmt.Sprintf("\n\nfunc (c %s) %s() %s {\n	return c.%s\n}", js.Upper, key.Upper, key.KindType, key.Lower)
		}
	}
	return funcs
}

// getGroupConfig 生成分组解析方法
func (g *GameConfig) getGroupConfig(js *jsonStruct) (hasGroup, groupAnalysis string) {
	if groupConfig, ok := g.groupList[js.FileName]; ok {
		hasGroup = "true"
		template := `
		if _, ok := %sConfig.groupMap[conf.%s]; !ok {
			%sConfig.groupMap[conf.%s] = make(map[int32]%s)
		}
		%sConfig.groupMap[conf.%s][conf.%s] = conf`
		groupAnalysis = fmt.Sprintf(template,
			js.Lower,
			groupConfig.GroupId,
			js.Lower,
			groupConfig.GroupId,
			js.Upper,
			js.Lower,
			groupConfig.GroupId,
			groupConfig.GroupKey,
		)
	} else {
		hasGroup = "false"
	}
	return
}

// isIntValue 是否整型（简单类型）
func (g *GameConfig) isIntValue(array []map[string]interface{}, key string) bool {
	flag := true
	for _, value := range array {
		num := value[key].(float64)
		if math.Floor(num) != num {
			flag = false
			break
		}
	}
	return flag
}

// isIntArray 是否整型（数组）
func (g *GameConfig) isIntArray(array []map[string]interface{}, key string) bool {
	flag := true
label:
	for _, value := range array {
		arrNum := value[key].([]interface{})
		for _, v := range arrNum {
			num := v.(float64)
			if math.Floor(num) != num {
				flag = false
				break label
			}
		}
	}
	return flag
}

// isIntMap 是否整型（数组对象）
func (g *GameConfig) isIntMap(array []map[string]interface{}, key, ke string) bool {
	flag := true
label:
	for _, value := range array {
		if value[key] == nil {
			continue
		}
		mappings := value[key].([]interface{})
		for _, v := range mappings {
			mapping := v.(map[string]interface{})
			num := mapping[ke].(float64)
			if math.Floor(num) != num {
				flag = false
				break label
			}
		}
	}
	return flag
}

// getSpecifiedType 获取指定字段类型
func (js *jsonStruct) getSpecifiedType(baseKey, extraKey string) string {
	keyType := ""
	if len(js.KeyTypeSets) > 0 {
		for _, keyTypeSet := range js.KeyTypeSets {
			if baseKey == keyTypeSet.BaseKeyName && extraKey == keyTypeSet.ExtraKeyName {
				keyType = keyTypeSet.KeyType
				break
			}
		}
	}
	return keyType
}

// baseStructFormat 结构体格式化
func baseStructFormat(baseStruct string) string {
	arr := strings.Split(baseStruct, "\n")
	// 统计字段最大长度
	maxLen := 0
	for _, str := range arr {
		if len(str) == 0 || string(str[0]) != "\t" {
			continue
		}
		field := strings.Split(str, " ")
		if maxLen < len(field[0]) {
			maxLen = len(field[0])
		}
	}
	// 补全字段空格符
	for index, str := range arr {
		if len(str) == 0 || string(str[0]) != "\t" {
			continue
		}
		fields := strings.Split(str, " ")
		length := len(fields[0])
		if length < maxLen {
			for i := 0; i < maxLen-length; i++ {
				fields[0] += " "
			}
		}
		arr[index] = strings.Join(fields, " ")
	}
	return strings.Join(arr, "\n")
}

// jsonStructFormat json结构体格式化
func jsonStructFormat(jsonStruct string) string {
	arr := strings.Split(jsonStruct, "\n")
	// 统计字段最大长度
	maxLen := 0
	for _, str := range arr {
		if len(str) == 0 || string(str[0]) != "\t" {
			continue
		}
		field := strings.Split(str, " ")
		if maxLen < len(field[0]) {
			maxLen = len(field[0])
		}
	}
	// 补全字段空格符
	for index, str := range arr {
		if len(str) == 0 || string(str[0]) != "\t" {
			continue
		}
		fields := strings.Split(str, " ")
		length := len(fields[0])
		if length < maxLen {
			for i := 0; i < maxLen-length; i++ {
				fields[0] += " "
			}
		}
		arr[index] = strings.Join(fields, " ")
	}
	// 统计字段加类型最大长度
	maxLenJson := 0
	for _, str := range arr {
		if len(str) == 0 || string(str[0]) != "\t" {
			continue
		}
		foot := str[maxLen+1:]
		fields := strings.Split(foot, " ")
		if maxLenJson < len(fields[0]) {
			maxLenJson = len(fields[0])
		}
	}
	for index, str := range arr {
		if len(str) == 0 || string(str[0]) != "\t" {
			continue
		}
		head := str[:maxLen+1]
		foot := str[maxLen+1:]
		fields := strings.Split(foot, " ")
		length := len(fields[0])
		if length < maxLenJson {
			for i := 0; i < maxLenJson-length; i++ {
				fields[0] += " "
			}
		}
		arr[index] = head + strings.Join(fields, " ")
	}
	return strings.Join(arr, "\n")
}
