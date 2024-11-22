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
	tasksCount := len(tasks)

	if tasksCount == 0 || n <= 0 {
		return nil
	}

	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	var (
		countWorkers  int
		errorsCount   int64
		limitExceeded bool

		chTask = make(chan Task)
		wg     = new(sync.WaitGroup)
	)

	if tasksCount <= n {
		countWorkers = tasksCount
	} else {
		countWorkers = n
	}

	wg.Add(countWorkers)
	for i := 0; i < countWorkers; i++ {
		go func() {
			for task := range chTask {
				err := task()
				if err != nil {
					atomic.AddInt64(&errorsCount, 1)
				}
			}
			wg.Done()
		}()
	}

	go func() {
		defer close(chTask)

		for i := 0; i < len(tasks); i++ {
			if atomic.LoadInt64(&errorsCount) >= int64(m) {
				limitExceeded = true
				return
			}

			chTask <- tasks[i]
		}
	}()

	wg.Wait()

	if limitExceeded {
		return ErrErrorsLimitExceeded
	}

	return nil
}
