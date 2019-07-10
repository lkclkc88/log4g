package log4g

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
	t := time.NewTimer(timeOut)
	select {
	case r := <-q.data:
		atomic.AddInt32(&q.size, -1)
		return r
	case <-t.C:
		{
			return nil
		}
	}
}

func (q *Queue) Size() int32 {
	return atomic.LoadInt32(&q.size)
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
	t := time.NewTimer(timeOut)
	select {
	case q.data <- obj:
		atomic.AddInt32(&q.size, 1)
		return true
	case <-t.C:
		{
			return false
		}
	}
}
