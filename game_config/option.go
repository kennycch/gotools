package game_config

// 目录路径
func SetDir(dirPath string) option {
	return func(g *GameConfig) {
		g.dirPath = dirPath
	}
}

// 文件路径
func SetFile(filePath string) option {
	return func(g *GameConfig) {
		g.filePath = filePath
	}
}

// 目标路径
func SetTarget(targetPath string) option {
	return func(g *GameConfig) {
		g.targetPath = targetPath
	}
}

// 文件黑名单
func SetBlackList(fileName ...string) option {
	return func(g *GameConfig) {
		g.blackList = fileName
	}
}

// 默认主键
type Key struct {
	JsonFileName string
	PrimaryKey   string
}

func SetPrimaryKey(keys ...Key) option {
	return func(g *GameConfig) {
		for _, key := range keys {
			g.primaryKeyMap[key.JsonFileName] = key.PrimaryKey
		}
	}
}
