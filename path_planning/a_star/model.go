package a_star

import (
	"github.com/kennycch/gotools/general"
	"github.com/kennycch/gotools/list"
)

const (
	// 直线节点G值
	STRAIGHT = 10
	// 斜线节点G值
	DIAGONAL = 14
)

// 路径导航
type AStar[V general.Signed] struct {
	mapRow         V                             // 地图列数（x轴）
	maxNode        V                             // 地图最大坐标
	diagonal       bool                          // 是否可以对角线移动
	closeList      map[V]*Node[V]                // 已探索节点
	openList       *list.List[*Node[V], float64] // 备选中心节点
	obstacles      map[V]struct{}                // 障碍物
	stratIndex     V                             // 开始节点
	currIndex      V                             // 当前节点
	targetIndex    V                             // 目标节点
	isSearchTarget bool                          // 是否已找到目标节点
}

// 节点
type Node[V general.Signed] struct {
	historyCost  float64 // 历史代价（G值）
	estimateCost float64 // 未来预计代价（H值）
	totalCost    float64 // 总代价（F值）
	index        V       // 节点坐标
	parentIndex  V       // 上一个节点
}
