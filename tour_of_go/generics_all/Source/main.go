package main

import (
	"fmt"

	"github.com/playdate-go/pdgo"
)

var pd *pdgo.PlaydateAPI

// ========== CORE 8 PATTERNS (EXISTING) ==========

// 1. Generic Functions
func Print[T any](v T) { pd.Graphics.DrawText(fmt.Sprint(v), 10, 10) }

func Min[T ~int|~float64](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// 2. Generic Stack
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var z T
		return z, false
	}
	v := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return v, true
}

// 3. Constraint interface
type Number interface{ ~int | ~float64 | ~string }

func SumNumbers[T Number](nums []T) T {
	var total T
	for _, n := range nums {
		total += n
	}
	return total
}

// 4. Generic Queue
type Queue[T comparable] struct {
	items []T
}

func (q *Queue[T]) Enqueue(v T) { q.items = append(q.items, v) }

func (q *Queue[T]) Contains(v T) bool {
	for _, item := range q.items {
		if item == v {
			return true
		}
	}
	return false
}

// 5. Approximate types
type Ordered interface{ ~int | ~float64 }

func Max[T Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// 6. Multiple params
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// 7. Reverse slice
type Slice[T any] interface{ ~[]T }

func Reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// 8. Generic Set
type Set[T comparable] struct {
	m map[T]struct{}
}

func NewSet[T comparable](vals ...T) *Set[T] {
	s := &Set[T]{m: make(map[T]struct{})}
	for _, v := range vals {
		s.m[v] = struct{}{}
	}
	return s
}

func (s *Set[T]) Add(v T) { s.m[v] = struct{}{} }

// ========== NEW ADVANCED FEATURES ==========

// 9. Type Sets in Interfaces
type Stringer interface {
	String() string
	~string
}

// 10. Generic Interface Types (Go 1.21+)
type Collection[T any] interface {
	Add(T)
	Get() T
}

// 11. Variadic Generic Types (Go 1.19+)
type Tuple[T1, T2 any] struct {
	First  T1
	Second T2
}

// 12. Generic Embeddings
type RingBuffer[T any] struct {
	Collection[T] // Embed generic interface
	Size           int
}

// 13. Higher-Order Generic Functions
func Map[I ~[]E, E any, F any](s I, f func(E) F) []F {
	result := make([]F, len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}

// joinInts formats []int as "[1,2,3]" without using %v (crashes on device)
func joinInts(s []int) string {
	r := "["
	for i, v := range s {
		if i > 0 {
			r += ","
		}
		r += fmt.Sprint(v)
	}
	return r + "]"
}

// joinStrs formats []string as "[a,b]" without using %v (crashes on device)
func joinStrs(s []string) string {
	r := "["
	for i, v := range s {
		if i > 0 {
			r += ","
		}
		r += v
	}
	return r + "]"
}

// 9. Type Sets: concrete type satisfying Stringer interface
type myString string

func (m myString) String() string { return string(m) }

// Demo function - draws ALL examples to screen (fits Playdate 400x240)
func demoGenerics() {
	// Title
	pd.Graphics.DrawText("Go Generics COMPLETE!", 10, 10)

	// 1. Generic Functions: Print[T], Min[T]
	Print("hello")
	pd.Graphics.DrawText(fmt.Sprintf("Min(5,3)=%d", Min(5, 3)), 10, 30)

	// 2. Generic Stack[T]
	s := Stack[int]{}
	s.Push(10)
	s.Push(20)
	val, _ := s.Pop()
	pd.Graphics.DrawText(fmt.Sprintf("Stack.Pop=%d", val), 10, 50)

	// 3. Constraint interface: SumNumbers[T]
	total := SumNumbers([]int{1, 2, 3})
	pd.Graphics.DrawText(fmt.Sprintf("Sum=%d", total), 10, 70)

	// 4. Generic Queue[T]
	q := Queue[string]{}
	q.Enqueue("PdGo")
	pd.Graphics.DrawText(fmt.Sprintf("Has:%t", q.Contains("PdGo")), 10, 90)

	// 5. Max[T Ordered]
	pd.Graphics.DrawText(fmt.Sprintf("Max(10.5,7)=%.1f", Max(10.5, 7)), 10, 110)

	// 6. Keys[K comparable, V any]
	m := map[string]int{"A": 1, "B": 2}
	keys := Keys(m)
	pd.Graphics.DrawText("Keys:" + joinStrs(keys), 10, 130)

	// 7. Reverse[S ~[]E, E any]
	nums := []int{3, 1, 2}
	Reverse(nums)
	pd.Graphics.DrawText("Rev:" + joinInts(nums), 10, 150)

	// 8. Generic Set[T comparable]
	set := NewSet[string]("A", "B")
	set.Add("C")
	pd.Graphics.DrawText(fmt.Sprintf("Set size=%d", len(set.m)), 10, 170)

	// ========== NEW ADVANCED ==========

	// 9. Type Sets in Interfaces (call String() directly — fmt.Sprint crashes on device with custom Stringer)
	pd.Graphics.DrawText(string(myString("test")), 10, 0)

	// 11. Variadic Generic Types: Tuple[T1, T2]
	t := Tuple[string, int]{First: "Hi", Second: 42}
	pd.Graphics.DrawText(fmt.Sprintf("Tuple:%s/%d", t.First, t.Second), 10, 190)

	// 13. Higher-Order Generic Functions: Map[I ~[]E, E any, F any]
	doubled := Map([]int{1, 2, 3}, func(x int) int { return x * 2 })
	pd.Graphics.DrawText("Map:" + joinInts(doubled), 10, 210)
}

// initGame is called once when the game starts (REPLACEMENT FOR MAIN)
func initGame() {
	demoGenerics()
}

// update is called every frame
func update() int {
	return 1
}

func main() {}
