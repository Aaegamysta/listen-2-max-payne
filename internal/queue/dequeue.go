package queue

import "github.com/aaegamysta/listen-2-max-payne/internal/db"

type Dequeue interface {
	Enqueue(excerpt db.Excerpt) bool
	Dequeue() (db.Excerpt, bool)
	Peek() db.Excerpt
	Empty() bool
	Full() bool
}

type node struct {
	next *node
	data db.Excerpt
}

type linkedListBackedDequeue struct {
	front  *node
	rear   *node
	length int
	size   int
}

func New(size int) Dequeue {
	return &linkedListBackedDequeue{
		front:  nil,
		rear:   nil,
		size:   size,
		length: 0,
	}
}

func (l *linkedListBackedDequeue) Enqueue(excerpt db.Excerpt) bool {
	if l.length == l.size {
		return false
	}
	if l.front == nil {
		firstEntry := &node{
			next: nil,
			data: excerpt,
		}
		l.front, l.rear = firstEntry, firstEntry
	}
	ptr := l.rear
	ptr.next = &node{
		next: nil,
		data: excerpt,
	}
	l.rear = ptr
	l.length++
	return true
}

// Empty implements Dequeue.
func (l *linkedListBackedDequeue) Empty() bool {
	return l.front == nil
}

// Full implements Dequeue.
func (l *linkedListBackedDequeue) Full() bool {
	return l.length == l.size
}

// Dequeue implements Dequeue.
func (l *linkedListBackedDequeue) Dequeue() (db.Excerpt, bool) {
	// if the double ended queue is empty
	if doubleEndedQueueEmpty := l.front == nil && l.rear == nil; doubleEndedQueueEmpty {
		return db.Excerpt{}, false
	}
	// if there is only one element left, we have guarded against accessing the nil pointer that the front's next point to by seeing
	// if it points to the same
	if hasOnlyOneElement := l.front == l.rear && l.front != nil && l.rear != nil; hasOnlyOneElement {
		excerpt := l.front.data
		l.front, l.rear = nil, nil
		l.length = 0
		return excerpt, true
	}
	excerpt := l.front.data
	l.front = l.front.next
	l.length--
	return excerpt, true
}

func (l *linkedListBackedDequeue) Peek() db.Excerpt {
	return l.front.data
}
