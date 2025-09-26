package lrucache

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
	first *ListItem
	last  *ListItem
	len   int
}

func NewList() List {
	return new(list)
}

// Len Длина списка.
func (l *list) Len() int {
	return l.len
}

// Front первый элемент списка.
func (l *list) Front() *ListItem {
	return l.first
}

// Back последний элемент списка.
func (l *list) Back() *ListItem {
	return l.last
}

// PushFront добавить значение в начало.
func (l *list) PushFront(v interface{}) *ListItem {
	newItem := &ListItem{Value: v}
	if l.len == 0 {
		l.first = newItem
		l.last = newItem
	} else {
		newItem.Next = l.first
		l.first.Prev = newItem
		l.first = newItem
	}
	l.len++
	return newItem
}

// PushBack добавить значение в конец.
func (l *list) PushBack(v interface{}) *ListItem {
	newItem := &ListItem{Value: v}
	if l.len == 0 {
		l.first = newItem
		l.last = newItem
	} else {
		newItem.Prev = l.last
		l.last.Next = newItem
		l.last = newItem
	}
	l.len++
	return newItem
}

// Remove удалить элемент.
func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.first = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.last = i.Prev
	}
	l.len--
}

// MoveToFront переместить элемент в начало.
func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}
