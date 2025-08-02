package timer

import (
	"container/list"
	"sync"
	"time"

	"github.com/alphadose/haxmap"
)

// TimeWheel 时间轮结构体
type TimeWheel struct {
	interval    time.Duration              // 每个槽位的时间间隔
	slotsNum    int                        // 槽位数量
	slots       []*list.List               // 每个槽位的任务链表
	slotLocks   []sync.Mutex               // 每个槽位的互斥锁
	currentPos  int                        // 当前指针位置
	currentTime int64                      // 当前时间（毫秒）
	ticker      *time.Ticker               // 定时器
	stopChan    chan struct{}              // 停止信号
	taskMap     *haxmap.Map[uint64, *Task] // 任务ID到任务的映射
	idGen       uint64                     // 任务ID生成器
}

// Task 延时任务结构体
type Task struct {
	id         uint64        // 任务唯一ID
	expiration int64         // 过期时间戳（毫秒）
	callback   func()        // 任务回调函数
	slotPos    int           // 所在槽位
	element    *list.Element // 在链表中的元素
}
