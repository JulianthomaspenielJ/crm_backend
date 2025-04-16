package queue

type Queue interface {
	Publish(task Task) error
	Consume() error
}

type Task func() error
