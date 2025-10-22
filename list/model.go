package list

import (
	"sync"

	"github.com/kennycch/gotools/general"
)

type List[T any, V general.Ordered] struct {
	lists    []T
	lock     *sync.RWMutex
	sortFunc func(T) V
}
