package list

import "sync"

type listType interface {
	GetScore() float64
}

type List[T listType] struct {
	lists []T
	lock  *sync.RWMutex
}
