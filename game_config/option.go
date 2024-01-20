package game_config

// 目录路径
type Dir struct {
	dirPath string
}

func (d *Dir) apply(gameConfig *GameConfig) {
	gameConfig.dirPath = d.dirPath
}

func SetDir(dirPath string) *Dir {
	dir := &Dir{
		dirPath: dirPath,
	}
	return dir
}

// 文件路径
type File struct {
	filePath string
}

func (f *File) apply(gameConfig *GameConfig) {
	gameConfig.filePath = f.filePath
}

func SetFile(filePath string) *File {
	file := &File{
		filePath: filePath,
	}
	return file
}

// 目标路径
type Target struct {
	targetPath string
}

func (t *Target) apply(gameConfig *GameConfig) {
	gameConfig.targetPath = t.targetPath
}

func SetTarget(targetPath string) *Target {
	target := &Target{
		targetPath: targetPath,
	}
	return target
}

// 文件黑名单
type BlackList struct {
	blackList []string
}

func (b *BlackList) apply(gameConfig *GameConfig) {
	gameConfig.blackList = b.blackList
}

func SetBlackList(fileName ...string) *BlackList {
	black := &BlackList{
		blackList: fileName,
	}
	return black
}

// 默认主键
type PrimaryKey struct {
	primaryKeyMap map[string]string
}

type Key struct {
	JsonFileName string
	PrimaryKey   string
}

func (p *PrimaryKey) apply(gameConfig *GameConfig) {
	gameConfig.primaryKeyMap = p.primaryKeyMap
}

func SetPrimaryKey(keys ...Key) *PrimaryKey {
	primaryKey := &PrimaryKey{
		primaryKeyMap: make(map[string]string),
	}
	for _, key := range keys {
		primaryKey.primaryKeyMap[key.JsonFileName] = key.PrimaryKey
	}
	return primaryKey
}
