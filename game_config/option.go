package game_config

import "github.com/kennycch/gotools/general"

// SetDir 目录路径
func SetDir(dirPath string) option {
	return func(g *GameConfig) {
		g.dirPath = dirPath
	}
}

// SetFile 文件路径
func SetFile(filePath string) option {
	return func(g *GameConfig) {
		g.filePath = filePath
	}
}

// SetTarget 目标路径
func SetTarget(targetPath string) option {
	return func(g *GameConfig) {
		g.targetPath = targetPath
	}
}

// SetBlackList 文件黑名单
func SetBlackList(fileName ...string) option {
	return func(g *GameConfig) {
		g.blackList = fileName
	}
}

// SetPrimaryKey 设置主键
func SetPrimaryKey(keys ...Key) option {
	return func(g *GameConfig) {
		for _, key := range keys {
			g.primaryKeyMap[key.JsonFileName] = key.PrimaryKey
		}
	}
}

// SetGroup 设置分组
func SetGroup(groups ...Group) option {
	return func(g *GameConfig) {
		for _, group := range groups {
			// 过滤禁用键
			if general.InArray(disableKey, group.GroupId) {
				group.GroupId += "Dk"
			}
			if general.InArray(disableKey, group.GroupKey) {
				group.GroupKey += "Dk"
			}
			g.groupList[group.JsonFileName] = group
		}
	}
}

func SetKeyType(keyTypes ...KeyType) option {
	return func(g *GameConfig) {
		for _, keyType := range keyTypes {
			g.keyTypeMap[keyType.JsonFileName] = keyType
		}
	}
}
