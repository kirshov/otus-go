package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	List
	length int
	first  *ListItem
	last   *ListItem
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	if l.first != nil {
		return l.first
	}

	return nil
}

func (l *list) Back() *ListItem {
	if l.last == nil && l.length == 1 {
		return l.first
	}

	return l.last
}

func (l *list) PushFront(v interface{}) *ListItem {
	n := &ListItem{Value: v}
	l.length++
	if l.first != nil {
		l.first.Prev = n
		n.Next = l.first
	}

	l.first = n
	return n
}

func (l *list) PushBack(v interface{}) *ListItem {
	l.length++
	n := &ListItem{Value: v}
	if l.last == nil {
		if l.Front() != nil {
			n.Prev = l.first
			l.first.Next = n
		}
	} else {
		l.last.Next = n
		n.Prev = l.last
	}

	l.last = n
	return n
}

func (l *list) Remove(i *ListItem) {
	l.length--

	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	i = nil
}

func (l *list) MoveToFront(i *ListItem) {
	l.PushFront(i.Value)
	l.Remove(i)
}

func NewList() List {
	return new(list)
}
