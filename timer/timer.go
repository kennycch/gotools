package timer

import (
	"container/list"
	"sync"
	"sync/atomic"
	"time"

	"github.com/alphadose/haxmap"
	"github.com/kennycch/gotools/general"
	"github.com/kennycch/gotools/worker"
)

// New 创建时间轮实例
func New(interval time.Duration, slotsNum int) *TimeWheel {
	tw := &TimeWheel{
		interval:  interval,
		slotsNum:  slotsNum,
		slots:     make([]*list.List, slotsNum),
		slotLocks: make([]sync.Mutex, slotsNum),
		stopChan:  make(chan struct{}),
		ticker:    time.NewTicker(interval),
		taskMap:   haxmap.New[uint64, *Task](),
	}
	for i := 0; i < slotsNum; i++ {
		tw.slots[i] = list.New()
	}
	atomic.StoreInt64(&tw.currentTime, general.NowNano()/int64(time.Millisecond))
	worker.AddTask(tw.run)
	return tw
}

// Stop 停止时间轮
func (tw *TimeWheel) Stop() {
	close(tw.stopChan)
}

// Add 添加延时任务，返回任务ID
func (tw *TimeWheel) Add(delay time.Duration, callback func()) uint64 {
	id := atomic.AddUint64(&tw.idGen, 1)
	now := atomic.LoadInt64(&tw.currentTime)
	expiration := now + int64(delay/time.Millisecond)
	// 计算目标槽位
	delta := expiration - now
	if delta < 0 {
		delta = 0
	}
	steps := int(delta)
	pos := (tw.currentPos + steps) % tw.slotsNum
	// 创建任务并插入槽位
	task := &Task{
		id:         id,
		expiration: expiration,
		callback:   callback,
		slotPos:    pos,
	}

	tw.slotLocks[pos].Lock()
	task.element = tw.slots[pos].PushBack(task)
	tw.slotLocks[pos].Unlock()
	// 记录任务到映射表
	tw.taskMap.Set(id, task)
	return id
}

// Cancel 取消指定任务
func (tw *TimeWheel) Cancel(id uint64) bool {
	task, ok := tw.taskMap.Get(id)
	if !ok {
		return false
	}
	tw.slotLocks[task.slotPos].Lock()
	tw.slots[task.slotPos].Remove(task.element)
	tw.slotLocks[task.slotPos].Unlock()
	tw.taskMap.Del(id)
	return true
}

// run 运行时间轮主循环
func (tw *TimeWheel) run() {
	defer tw.ticker.Stop()
	for {
		select {
		case <-tw.ticker.C:
			now := time.Now().UnixNano() / int64(time.Millisecond)
			oldTime := atomic.LoadInt64(&tw.currentTime)
			elapsed := now - oldTime
			if elapsed > 0 {
				tw.advance(int(elapsed))
			}
		case <-tw.stopChan:
			return
		}
	}
}

// advance 推进时间轮
func (tw *TimeWheel) advance(steps int) {
	if steps <= 0 {
		return
	}
	// 更新当前时间
	atomic.AddInt64(&tw.currentTime, int64(steps))
	startPos := tw.currentPos
	slotsNum := tw.slotsNum
	// 处理所有跳过的槽位
	for i := 0; i < steps; i++ {
		pos := (startPos + i + 1) % slotsNum
		tw.processSlot(pos)
	}
	// 更新指针位置
	tw.currentPos = (startPos + steps) % slotsNum
}

// processSlot 处理指定槽位的任务
func (tw *TimeWheel) processSlot(pos int) {
	// 获取槽位锁并提取任务
	tw.slotLocks[pos].Lock()
	l := tw.slots[pos]
	tasks := make([]*Task, 0, l.Len())
	for e := l.Front(); e != nil; e = e.Next() {
		tasks = append(tasks, e.Value.(*Task))
	}
	l.Init() // 清空槽位
	tw.slotLocks[pos].Unlock()
	// 处理所有任务
	now := atomic.LoadInt64(&tw.currentTime)
	for _, task := range tasks {
		if task.expiration <= now {
			// 执行到期任务
			worker.AddTask(func() {
				task.callback()
			})
			tw.taskMap.Del(task.id)
		} else {
			// 重新计算并插入到新槽位
			delta := task.expiration - now
			if delta < 0 {
				delta = 0
			}
			newPos := (tw.currentPos + int(delta)) % tw.slotsNum
			task.slotPos = newPos

			tw.slotLocks[newPos].Lock()
			task.element = tw.slots[newPos].PushBack(task)
			tw.slotLocks[newPos].Unlock()
		}
	}
}
