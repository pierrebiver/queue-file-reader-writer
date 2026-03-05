package main

import (
	"sync"
	"testing"
)

func TestPushPop(t *testing.T) {
	q := NewQueue(1)
	q.Push("hello")
	val := q.Pop()
	if val != "hello" {
		t.Fatalf("expected %q, got %q", "hello", val)
	}
}

func TestFIFOOrder(t *testing.T) {
	items := []string{"first", "second", "third"}
	q := NewQueue(len(items))
	for _, item := range items {
		q.Push(item)
	}
	for _, want := range items {
		got := q.Pop()
		if got != want {
			t.Fatalf("expected %q, got %q", want, got)
		}
	}
}

func TestConcurrentPushPop(t *testing.T) {
	const n = 10
	q := NewQueue(n)
	var wg sync.WaitGroup
	wg.Add(n * 2)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			q.Push("item")
		}()
	}

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			q.Pop()
		}()
	}

	wg.Wait()
}
