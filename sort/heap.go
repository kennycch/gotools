package sort

import "github.com/kennycch/gotools/general"

/*
堆排序
array：原数组
newArray：排序后数组
*/
func Heap[V any, T general.Number](array []V, sortType SortType, sortValue func(vlaue V) T) (newArray []V) {
	newArray = make([]V, len(array))
	copy(newArray, array)
	// 长度少于2直接返回
	if len(newArray) < 2 {
		return
	}
	h := newHeap(newArray, sortType, sortValue)
	// 元素入堆
	for _, value := range newArray {
		h.push(value)
	}
	// 堆顶移至底部，完成排序
	for i := 0; i < len(newArray); i++ {
		h.pop()
	}
	newArray = h.array
	return
}

type heapSort[V any, T general.Number] struct {
	sortType SortType // 排序类型，决定是用大根堆还是小根堆
	size     int      // 堆的大小
	/*
		使用数组来模拟树
		假设某个节点下标为i
		父节点为(i - 1) / 2
		两个子节点分别为i * 2 + 1, i * 2 + 2
	*/
	array     []V
	sortValue func(V) T
}

// 创建堆
func newHeap[V any, T general.Number](array []V, sortType SortType, sortValue func(vlaue V) T) *heapSort[V, T] {
	return &heapSort[V, T]{
		sortType:  sortType,
		array:     make([]V, len(array)),
		sortValue: sortValue,
	}
}

// 元素入堆
func (h *heapSort[V, T]) push(value V) {
	// 没有元素时，直接置于堆顶
	if h.size == 0 {
		h.array[0] = value
		h.size++
		return
	}
	// 定义要插入节点的下标
	i := h.size
	// 定位入堆元素的最终下标
	for i > 0 {
		// 该元素父亲节点的下标
		p := (i - 1) / 2
		// 根据排序类型决定插入元素位置
		if (h.sortValue(h.array[p]) >= h.sortValue(value) && h.sortType == ASC) ||
			(h.sortValue(h.array[p]) <= h.sortValue(value) && h.sortType == DESC) { // 如果入堆元素大于/小于父节点，直接退出循环，父节点保持不变
			break
		} else { // 插入元素与父节点换位
			h.array[i] = h.array[p]
			i = p
		}
	}
	h.array[i] = value
	h.size++
}

// 堆顶移至底部
func (h *heapSort[V, T]) pop() {
	// 堆大小为0时不作任何操作
	if h.size == 0 {
		return
	}
	// 获取根节点和堆底元素
	top := h.array[0]
	// 堆大小减一，防止下次移动根节点元素依然参与
	h.size--
	// 获取堆底元素
	bottom := h.array[h.size]
	// 根节点元素置于底部
	h.array[h.size] = top
	// 开始翻转剩余元素，直至最小/大元素置于堆顶
	i := 0
	for {
		// 定位两个子节点
		left, right := i*2+1, i*2+2
		child := left
		// 左节点超出堆大小，证明不存在子节点，直接退出循环
		if left >= h.size {
			break
		}
		// 如果右节点比左节点小/大，由右节点来替换父节点
		if right < h.size && ((h.sortValue(h.array[left]) < h.sortValue(h.array[right]) && h.sortType == ASC) ||
			(h.sortValue(h.array[left]) > h.sortValue(h.array[right]) && h.sortType == DESC)) {
			child = right
		}
		// 底部值小/大于两个子节点，直接退出循环
		if (h.sortValue(bottom) >= h.sortValue(h.array[child]) && h.sortType == ASC) ||
			(h.sortValue(bottom) <= h.sortValue(h.array[child]) && h.sortType == DESC) {
			break
		}
		// 将最符合条件的子节点替换掉父节点
		h.array[i] = h.array[child]
		// 继续下一次比较
		i = child
	}
	// 将底部元素置于不会再被替换的位置
	h.array[i] = bottom
}
