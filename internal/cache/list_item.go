package cache

type ListItem[T any] struct {
	Value T
	Next  *ListItem[T]
	Prev  *ListItem[T]
}

func (next *ListItem[T]) wireToPrev(prev *ListItem[T]) {
	if next != nil {
		next.Prev = prev
	}
	if prev != nil {
		prev.Next = next
	}
}
