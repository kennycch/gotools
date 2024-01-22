package a_start

const (
	// 直线节点G值
	STRAIGHT costType = 10
	// 斜线节点G值
	OBLIQUE costType = 14
)

type costType int

// 路径导航
type PathFinding struct {
	mapRow     int           // 地图列数（x轴）
	mapCap     int           // 地图行数（y轴）
	maps       []int         // 地图
	closeList  []int         // 已选择的路径
	obstacles  []int         // 障碍物
	nearNode   map[int]*cost // 临近节点
	currNode   int           // 当前节点
	targetNode int           // 目标节点
}

// 代价
type cost struct {
	historyCost  uint64 // 历史代价（G值）
	estimateCost uint64 // 未来预计代价（H值）
	totalCost    uint64 // 总代价（F值）
}
