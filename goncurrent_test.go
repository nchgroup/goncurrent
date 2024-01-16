package goncurrent

import (
	"sync/atomic"
	"testing"
	"time"
)

// TestThreadPoolInitialization checks if the ThreadPool is initialized correctly.
func TestThreadPoolInitialization(t *testing.T) {
	workerCount := 5
	tp := NewThreadPool(workerCount)

	if tp.workerCount != workerCount {
		t.Errorf("expected workerCount %d, got %d", workerCount, tp.workerCount)
	}
}

// TestThreadPoolTaskExecution verifies that a task is executed.
func TestThreadPoolTaskExecution(t *testing.T) {
	var counter int32
	task := func() {
		atomic.AddInt32(&counter, 1)
	}

	Execute(1, []func(){task})

	if counter != 1 {
		t.Errorf("expected counter to be 1, got %d", counter)
	}
}

// TestThreadPoolConcurrency checks if multiple tasks are executed concurrently.
func TestThreadPoolConcurrency(t *testing.T) {
	var counter int32
	task := func() {
		atomic.AddInt32(&counter, 1)
		time.Sleep(100 * time.Millisecond)
	}

	Execute(5, []func(){task, task, task, task, task})

	if counter != 5 {
		t.Errorf("expected counter to be 5, got %d", counter)
	}
}

// TestThreadPoolShutdown verifies that no new tasks are accepted after shutdown.
func TestThreadPoolShutdown(t *testing.T) {
	tp := NewThreadPool(1)
	tp.Start()
	tp.Shutdown()

	done := make(chan struct{})
	go func() {
		tp.Submit(func() {}) // This should not be executed
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Error("Submit did not return immediately after shutdown")
	}
}

// TestInvalidThreadCount checks the error handling for invalid thread counts.
func TestInvalidThreadCount(t *testing.T) {
	err := Execute(0, []func(){func() {}})
	if err == nil {
		t.Error("expected error for invalid thread count, got nil")
	}
}
