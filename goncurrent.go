package goncurrent

import (
	"errors"
	"sync"
)

// ThreadPool is a structure that manages a pool of goroutines to execute tasks
type ThreadPool struct {
	workerCount int
	taskQueue   chan func()
	wg          sync.WaitGroup
	closed      chan struct{}
}

// Execute runs a set of tasks using a fixed number of goroutines.
// It takes the number of threads (goroutines) and a slice of tasks (functions) as input.
func Execute(threads int, tasks []func()) error {
	if threads < 1 {
		return errors.New("invalid number of threads")
	}
	if len(tasks) == 0 {
		return nil // Nothing to execute
	}

	// Initialize the ThreadPool
	pool := NewThreadPool(threads)
	pool.Start()

	// Submit all tasks to the ThreadPool
	for _, task := range tasks {
		pool.Submit(task)
	}

	// Close the ThreadPool and wait for all tasks to complete
	pool.Shutdown()
	pool.Wait()
	return nil
}

// NewThreadPool creates a new ThreadPool with the specified number of workers (goroutines).
func NewThreadPool(workerCount int) *ThreadPool {
	return &ThreadPool{
		workerCount: workerCount,
		taskQueue:   make(chan func(), workerCount), // Using a buffered channel
		closed:      make(chan struct{}),
	}
}

// Start begins the execution of the ThreadPool, launching worker goroutines.
func (tp *ThreadPool) Start() {
	for i := 0; i < tp.workerCount; i++ {
		tp.wg.Add(1)
		go func() {
			defer tp.wg.Done()
			for {
				select {
				case task, ok := <-tp.taskQueue:
					if !ok {
						return // Channel closed, exit goroutine
					}
					task()
				case <-tp.closed:
					return // Shutdown signal received, exit goroutine
				}
			}
		}()
	}
}

// Submit adds a new task to the ThreadPool.
// If the ThreadPool is already closed, the task is not added.
func (tp *ThreadPool) Submit(task func()) {
	select {
	case tp.taskQueue <- task:
	case <-tp.closed:
		// ThreadPool is closed, do not accept new tasks
	}
}

// Wait blocks until all tasks have been completed.
func (tp *ThreadPool) Wait() {
	tp.wg.Wait()
}

// Shutdown gracefully shuts down the ThreadPool.
// It stops accepting new tasks and closes the taskQueue channel after all existing tasks are processed.
func (tp *ThreadPool) Shutdown() {
	close(tp.closed) // Signal to goroutines to stop
	close(tp.taskQueue) // Close the task queue
}
