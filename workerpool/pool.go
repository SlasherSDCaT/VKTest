package workerpool

import (
	"fmt"
	"sync"
)

type Pool struct {
	taskChan    chan *Task
	workers     []*Worker
	workerCount int
	mu          sync.Mutex
}

// NewPool создает новый пул с указанным размером канала задач.
func NewPool(bufferSize int) *Pool {
	return &Pool{
		taskChan: make(chan *Task, bufferSize),
	}
}

// AddWorker добавляет нового воркера в пул.
func (p *Pool) AddWorker() {
	p.mu.Lock()
	defer p.mu.Unlock()

	worker := NewWorker(p.workerCount+1, p.taskChan)
	worker.Start()
	p.workers = append(p.workers, worker)
	p.workerCount++
	fmt.Printf("Added worker %d\n", worker.ID)
}

// RemoveWorker удаляет последнего добавленного воркера.
func (p *Pool) RemoveWorker() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.workers) == 0 {
		fmt.Println("No workers to remove.")
		return
	}

	worker := p.workers[len(p.workers)-1]
	worker.Stop()
	p.workers = p.workers[:len(p.workers)-1]
	p.workerCount--
	fmt.Printf("Removed worker %d\n", worker.ID)
}

// AddTask добавляет задачу в канал для обработки воркерами.
func (p *Pool) AddTask(task *Task) {
	p.taskChan <- task
}

// Stop останавливает всех воркеров и закрывает канал задач.
func (p *Pool) Stop() {
	for _, worker := range p.workers {
		worker.Stop()
	}
	close(p.taskChan)
	p.workers = nil
	p.workerCount = 0
}
