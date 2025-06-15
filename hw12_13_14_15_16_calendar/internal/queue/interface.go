package queue

import "io"

type Handler interface {
	Stop() error
	InitQueue(name string) (io.Closer, error)
	Publish(queue string, data []byte) error
	Consume(queue string) (<-chan []byte, error)
}
