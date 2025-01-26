package cache

import "sync"

type LinkedList[T interface{}] interface {
	Len() int
	Front() *ListItem[T]
	Back() *ListItem[T]
	PushFront(v T) *ListItem[T]
	PushBack(v T) *ListItem[T]
	Remove(i *ListItem[T])
	MoveToFront(i *ListItem[T])
}

type list[T interface{}] struct {
	mut   sync.RWMutex
	front *ListItem[T]
	back  *ListItem[T]
	len   int
}

func NewList[T interface{}]() LinkedList[T] {
	return new(list[T])
}

func (l *list[T]) Len() int {
	l.mut.RLock()
	defer l.mut.RUnlock()

	return l.len
}

func (l *list[T]) Front() *ListItem[T] {
	l.mut.RLock()
	defer l.mut.RUnlock()

	return l.front
}

func (l *list[T]) Back() *ListItem[T] {
	l.mut.RLock()
	defer l.mut.RUnlock()

	return l.back
}

func (l *list[T]) PushFront(v T) *ListItem[T] {
	l.mut.Lock()
	defer l.mut.Unlock()

	i := &ListItem[T]{Value: v, Next: l.front}

	if l.front != nil {
		l.front.wire(i)
	} else {
		l.back = i
	}
	l.front = i
	l.len++
	return i
}

func (l *list[T]) PushBack(v T) *ListItem[T] {
	l.mut.Lock()
	defer l.mut.Unlock()

	i := &ListItem[T]{Value: v, Next: nil, Prev: l.back}
	if l.back != nil {
		i.wire(l.back)
	} else {
		l.front = i
	}
	l.back = i
	l.len++
	return i
}

func (l *list[T]) Remove(i *ListItem[T]) {
	l.mut.Lock()
	defer l.mut.Unlock()

	if i == nil {
		return
	}

	prevItem := i.Prev
	nextItem := i.Next
	nextItem.wire(prevItem)

	if nextItem == nil {
		l.back = prevItem
	}
	if prevItem == nil {
		l.front = nextItem
	}
	i.Prev = nil
	i.Next = nil

	l.len--
}

func (l *list[T]) MoveToFront(i *ListItem[T]) {
	l.mut.Lock()
	defer l.mut.Unlock()

	if i == nil {
		return
	}

	prevItem := i.Prev
	if prevItem == nil {
		return
	}

	i.Next.wire(prevItem)
	if i.Next == nil {
		l.back = prevItem
	}

	l.front.wire(i)
	l.front = i
	i.Prev = nil
}
