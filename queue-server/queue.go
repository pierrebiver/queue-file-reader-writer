package main

import "sync"

type Queue struct {
	mu    sync.Mutex
	items []string
}

func (q *Queue) Push(item string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, item)
}

func (q *Queue) Pop() (string, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	itemIsEmpty := len(q.items) == 0
	if itemIsEmpty {
		return "", false
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}
