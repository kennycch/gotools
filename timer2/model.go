package timer2

import (
	"context"
	"sync"

	"github.com/alphadose/haxmap"
	"github.com/kennycch/gotools/list"
)

type timerWheel struct {
	// 任务唯一Id
	taskId uint64
	// 任务唯一Id映射
	taskMap *haxmap.Map[uint64, func()]
	// 任务链表
	actionList *list.List[*TimerTask, int64]
	// 检测间隔
	interval int64
	// 取消任务方法
	cancel context.CancelFunc
	// 写入锁
	lock *sync.Mutex
}

type TimerTask struct {
	// 任务唯一Id
	taskId uint64
	// 执行时间
	actionTime int64
}

var (
	// 是否已初始化
	isInit bool = false
	// 单例模式
	once = &sync.Once{}
	// 时间轮管理器
	timerWheelManager *timerWheel
)
