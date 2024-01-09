package worker

import (
	"github.com/panjf2000/ants/v2"
	"github.com/panjf2000/gnet/pkg/pool/goroutine"
)

// 把任务压入任务缓冲池
func AddTask(data func()) {
	workerPool.Submit(data)
}

// 创建缓冲池对象
func newWorkPool() *ants.Pool {
	workerPool, _ := ants.NewPool(goroutine.DefaultAntsPoolSize,
		ants.WithOptions(ants.Options{ExpiryDuration: goroutine.ExpiryDuration,
			Nonblocking:  goroutine.Nonblocking,
			PanicHandler: nil,
		}))
	return workerPool
}
