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
	if q.rear == len(q.arr)-1 {
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
	if q.rear == len(q.arr)-1 {
		return true
	}
	return false
}

type node struct {
	next  *node
	value *data.Excerpt
}

type linkedListBackedQueue struct {
	head  *node
	front *node
	rear  *node
	length int
	size  int
}

func NewLinkedListBackedQueue(size int) Queue {
	var head, front, rear *node
	return &linkedListBackedQueue{
		head:  head,
		front: front,
		rear:  rear,
		length: 0,
		size:  size,
	}
}

func (q *linkedListBackedQueue) Empty() bool {
	if q.rear == nil {
		return true
	}
	return false
}

func (q *linkedListBackedQueue) Full() bool {
	if q.length == q.size {
		return true
	}
	return false
}

func (q *linkedListBackedQueue) Dequeue() (data.Excerpt, bool) {
	if q.front == nil {
		return data.Excerpt{}, false
	}
	e := q.front
	q.front = q.front.next
	q.length--
	return *e.value, false
}

func (q *linkedListBackedQueue) Enqueue(e data.Excerpt) (int, bool) {
	if q.Full() {
		return q.size, false
	}
	if q.rear == nil {
		q.rear = &node{
			next: nil,
			value: &e,
		}
		q.front = q.rear 
		q.length++
		return q.size, true
	}
	n := q.rear
	n.next = &node{
		next:  nil,
		value: &e,
	}	
	q.rear = n.next
	q.length++
	return q.size, true
}
