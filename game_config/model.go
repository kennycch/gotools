package game_config

var (
	// 一唯数组
	normalArray = []string{
		"[]bool",
		"[]int",
		"[]int32",
		"[]int64",
		"[]float64",
		"[]string",
	}
	// 禁用键
	disableKey = []string{
		"type",
		"select",
		"switch",
		"var",
		"func",
		"const",
		"range",
	}
)

type option func(g *GameConfig)

type GameConfig struct {
	dirPath       string
	filePath      string
	targetPath    string
	blackList     []string
	groupList     map[string]Group
	primaryKeyMap map[string]string
}

type jsonStruct struct {
	FileName    string
	Upper       string
	Lower       string
	BaseStruct  string
	JsonStruct  string
	HasArray    bool
	ExtraStruct []string
	Content     []byte
	Keys        []key
}

type key struct {
	Upper      string
	Lower      string
	KindType   string
	BaseStruct string
	JsonStruct string
}

// 默认主键
type Key struct {
	JsonFileName string
	PrimaryKey   string
}

type Group struct {
	JsonFileName string
	GroupId      string
	GroupKey     string
}
