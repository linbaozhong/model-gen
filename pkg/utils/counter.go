package utils

import (
	"sync"
	"sync/atomic"
)

type counter struct {
	counter uint64
	lock    sync.Mutex
}

var Counter counter

// 获取计数器的返回值，每次调用，返回值加1，到10000后重新从1计数
func GetCounter() uint64 {

	Counter.lock.Lock()

	// 不加锁，可能出现获取多次1的后果
	if Counter.counter >= 10000 {

		Counter.counter = 0
	}
	Counter.lock.Unlock()

	return atomic.AddUint64(&Counter.counter, 1)

}
