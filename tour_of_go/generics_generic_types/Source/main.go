package main

import (
	"fmt"
	"github.com/playdate-go/pdgo"
)

var pd *pdgo.PlaydateAPI
	

// List represents a singly-linked list that holds
// values of any type.
type List[T any] struct {
	next *List[T]
	val  T
}

// PushFront prepends a value to the front of the list
// and returns the new head.
func (l *List[T]) PushFront(v T) *List[T] {
	return &List[T]{val: v, next: l}
}

// PushBack appends a value to the end of the list.
func (l *List[T]) PushBack(v T) *List[T] {
	if l == nil {
		return &List[T]{val: v}
	}
	cur := l
	for cur.next != nil {
		cur = cur.next
	}
	cur.next = &List[T]{val: v}
	return l
}

// Len returns the number of elements in the list.
func (l *List[T]) Len() int {
	n := 0
	for cur := l; cur != nil; cur = cur.next {
		n++
	}
	return n
}

// Each calls fn for every value in the list.
func (l *List[T]) Each(fn func(T)) {
	for cur := l; cur != nil; cur = cur.next {
		fn(cur.val)
	}
}

// initGame is called once when the game starts
func initGame() {
	// Build a list: 10 -> 20 -> 30
	l := new(List[int])
	l = l.PushFront(20).PushFront(10)
	l = l.PushBack(30)
	pd.Graphics.DrawText("values: ", 50, 50)
	
	yPos := 70
	
	l.Each(func(v int) {
		pd.Graphics.DrawText(fmt.Sprint(v), 50, yPos)
		yPos+= 20
	})
	
	yPos+= 20
	
	pd.Graphics.DrawText(fmt.Sprint("length:", l.Len()), 50, yPos)
}

// update is called every frame
func update() int {
	return 1
}

func main() {}

