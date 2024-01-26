package jps

import "github.com/kennycch/gotools/list"

const (
	// 直线节点G值
	STRAIGHT = 10
	// 斜线节点G值
	DIAGONAL = 14
)

// 路径导航
type Jps struct {
	mapRow         int               // 地图列数（x轴）
	maxNode        int               // 地图最大坐标
	closeList      map[int]*Node     // 已探索节点
	openList       *list.List[*Node] // 备选中心节点
	obstacles      map[int]struct{}  // 障碍物
	stratIndex     int               // 开始节点
	currIndex      int               // 当前节点
	targetIndex    int               // 目标节点
	isSearchTarget bool              // 是否已找到目标节点
}

// 节点
type Node struct {
	historyCost     float64   // 历史代价（G值）
	estimateCost    float64   // 未来预计代价（H值）
	totalCost       float64   // 总代价（F值）
	index           int       // 节点坐标
	parentIndex     int       // 上一个节点
	forcedNeighbour []int     // 强迫邻居节点
	searchWay       *searchWay // 探索方向（起始节点为nil）
}

// 探索方向
type searchWay struct {
	x int
	y int
}
