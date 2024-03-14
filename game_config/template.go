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
		lock:      &sync.RWMutex{},
	}
)

type %sManager struct {
	configMap map[int32]%s
	lock      *sync.RWMutex
}

%s
%s
// 注册cl
func init() {
	AddCl(%s{})
}

// 结构体名称
func (c %s) StructName() string {
	return "%s"
}

// 文件名称
func (c %s) FileName() string {
	return "%s"
}

// 获取配置
func (c %s) GetConfigByKey(id int32) (ICl, bool) {
	%sConfig.lock.RLock()
	defer %sConfig.lock.RUnlock()
	config, ok := %sConfig.configMap[id]
	return config, ok
}

// 全部配置迭代器
func (c %s) IteratorConfigs(f func(key int32, value ICl) bool) {
	%sConfig.lock.RLock()
	defer %sConfig.lock.RUnlock()
	for k, v := range %sConfig.configMap {
		if !f(k, v) {
			break
		}
	}
}

// 解析Json
func (c %s) Analysis(content []byte) {
	%sConfig.lock.Lock()
	defer %sConfig.lock.Unlock()
	%sConfig.configMap = make(map[int32]%s)
	temp := []%sJson{}
	json.Unmarshal(content, &temp)
	for _, cj := range temp {
		%sConfig.configMap[cj.getKey()] = cj.copy()
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

// 结构体名称
func (c %s) StructName() string {
	return "%s"
}

// 文件名称
func (c %s) FileName() string {
	return "%s"
}

// 获取配置
func (c %s) GetConfigByKey(id int32) (ICl, bool) {
	%sConfig.lock.RLock()
	defer %sConfig.lock.RUnlock()
	return %sConfig.config, true
}

// 全部配置迭代器（对象配置中不会进行任何操作）
func (c %s) IteratorConfigs(f func(key int32, value ICl) bool) {
	
}

// 解析Json
func (c %s) Analysis(content []byte) {
	%sConfig.lock.Lock()
	defer %sConfig.lock.Unlock()
	temp := %sJson{}
	json.Unmarshal(content, &temp)
	%sConfig.config = temp.copy()
}

%s%s%s
`

	managerTemplate = `package %s

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/kennycch/gotools/timer"
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
	FileName() string
	StructName() string
	Analysis([]byte)
	GetConfigByKey(int32) (ICl, bool)
	IteratorConfigs(f func(key int32, value ICl) bool)
}

func AddCl(cl ICl) {
	fileNameToCL[cl.FileName()] = cl
}

// 开始加载配置
func InitCl(dirPath string) {
	// 读取json文件
	filepath.WalkDir(dirPath, func(fileDir string, file fs.DirEntry, err error) error {
		readFileLoadMap(file, fileDir)
		return err
	})
	// 解析Json
	for fileName, content := range loadTmpJsonMap {
		if icl, ok := fileNameToCL[fileName]; ok {
			icl.Analysis(content)
		}
	}
	// 监听配置更改
	Listen(dirPath)
	// 清除缓存
	ClearTemp()
}

// 读取配置文件
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

// 获取单个配置
func GetGameConfig[T ICl](cl T, id int32) (T, bool) {
	icl, ok := cl.GetConfigByKey(id)
	if ok {
		cl = icl.(T)
	} else {
		fmt.Println("config not found", cl.StructName(), id)
	}
	return cl, ok
}

// 全部配置迭代器
func IteratorAllConfig[T ICl](cl T, f func(key int32, value ICl) bool) {
	cl.IteratorConfigs(f)
}

// 监听配置更改
func Listen(dirPath string) {
	timer.AddTicker(5*time.Minute, func() {
		// 重载配置
		filepath.WalkDir(".", func(dirPath string, file fs.DirEntry, err error) error {
			if _, ok := fileNameToCL[file.Name()]; ok {
				ReloadConfig(file, dirPath)
			}
			return err
		})
	})
}

// 清除缓存
func ClearTemp() {
	loadTmpJsonMap = make(map[string][]byte, 0)
}

// 重载配置
func ReloadConfig(file fs.DirEntry, fileDir string) {
	info, err := file.Info()
	if err != nil {
		return
	}
	if changeTime, ok := fileChangeTime[file.Name()]; ok && info.ModTime().Unix() == changeTime {
		return
	}
	readFileLoadMap(file, fileDir)
	if icl, ok := fileNameToCL[file.Name()]; ok {
		icl.Analysis(loadTmpJsonMap[file.Name()])
	}
}
	
`
)
