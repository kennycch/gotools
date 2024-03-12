package a_star

import (
	"math"

	"github.com/kennycch/gotools/general"
	"github.com/kennycch/gotools/list"
)

/*
生成路径导航对象
地图根据列与行数自动生成，左下角坐标为0，往右延伸
例：3 * 3
生成地图：
---------
|6, 7, 8|
|3, 4, 5|
|0, 1, 2|
---------
row：地图列数（x轴）
cap：地图行数（y轴）
obstacles：障碍物坐标
diagonal：是否可对角移动
aStar：A星导航对象
*/
func NewAStar(row, cap int, obstacles []int, diagonal bool) (aStar *AStar) {
	if row <= 0 || cap <= 0 {
		return nil
	}
	// 初始化导航对象
	aStar = &AStar{
		mapRow:    row,
		maxNode:   row*cap - 1,
		diagonal:  diagonal,
		obstacles: make(map[int]struct{}),
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

/*
巡航
startIndex：开始节点
targetIndex：目标节点
route：规划的路径
result：是否成功生成路径规划
*/
func (a *AStar) Planning(startIndex, targetIndex int) (path []int, result bool) {
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
	a.closeList = make(map[int]*Node)
	a.openList = list.NewList[*Node]()
	// 开始探索
	path, result = a.search()
	return
}

// 开始探索节点
func (a *AStar) search() (path []int, result bool) {
	path = make([]int, 0)
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

// 计算节点代价
func (a *AStar) getNode(index, parentIndex int) *Node {
	node := &Node{
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

// 探索相邻节点
func (a *AStar) searchAdjoinNode() {
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

// 探索水平节点
func (a *AStar) addHorizontalNode(index int) {
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

// 探索垂直节点
func (a *AStar) addVerticalNode(index int) {
	if a.isInMap(index) && a.diffRow(index, a.currIndex) == 0 && // 是否在地图范围、是否在同一行
		!a.isObstacle(index) && !a.isInCloseList(index) { // 是否是障碍物，是否已被探索过
		node := a.getNode(index, a.currIndex)
		a.closeList[index] = node
		a.openList.Insert(node)
		if index == a.targetIndex {
			a.isSearchTarget = true
		}
	}
}

// 探索斜上节点
func (a *AStar) addUpDiagonalNode(index int) {
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

// 探索斜下节点
func (a *AStar) addDownDiagonalNode(index int) {
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

// 计算预估代价
func (a *AStar) getEstimate(index int) float64 {
	estimate := float64(0)
	if a.diagonal {
		estimate = a.euclidDIstance(index)
	} else {
		estimate = a.manhattanDistance(index)
	}
	return estimate
}

// 曼哈顿距离
func (a *AStar) manhattanDistance(index int) float64 {
	xDistance := math.Abs(float64(a.diffRow(index, a.targetIndex))) // x轴距离
	yDistance := math.Abs(float64(a.diffCap(index, a.targetIndex))) // y轴距离
	// 计算x和y的距离总和
	estimate := general.Decimal(xDistance, yDistance, general.Add)
	// 最终结果要乘以十
	estimate = general.Decimal(estimate, 10, general.Multiply)
	return estimate
}

// 欧几里得距离
func (a *AStar) euclidDIstance(index int) float64 {
	xDistance := math.Abs(float64(a.diffRow(index, a.targetIndex))) // x轴距离
	yDistance := math.Abs(float64(a.diffCap(index, a.targetIndex))) // y轴距离
	// 计算直线距离
	estimate := math.Sqrt(general.Decimal(math.Pow(xDistance, 2), math.Pow(yDistance, 2), general.Add))
	// 最终结果要乘以十
	estimate = general.Decimal(estimate, 10, general.Multiply)
	return estimate
}

// 回溯节点
func (a *AStar) backtrack() []int {
	path := make([]int, 0)
	lastNode := a.closeList[a.targetIndex]
	path = append(path, lastNode.index)
	for {
		// 如果已经在起点，停止回溯
		parentNode := a.closeList[lastNode.parentIndex]
		path = append([]int{parentNode.index}, path...)
		if parentNode.parentIndex == -1 {
			break
		}
		lastNode = parentNode
	}
	return path
}

// 行差
func (a *AStar) diffCap(index, contrast int) int {
	return index/a.mapRow - contrast/a.mapRow
}

// 列差
func (a *AStar) diffRow(index, contrast int) int {
	return index%a.mapRow - contrast%a.mapRow
}

// 是否在地图里
func (a *AStar) isInMap(index int) bool {
	return index >= 0 && index <= a.maxNode
}

// 是否障碍物
func (a *AStar) isObstacle(index int) bool {
	_, ok := a.obstacles[index]
	return ok
}

// 是否在关闭列表
func (a *AStar) isInCloseList(index int) bool {
	_, ok := a.closeList[index]
	return ok
}

// 返回节点总代价
func (n *Node) GetScore() float64 {
	return n.totalCost
}
