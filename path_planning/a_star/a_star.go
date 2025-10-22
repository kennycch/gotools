package a_star

import (
	"math"

	"github.com/kennycch/gotools/general"
	"github.com/kennycch/gotools/list"
)

// NewAStar 生成路径导航对象
// 地图根据列与行数自动生成，左下角坐标为0，往右延伸
// 例：3 * 3
// 生成地图：
// ---------
// |6, 7, 8|
// |3, 4, 5|
// |0, 1, 2|
// ---------
//
//	@param row 地图列数（x轴）
//	@param cap 地图行数（y轴）
//	@param obstacles 障碍物坐标
//	@param diagonal 是否可对角移动
//	@return aStar A星导航对象
func NewAStar[V general.Signed](row, cap V, obstacles []V, diagonal bool) (aStar *AStar[V]) {
	if row <= 0 || cap <= 0 {
		return nil
	}
	// 初始化导航对象
	aStar = &AStar[V]{
		mapRow:    row,
		maxNode:   row*cap - 1,
		diagonal:  diagonal,
		obstacles: make(map[V]struct{}),
	}
	// 过滤地图外的障碍物
	for _, obstacle := range obstacles {
		if obstacle < 0 || obstacle > aStar.maxNode {
			continue
		}
		aStar.obstacles[obstacle] = struct{}{}
	}
	return aStar
}

// Planning 巡航
//
//	@receiver a
//	@param startIndex 开始节点
//	@param targetIndex 目标节点
//	@return path 规划的路径
//	@return result 是否成功生成路径规划
func (a *AStar[V]) Planning(startIndex, targetIndex V) (path []V, result bool) {
	if startIndex == targetIndex ||
		!a.isInMap(startIndex) || !a.isInMap(targetIndex) ||
		a.isObstacle(startIndex) || a.isObstacle(targetIndex) {
		return
	}
	// 赋值变量
	a.stratIndex = startIndex
	a.currIndex = startIndex
	a.targetIndex = targetIndex
	a.isSearchTarget = false
	// 清空列表
	a.closeList = make(map[V]*Node[V])
	a.openList = list.NewList(func(element *Node[V]) float64 {
		return element.totalCost
	})
	// 开始探索
	path, result = a.search()
	return
}

// search 开始探索节点
//
//	@receiver a
//	@return path 巡航路径
//	@return result 是否巡航成功
func (a *AStar[V]) search() (path []V, result bool) {
	path = make([]V, 0)
	// 将开始节点放入关闭列表
	a.closeList[a.stratIndex] = a.getNode(a.stratIndex, a.currIndex)
	for {
		// 将开始节点相邻节点放入开放列表和关闭列表
		a.searchAdjoinNode()
		// 如果已找到目标节点，终止循环
		if a.isSearchTarget {
			// 回溯节点
			path = a.backtrack()
			result = true
			break
		}
		// 如果已没有可作为探索中心的节点，终止探索
		nextNode, err := a.openList.PopFront()
		if err != nil {
			break
		}
		a.currIndex = nextNode.index
	}
	return
}

// getNode 计算节点代价
//
//	@receiver a
//	@param index 节点坐标
//	@param parentIndex 父节点坐标
//	@return *Node 生成的节点对象
func (a *AStar[V]) getNode(index, parentIndex V) *Node[V] {
	node := &Node[V]{
		index: index,
	}
	if index == parentIndex { // 自身就是父节点
		node.historyCost = 0
		node.parentIndex = -1
	} else {
		parentNode := a.closeList[parentIndex]
		node.parentIndex = parentIndex
		if a.diffCap(index, parentIndex) == 0 || a.diffRow(index, parentIndex) == 0 { // 水平/垂直节点
			node.historyCost = general.Decimal(parentNode.historyCost, STRAIGHT, general.Add)
		} else { // 对角节点
			node.historyCost = general.Decimal(parentNode.historyCost, DIAGONAL, general.Add)
		}
	}
	// 预估代价
	node.estimateCost = a.getEstimate(index)
	node.totalCost = general.Decimal(node.historyCost, node.estimateCost, general.Add)
	return node
}

// searchAdjoinNode 探索相邻节点
//
//	@receiver a
func (a *AStar[V]) searchAdjoinNode() {
	// 左节点
	a.addHorizontalNode(a.currIndex - 1)
	// 右节点
	if !a.isSearchTarget {
		a.addHorizontalNode(a.currIndex + 1)
	}
	// 上节点
	if !a.isSearchTarget {
		a.addVerticalNode(a.currIndex + a.mapRow)
	}
	// 下节点
	if !a.isSearchTarget {
		a.addVerticalNode(a.currIndex - a.mapRow)
	}
	// 对角节点
	if a.diagonal {
		// 左上节点
		if !a.isSearchTarget {
			a.addUpDiagonalNode(a.currIndex + a.mapRow - 1)
		}
		// 右上节点
		if !a.isSearchTarget {
			a.addUpDiagonalNode(a.currIndex + a.mapRow + 1)
		}
		// 左下节点
		if !a.isSearchTarget {
			a.addDownDiagonalNode(a.currIndex - a.mapRow - 1)
		}
		// 右下节点
		if !a.isSearchTarget {
			a.addDownDiagonalNode(a.currIndex - a.mapRow + 1)
		}
	}
}

// addHorizontalNode 探索水平节点
//
//	@receiver a
//	@param index 节点坐标
func (a *AStar[V]) addHorizontalNode(index V) {
	if a.isInMap(index) && a.diffCap(index, a.currIndex) == 0 && // 是否在地图范围、是否在同一行
		!a.isObstacle(index) && !a.isInCloseList(index) { // 是否是障碍物，是否已被探索过
		node := a.getNode(index, a.currIndex)
		a.closeList[index] = node
		a.openList.Insert(node)
		if index == a.targetIndex {
			a.isSearchTarget = true
		}
	}
}

// addVerticalNode 探索垂直节点
//
//	@receiver a
//	@param index 节点坐标
func (a *AStar[V]) addVerticalNode(index V) {
	if a.isInMap(index) && // 是否在地图范围
		!a.isObstacle(index) && !a.isInCloseList(index) { // 是否是障碍物，是否已被探索过
		node := a.getNode(index, a.currIndex)
		a.closeList[index] = node
		a.openList.Insert(node)
		if index == a.targetIndex {
			a.isSearchTarget = true
		}
	}
}

// addUpDiagonalNode 探索斜上节点
//
//	@receiver a
//	@param index 节点坐标
func (a *AStar[V]) addUpDiagonalNode(index V) {
	if a.isInMap(index) && index/a.mapRow-a.currIndex/a.mapRow == 1 && // 是否在地图范围、是否在当前节点上一行
		!a.isObstacle(index) && !a.isInCloseList(index) { // 是否是障碍物，是否已被探索过
		node := a.getNode(index, a.currIndex)
		a.closeList[index] = node
		a.openList.Insert(node)
		if index == a.targetIndex {
			a.isSearchTarget = true
		}
	}
}

// addDownDiagonalNode 探索斜下节点
//
//	@receiver a
//	@param index 节点坐标
func (a *AStar[V]) addDownDiagonalNode(index V) {
	if a.isInMap(index) && index/a.mapRow-a.currIndex/a.mapRow == -1 && // 是否在地图范围、是否在当前节点下一行
		!a.isObstacle(index) && !a.isInCloseList(index) { // 是否是障碍物，是否已被探索过
		node := a.getNode(index, a.currIndex)
		a.closeList[index] = node
		a.openList.Insert(node)
		if index == a.targetIndex {
			a.isSearchTarget = true
		}
	}
}

// getEstimate 计算预估代价
//
//	@receiver a
//	@param index 节点坐标
//	@return float64 预估代价
func (a *AStar[V]) getEstimate(index V) float64 {
	estimate := float64(0)
	if a.diagonal {
		estimate = a.euclidDIstance(index)
	} else {
		estimate = a.manhattanDistance(index)
	}
	return estimate
}

// manhattanDistance 曼哈顿距离
//
//	@receiver a
//	@param index 节点坐标
//	@return float64 预估代价
func (a *AStar[V]) manhattanDistance(index V) float64 {
	xDistance := math.Abs(float64(a.diffRow(index, a.targetIndex))) // x轴距离
	yDistance := math.Abs(float64(a.diffCap(index, a.targetIndex))) // y轴距离
	// 计算x和y的距离总和
	estimate := general.Decimal(xDistance, yDistance, general.Add)
	// 最终结果要乘以十
	estimate = general.Decimal(estimate, 10, general.Multiply)
	return estimate
}

// euclidDIstance 欧几里得距离
//
//	@receiver a
//	@param index 节点坐标
//	@return float64 预估代价
func (a *AStar[V]) euclidDIstance(index V) float64 {
	xDistance := math.Abs(float64(a.diffRow(index, a.targetIndex))) // x轴距离
	yDistance := math.Abs(float64(a.diffCap(index, a.targetIndex))) // y轴距离
	// 计算直线距离
	estimate := math.Sqrt(general.Decimal(math.Pow(xDistance, 2), math.Pow(yDistance, 2), general.Add))
	// 最终结果要乘以十
	estimate = general.Decimal(estimate, 10, general.Multiply)
	return estimate
}

// backtrack 回溯节点
//
//	@receiver a
//	@return []V 巡航路径
func (a *AStar[V]) backtrack() []V {
	path := make([]V, 0)
	currNode := a.closeList[a.targetIndex]
	path = append(path, currNode.index)
	for {
		// 如果已经在起点，停止回溯
		parentNode := a.closeList[currNode.parentIndex]
		path = append([]V{parentNode.index}, path...)
		if parentNode.parentIndex == -1 {
			break
		}
		currNode = parentNode
	}
	return path
}

// diffCap 行差
//
//	@receiver a
//	@param index 当前节点
//	@param contrast 对比节点
//	@return V 行差
func (a *AStar[V]) diffCap(index, contrast V) V {
	return index/a.mapRow - contrast/a.mapRow
}

// diffRow 列差
//
//	@receiver a
//	@param index 当前节点
//	@param contrast 对比节点
//	@return V 列差
func (a *AStar[V]) diffRow(index, contrast V) V {
	return index%a.mapRow - contrast%a.mapRow
}

// isInMap 是否在地图里
//
//	@receiver a
//	@param index 节点坐标
//	@return bool 是否在地图里
func (a *AStar[V]) isInMap(index V) bool {
	return index >= 0 && index <= a.maxNode
}

// isObstacle 是否障碍物
//
//	@receiver a
//	@param index 节点坐标
//	@return bool 是否障碍物
func (a *AStar[V]) isObstacle(index V) bool {
	_, ok := a.obstacles[index]
	return ok
}

// isInCloseList 是否在关闭列表
//
//	@receiver a
//	@param index 节点坐标
//	@return bool 是否在关闭列表
func (a *AStar[V]) isInCloseList(index V) bool {
	_, ok := a.closeList[index]
	return ok
}
