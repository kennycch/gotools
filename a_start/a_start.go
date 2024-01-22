package a_start

import (
	"math"

	"github.com/kennycch/gotools/general"
)

/*
生成路径导航对象
row：地图列数（x轴）
cap：地图行数（y轴）
obstacles：障碍物坐标
pathFinding：导航对象
*/
func NewPathFinding(row, cap int, obstacles []int) (pathFinding *PathFinding) {
	if row <= 0 || cap <= 0 {
		return nil
	}
	// 导航对象
	pathFinding = &PathFinding{
		mapRow:    row,
		mapCap:    cap,
		maps:      make([]int, 0),
		closeList: make([]int, 0),
		obstacles: make([]int, 0),
		nearNode:  make(map[int]*cost),
	}
	// 初始化地图节点
	for i := 0; i < row*cap; i++ {
		pathFinding.maps = append(pathFinding.maps, i)
	}
	// 过滤掉不在地图里的障碍物
	for _, obstacle := range obstacles {
		if general.InArray(pathFinding.maps, obstacle) {
			pathFinding.obstacles = append(pathFinding.obstacles, obstacle)
		}
	}
	return pathFinding
}

/*
巡航
startNode：开始节点
targetNode：目标节点
route：移动路径
result：是否成功生成导航路径
*/
func (p *PathFinding) Finding(startNode, targetNode int) (route []int, result bool) {
	route = make([]int, 0)
	p.closeList = make([]int, 0)
	// 判断起点与终点是否都在地图中
	if !general.InArray(p.maps, startNode) ||
		!general.InArray(p.maps, targetNode) {
		return
	}
	// 节点赋值
	p.currNode = startNode
	p.targetNode = targetNode
	p.closeList = append(p.closeList, startNode)
	// 开始导航
	if result = p.findingAction(); result {
		route = p.closeList
	}
	return route, result
}

// 执行导航
func (p *PathFinding) findingAction() bool {
	flag := false
	for {
		// 获取当前节点临近节点
		p.getNearNode()
		// 如果没有临近节点停止导航
		if len(p.nearNode) == 0 {
			break
		}
		// 选取下一个节点
		minCost := uint64(math.MaxUint64)
		nextNode := len(p.maps)
		for node, nearCost := range p.nearNode {
			if nearCost.totalCost <= minCost {
				nextNode = node
				minCost = nearCost.totalCost
			}
		}
		p.closeList = append(p.closeList, nextNode)
		p.currNode = nextNode
		if p.currNode == p.targetNode {
			flag = true
			break
		}
	}
	return flag
}

// 获取并计算临近节点
func (p *PathFinding) getNearNode() {
	p.nearNode = make(map[int]*cost)
	// 左节点
	p.horizontalHandle(p.currNode - 1)
	// 右节点
	p.horizontalHandle(p.currNode + 1)
	// 上节点
	p.verticalHandle(p.currNode - p.mapRow)
	// 下节点
	p.verticalHandle(p.currNode + p.mapRow)
	// 左上节点
	p.upOblique(p.currNode - p.mapRow - 1)
	// 右上节点
	p.upOblique(p.currNode - p.mapRow + 1)
	// 左下节点
	p.downOblique(p.currNode + p.mapRow - 1)
	// 右下节点
	p.downOblique(p.currNode + p.mapRow + 1)
}

// 水平节点处理
func (p *PathFinding) horizontalHandle(node int) {
	if node >= 0 && node < len(p.maps) && // 是否超出地图范围
		node/p.mapRow == p.currNode/p.mapRow && // 是否与当前节点同一行
		!general.InArray(p.obstacles, node) && // 是否不在障碍物中
		!general.InArray(p.closeList, node) { // 是否不在历史路径中
		p.addNearNode(node, STRAIGHT)
	}
}

// 垂直节点处理
func (p *PathFinding) verticalHandle(node int) {
	if node >= 0 && node < len(p.maps) && // 是否超出地图范围
		!general.InArray(p.obstacles, node) && // 是否不在障碍物中
		!general.InArray(p.closeList, node) { // 是否不在历史路径中
		p.addNearNode(node, STRAIGHT)
	}
}

// 上斜节点处理
func (p *PathFinding) upOblique(node int) {
	if node >= 0 && node < len(p.maps) && // 是否超出地图范围
		(node/p.mapRow)-(p.currNode/p.mapRow) == -1 && // 是否在当前节点上一行
		!general.InArray(p.obstacles, node) && // 是否不在障碍物中
		!general.InArray(p.closeList, node) { // 是否不在历史路径中
		p.addNearNode(node, OBLIQUE)
	}
}

// 下斜节点处理
func (p *PathFinding) downOblique(node int) {
	if node >= 0 && node < len(p.maps) && // 是否超出地图范围
		(node/p.mapRow)-(p.currNode/p.mapRow) == 1 && // 是否在当前节点上一行
		!general.InArray(p.obstacles, node) && // 是否不在障碍物中
		!general.InArray(p.closeList, node) { // 是否不在历史路径中
		p.addNearNode(node, OBLIQUE)
	}
}

// 添加临近节点
func (p *PathFinding) addNearNode(node int, historyCost costType) {
	nearCost := &cost{
		historyCost:  uint64(historyCost),
		estimateCost: p.getEstimateCost(node),
	}
	nearCost.totalCost = nearCost.historyCost + nearCost.estimateCost
	p.nearNode[node] = nearCost
}

// 获取节点未来预计代价（曼哈顿算法）
func (p *PathFinding) getEstimateCost(node int) uint64 {
	cost := uint64(0)
	if node != p.targetNode {
		// 计算水平距离
		horizontalCost := uint64(math.Abs(float64(node%p.mapRow - p.targetNode%p.mapRow)))
		// 计算垂直距离
		verticalCost := uint64(math.Abs(float64(node/p.mapRow - p.targetNode/p.mapCap)))
		// 最终未来预计代价要乘十
		cost = (horizontalCost + verticalCost) * 10
	}
	return cost
}
