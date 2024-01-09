package timer

import (
	"context"
	"gotools/worker"
	"math"
	"sync/atomic"
	"time"
)

// 新增循环任务
func AddTicker(times time.Duration, action func()) (taskId uint64) {
	// 自增任务ID
	if tickerManager.taskId == math.MaxUint64 {
		tickerManager.taskId = 0
	}
	taskId = atomic.AddUint64(&tickerManager.taskId, 1)
	// 设置中途拦截
	ctx, cancel := context.WithCancel(context.Background())
	// 注册任务
	tickerManager.taskMap.Set(taskId, cancel)
	// 任务执行
	tickerAction(times, action, ctx, taskId)
	return
}

// 循环任务执行
func tickerAction(times time.Duration, action func(), ctx context.Context, taskId uint64) {
	worker.AddTask(func() {
		// 生成延时器
		ticker := time.NewTicker(times)
		defer func() {
			// 关闭定时器
			ticker.Stop()
			// 注销任务
			tickerManager.taskMap.Del(taskId)
		}()
	label:
		for {
			select {
			// 时间到，执行任务
			case <-ticker.C:
				action()
			// 中途截断任务
			case <-ctx.Done():
				break label
			}
		}
	})
}

// 拦截延时任务
func DeleteTicker(taskId uint64) {
	if cancel, ok := tickerManager.taskMap.Get(taskId); ok {
		cancel()
	}
}
