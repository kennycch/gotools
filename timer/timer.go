package timer

import (
	"context"
	"gotools/worker"
	"math"
	"sync/atomic"
	"time"
)

// 新增延时任务
func AddTimer(times time.Duration, action func()) (taskId uint64) {
	// 自增任务ID
	if timerManager.taskId == math.MaxUint64 {
		timerManager.taskId = 0
	}
	taskId = atomic.AddUint64(&timerManager.taskId, 1)
	// 设置中途拦截
	ctx, cancel := context.WithCancel(context.Background())
	// 注册任务
	timerManager.taskMap.Set(taskId, cancel)
	// 任务执行
	timerAction(times, action, ctx, taskId)
	return
}

// 延时任务执行
func timerAction(times time.Duration, action func(), ctx context.Context, taskId uint64) {
	worker.AddTask(func() {
		// 生成延时器
		timer := time.NewTimer(times)
		defer func() {
			// 关闭定时器
			timer.Stop()
			// 注销任务
			timerManager.taskMap.Del(taskId)
		}()
		select {
		case <-timer.C:
			action()
		case <-ctx.Done():
		}
	})
}

// 拦截延时任务
func DeleteTimer(taskId uint64) {
	if cancel, ok := timerManager.taskMap.Get(taskId); ok {
		cancel()
	}
}
