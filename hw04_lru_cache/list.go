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
	front *ListItem
	back  *ListItem
	len   int
}

func NewList() List {
	return new(list)
}

func (c *list) Len() int {
	return c.len
}

func (c *list) Front() *ListItem {
	return c.front
}

func (c *list) Back() *ListItem {
	return c.back
}

func (c *list) PushFront(v interface{}) *ListItem {
	defer func() {
		c.len++
	}()

	item := ListItem{
		Value: v,
	}

	if c.len == 0 {
		c.front = &item
		c.back = &item
		return &item
	}

	c.front.Prev = &item
	item.Next = c.front
	c.front = &item
	return &item
}

func (c *list) PushBack(v interface{}) *ListItem {
	defer func() {
		c.len++
	}()

	item := ListItem{
		Value: v,
	}

	if c.len == 0 {
		c.front = &item
		c.back = &item
		return &item
	}

	c.back.Next = &item
	item.Prev = c.back
	c.back = &item
	return &item
}

func (c *list) Remove(i *ListItem) {
	defer func() {
		c.len--
	}()

	if i == c.front {
		c.front = c.front.Next
		c.front.Prev = nil
		return
	}

	if i == c.back {
		c.back = c.back.Prev
		c.back.Next = nil
		return
	}

	i.Prev.Next = i.Next
	i.Next.Prev = i.Prev
}

func (c *list) MoveToFront(i *ListItem) {
	if c.front == i {
		return
	}

	c.Remove(i)
	i.Prev = nil
	i.Next = c.front
	c.front.Prev = i
	c.front = i
	c.len++
}
