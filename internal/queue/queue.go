package queue

import "github.com/aaegamysta/listen-2-max-payne/internal/data"

type Queue interface {
	Enqueue(data.Excerpt) (int, bool)
	Dequeue() (data.Excerpt, bool)
	Empty() bool
	Full() bool
}

type arrayBackedQueue struct {
	arr   []*data.Excerpt
	front int
	rear  int
	size  int
}

func NewArrayBackedQueue(size int) Queue {
	arr := make([]*data.Excerpt, size)
	for i := 0; i < size; i++ {
		arr[i] = nil
	}
	return &arrayBackedQueue{
		arr:   arr,
		front: -1,
		rear:  -1,
		size:  size,
	}
}

func (q *arrayBackedQueue) Dequeue() (data.Excerpt, bool) {
	if q.rear == q.front {
		return data.Excerpt{}, false
	}
	q.front++
	e := q.arr[q.front]
	q.arr[q.front] = nil
	return *e, true
}

func (q *arrayBackedQueue) Enqueue(e data.Excerpt) (int, bool) {
	if q.rear == len(q.arr) - 1 {
		return q.rear, false
	}
	q.rear++
	q.arr[q.rear] = &e	
	return q.rear, true
}

func (q *arrayBackedQueue) Empty() bool {
	if q.front == q.rear {
		return true
	}
	return false
}

func (q *arrayBackedQueue) Full() bool {
	if q.rear == len(q.arr) - 1 {
		return true
	}
	return false
}
