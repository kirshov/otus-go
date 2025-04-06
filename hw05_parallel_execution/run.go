package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var countErrors int32
	var once sync.Once
	wg := sync.WaitGroup{}
	ch := make(chan Task, n)
	done := make(chan struct{})

	for i := 0; i < n; i++ {
		wg.Add(1)

		go handleTask(&wg, ch, done, m, &once, &countErrors)
	}

	go func() {
		defer close(ch)
		for _, task := range tasks {
			select {
			case ch <- task:
			case <-done:
				return
			}
		}
	}()

	wg.Wait()
	if m > 0 && int(atomic.LoadInt32(&countErrors)) >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func handleTask(wg *sync.WaitGroup, ch chan Task, done chan struct{}, m int, once *sync.Once, countErrors *int32) {
	defer wg.Done()
	for {
		select {
		case task, ok := <-ch:
			if !ok {
				return
			}

			if err := task(); err != nil {
				if m > 0 && int(atomic.AddInt32(countErrors, 1)) >= m {
					once.Do(func() { close(done) })
					return
				}
			}
		case <-done:
			return
		}
	}
}
