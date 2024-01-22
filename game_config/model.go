package game_config

var (
	// 一唯数组
	normalArray = []string{
		"[]bool",
		"[]int",
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

type GameConfig struct {
	dirPath       string
	filePath      string
	targetPath    string
	blackList     []string
	primaryKeyMap map[string]string
}

type setOption interface {
	apply(*GameConfig)
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