package main

type Queue struct {
	items chan string
}

func NewQueue(capacity int) *Queue {
	return &Queue{items: make(chan string, capacity)}
}

func (q *Queue) Push(item string) {
	q.items <- item
}

func (q *Queue) Pop() string {
	return <-q.items
}
