package timer2

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/alphadose/haxmap"
	"github.com/kennycch/gotools/general"
	"github.com/kennycch/gotools/list"
	"github.com/kennycch/gotools/worker"
)

// AddTimer 新增延时任务
//
//	@param times 时间
//	@param action 要执行的任务
//	@return taskId 任务唯一Id
func AddTimer(times time.Duration, action func()) (taskId uint64) {
	checkInit()
	// 唯一Id
	taskId = atomic.AddUint64(&timerWheelManager.taskId, 1)
	// 任务唯一Id映射
	timerWheelManager.taskMap.Set(taskId, action)
	// 任务入列
	task := &TimerTask{
		taskId:     taskId,
		actionTime: general.NowAdd(times).UnixMilli(),
	}
	timerWheelManager.actionList.Insert(task)
	return
}

// DeleteTimer 拦截延时任务
//
//	@param taskId 任务唯一Id
func DeleteTimer(taskId uint64) {
	checkInit()
	// 删除任务唯一Id对应的任务
	timerWheelManager.taskMap.Del(taskId)
}

// initTimerWheel 初始化时间轮
//
//	@param interval 间隔时长（毫秒）
func initTimerWheel(interval int64) {
	once.Do(func() {
		// 结束任务
		ctx, cancel := context.WithCancel(context.Background())
		// 初始化对象
		timerWheelManager = &timerWheel{
			taskMap: haxmap.New[uint64, func()](),
			actionList: list.NewList(func(t *TimerTask) int64 {
				return t.actionTime
			}),
			interval: interval,
			cancel:   cancel,
			lock:     &sync.Mutex{},
		}
		// 标记已初始化
		isInit = true
		// 开始监听
		worker.AddTask(func() {
			loop(ctx)
		})
	})
}

// loop 监听任务
//
//	@param ctx 上下文
func loop(ctx context.Context) {
	ticket := time.NewTicker(time.Duration(timerWheelManager.interval) * time.Millisecond)
label:
	for {
		select {
		case <-ticket.C:
			worker.AddTask(func() {
				now := general.NowMilli()
				// 抢夺锁
				timerWheelManager.lock.Lock()
				// 遍历list执行任务
				lastIndex := -1
				timerWheelManager.actionList.Iterator(func(i int, t *TimerTask) bool {
					if t.actionTime <= now {
						// 执行延时任务
						if f, ok := timerWheelManager.taskMap.Get(t.taskId); ok {
							worker.AddTask(f)
							timerWheelManager.taskMap.Del(t.taskId)
						}
						lastIndex = i
						return true
					}
					return false
				})
				if lastIndex >= 0 {
					timerWheelManager.actionList.DeleteBefore(lastIndex + 1)
				}
				timerWheelManager.lock.Unlock()
			})
		case <-ctx.Done():
			break label
		}
	}
}

// checkInit 检查初始化
func checkInit() {
	if !isInit {
		panic("timer is not init")
	}
}
