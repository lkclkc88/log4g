package imLog

import (
	"sync/atomic"
	"time"
)

//队列
type Queue struct {
	size int32
	data chan interface{}
}

func NewQueue(size uint64) *Queue {
	queue := Queue{data: make(chan interface{}, size), size: 0}
	return &queue
}

func (q *Queue) Poll(timeOut time.Duration) interface{} {
	tick := time.NewTicker(timeOut)
	defer tick.Stop()
	select {
	case r := <-q.data:
		atomic.AddInt32(&q.size, -1)
		return r
	case <-tick.C:
		{
			return nil
		}
	}
}

func (q *Queue) Get() interface{} {
	for {
		num := atomic.LoadInt32(&q.size)
		if num < 1 {
			return nil
		} else if atomic.CompareAndSwapInt32(&q.size, num, num-1) {
			return <-q.data
		}
	}
}

func (q *Queue) Offer(obj interface{}, timeOut time.Duration) bool {
	tick := time.NewTicker(timeOut)
	defer tick.Stop()
	select {
	case q.data <- obj:
		atomic.AddInt32(&q.size, 1)
		return true
	case <-tick.C:
		{
			return false
		}
	}
}
