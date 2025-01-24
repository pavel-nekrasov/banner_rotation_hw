package cache

type ListItem[T interface{}] struct {
	Value T
	Next  *ListItem[T]
	Prev  *ListItem[T]
}

func (next *ListItem[T]) wire(prev *ListItem[T]) {
	if next != nil {
		next.Prev = prev
	}
	if prev != nil {
		prev.Next = next
	}
}
