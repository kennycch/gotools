package game_config

var (
	// 数组模板
	arrayTemplate = `package %s

import (
	"encoding/json"
	"sync"%s
)

var (
	%sConfig = %sManager{
		configMap: make(map[int32]%s),
		groupMap:  make(map[int32]map[int32]%s),
		lock:      &sync.RWMutex{},
	}
)

type %sManager struct {
	configMap map[int32]%s
	groupMap  map[int32]map[int32]%s
	lock      *sync.RWMutex
}

%s
%s
// 注册cl
func init() {
	AddCl(%s{})
}

// structName 结构体名称
func (c %s) structName() string {
	return "%s"
}

// fileName 文件名称
func (c %s) fileName() string {
	return "%s"
}

// hasGroup 是否分组
func (c %s) hasGroup() bool {
	return %s
}

// getConfigByKey 获取配置
func (c %s) getConfigByKey(id int32) (ICl, bool) {
	%sConfig.lock.RLock()
	defer %sConfig.lock.RUnlock()
	config, ok := %sConfig.configMap[id]
	return config, ok
}

// getConfigByGroup 获取组配置
func (c %s) getConfigByGroup(groupId int32, groupKey int32) (ICl, bool) {
	%sConfig.lock.RLock()
	defer %sConfig.lock.RUnlock()
	group, ok := %sConfig.groupMap[groupId]
	if !ok {
		return nil, ok
	}
	config, ok := group[groupKey]
	return config, ok
}

// iteratorConfigs 全部配置迭代器
func (c %s) iteratorConfigs(f func(key int32, value ICl) bool) {
	%sConfig.lock.RLock()
	defer %sConfig.lock.RUnlock()
	for k, v := range %sConfig.configMap {
		if !f(k, v) {
			break
		}
	}
}

// analysis 解析Json
func (c %s) analysis(content []byte) {
	%sConfig.lock.Lock()
	defer %sConfig.lock.Unlock()
	%sConfig.configMap = make(map[int32]%s)
	%sConfig.groupMap = make(map[int32]map[int32]%s)
	temp := []%sJson{}
	json.Unmarshal(content, &temp)
	for _, cj := range temp {
		conf := cj.copy()
		%sConfig.configMap[cj.getKey()] = conf%s
	}
}

func (cj %sJson) getKey() int32 {
	return cj.%s
}

%s%s%s
`

	// 对象模板
	objectTemplate = `package %s

import (
	"encoding/json"
	"sync"%s
)

var (
	%sConfig = %sManager{
		config: %s{},
		lock:   &sync.RWMutex{},
	}
)

type %sManager struct {
	config %s
	lock   *sync.RWMutex
}

%s
%s
// 注册cl
func init() {
	AddCl(%s{})
}

// structName 结构体名称
func (c %s) structName() string {
	return "%s"
}

// fileName 文件名称
func (c %s) fileName() string {
	return "%s"
}

// hasGroup 是否分组
func (c %s) hasGroup() bool {
	return false
}

// getConfigByKey 获取配置
func (c %s) getConfigByKey(id int32) (ICl, bool) {
	%sConfig.lock.RLock()
	defer %sConfig.lock.RUnlock()
	return %sConfig.config, true
}

// getConfigByGroup 获取组配置（对象配置中不会进行任何操作）
func (c %s) getConfigByGroup(groupId int32, groupKey int32) (ICl, bool) {
	return nil, false
}

// iteratorConfigs 全部配置迭代器（对象配置中不会进行任何操作）
func (c %s) iteratorConfigs(f func(key int32, value ICl) bool) {

}

// analysis 解析Json
func (c %s) analysis(content []byte) {
	%sConfig.lock.Lock()
	defer %sConfig.lock.Unlock()
	temp := %sJson{}
	json.Unmarshal(content, &temp)
	%sConfig.config = temp.copy()
}

%s%s%s
`

	// 管理文件模板
	managerTemplate = `package %s

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

var (
	// 临时json内容
	loadTmpJsonMap = map[string][]byte{}
	// json文件修改时间
	fileChangeTime = map[string]int64{}
	// 结构体映射
	fileNameToCL = make(map[string]ICl)
)

type ICl interface {
	fileName() string
	structName() string
	analysis([]byte)
	getConfigByKey(int32) (ICl, bool)
	iteratorConfigs(f func(key int32, value ICl) bool)
	hasGroup() bool
	getConfigByGroup(int32, int32) (ICl, bool)
}

func AddCl(cl ICl) {
	fileNameToCL[cl.fileName()] = cl
}

// InitCl 开始加载配置
func InitCl(dirPath string) {
	// 读取json文件
	filepath.WalkDir(dirPath, func(fileDir string, file fs.DirEntry, err error) error {
		readFileLoadMap(file, fileDir)
		return err
	})
	// 解析Json
	for fileName, icl := range fileNameToCL {
		if content, ok := loadTmpJsonMap[fileName]; ok {
			icl.analysis(content)
		} else {
			panic(fmt.Errorf("config file not found, file name:%%s", icl.fileName()))
		}
	}
	// 清除缓存
	clearTemp()
}

// readFileLoadMap 读取配置文件
func readFileLoadMap(file fs.DirEntry, fileDir string) {
	content, err := os.ReadFile(fileDir)
	if err != nil {
		return
	}
	info, err := file.Info()
	if err != nil {
		return
	}
	loadTmpJsonMap[file.Name()] = content
	fileChangeTime[file.Name()] = info.ModTime().Unix()
}

// GetGameConfig 获取单个配置
func GetGameConfig[T ICl](cl T, id int32) (T, bool) {
	icl, ok := cl.getConfigByKey(id)
	if ok {
		cl = icl.(T)
	}
	return cl, ok
}

// GetGameConfigByGroup 按分组获取单个配置
func GetGameConfigByGroup[T ICl](cl T, groupId, groupKey int32) (T, bool) {
	if !cl.hasGroup() {
		return cl, false
	}
	icl, ok := cl.getConfigByGroup(groupId, groupKey)
	if ok {
		cl = icl.(T)
	}
	return cl, ok
}

// IteratorAllConfig 全部配置迭代器
func IteratorAllConfig[T ICl](cl T, f func(key int32, value ICl) bool) {
	cl.iteratorConfigs(f)
}

// clearTemp 清除缓存
func clearTemp() {
	loadTmpJsonMap = make(map[string][]byte, 0)
}
`
)
