package main

import (
	"fmt"
)

func Println[T any](val T) {
	fmt.Println(val)
}

type List[T any] struct {
	front *listNode[T]
	back  *listNode[T]
}

type listNode[T any] struct {
	elt  T
	next *listNode[T]
}

func (l *List[T]) Prepend(val T) {
	l.front = &listNode[T]{
		elt:  val,
		next: l.front
	}
}

func (l *List[T]) Append(val T) {
	l.back.next = &listNode[T]{
		elt:  val,
	}
	l.back = l.back.next
}

func (l *List[T]) Get(i int) (val T, found bool) {
	node := l.front
	for ; node != nil && i > 0; i-- {
		node = node.next
	}
	if i != 0 {
		return T{}, false
	}
	return node.elt, true
}
