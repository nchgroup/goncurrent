package goncurrent

import (
	"errors"
	"reflect"
	"sync"
)

type ThreadPool struct {
	workerCount int
	taskQueue   chan func()
	wg          sync.WaitGroup
}

func Execute(threads int, list interface{}, task func(interface{})) error {

	listType := reflect.TypeOf(list)
	if listType.Kind() != reflect.Slice {
		return errors.New("list debe ser un slice")
	}

	// Crear un ThreadPool con el número de trabajadores especificado
	pool := NewThreadPool(threads)

	// Iniciar el ThreadPool
	pool.Start()

	// Enviar tareas al ThreadPool
	listValue := reflect.ValueOf(list)
	for i := 0; i < listValue.Len(); i++ {
		item := listValue.Index(i).Interface()
		pool.Submit(func() {
			task(item)
		})
	}
	pool.Shutdown()
	// Esperar a que todas las tareas se completen
	pool.Wait()
	return nil
}

func NewThreadPool(workerCount int) *ThreadPool {
	return &ThreadPool{
		workerCount: workerCount,
		taskQueue:   make(chan func()),
	}
}

func (tp *ThreadPool) Start() {
	for i := 0; i < tp.workerCount; i++ {
		tp.wg.Add(1)
		go func() {
			defer tp.wg.Done()
			for task := range tp.taskQueue {
				task() // Ejecutar la tarea
			}
		}()
	}
}

func (tp *ThreadPool) Submit(task func()) {
	tp.taskQueue <- task
}

func (tp *ThreadPool) Wait() {
	tp.wg.Wait()
}

func (tp *ThreadPool) Shutdown() {
	close(tp.taskQueue)
}
