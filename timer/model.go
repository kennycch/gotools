package timer

import (
	"context"

	"github.com/alphadose/haxmap"
)

// 任务结构体
type task struct {
	taskId  uint64
	taskMap haxmap.Map[uint64, context.CancelFunc]
}

var (
	// 延时任务管理器
	timerManager = &task{
		taskId:  0,
		taskMap: *haxmap.New[uint64, context.CancelFunc](),
	}
	// 循环任务管理器
	tickerManager = &task{
		taskId:  0,
		taskMap: *haxmap.New[uint64, context.CancelFunc](),
	}
)
