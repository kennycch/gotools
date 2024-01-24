package a_star

import "github.com/kennycch/gotools/list"

const (
	// 直线节点G值
	STRAIGHT = 10
	// 斜线节点G值
	DIAGONAL = 14
)

// 路径导航
type AStar struct {
	mapRow         int               // 地图列数（x轴）
	maxNode        int               // 地图最大坐标
	diagonal       bool              // 是否可以对角线移动
	closeList      map[int]*Node     // 已探索节点
	openList       *list.List[*Node] // 作为探索中心的节点
	obstacles      map[int]struct{}  // 障碍物
	stratIndex     int               // 开始节点
	currIndex      int               // 当前节点
	targetIndex    int               // 目标节点
	isSearchTarget bool              // 是否已找到目标节点
}

// 节点
type Node struct {
	historyCost  float64 // 历史代价（G值）
	estimateCost float64 // 未来预计代价（H值）
	totalCost    float64 // 总代价（F值）
	index        int     // 节点坐标
	parentIndex  int     // 上一个节点
}
