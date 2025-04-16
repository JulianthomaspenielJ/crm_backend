package queue

import (
	"fmt"
	"sync"
)

type MemoryQueue struct {
	queue chan Task
	wg    sync.WaitGroup
}

func NewMemoryQueue() *MemoryQueue {
	return &MemoryQueue{
		queue: make(chan Task, 100),
	}
}

func (q *MemoryQueue) Publish(task Task) error {
	q.queue <- task
	return nil
}
func (q *MemoryQueue) Consume() error {
	for task := range q.queue {
		q.wg.Add(1)
		go func(t Task) {
			defer q.wg.Done()
			if err := t(); err != nil {
				fmt.Println("Error executing task:", err)
			}
		}(task)
	}
	q.wg.Wait()
	return nil
}
