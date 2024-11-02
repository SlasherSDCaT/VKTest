package workerpool

import "fmt"

// Worker представляет воркера, который обрабатывает задачи.
type Worker struct {
	ID       int
	taskChan chan *Task
	quit     chan bool
}

// NewWorker создаёт нового воркера с заданным каналом задач.
func NewWorker(id int, taskChan chan *Task) *Worker {
	return &Worker{
		ID:       id,
		taskChan: taskChan,
		quit:     make(chan bool),
	}
}

// Start запускает воркера, который ждёт задачи на выполнение.
func (w *Worker) Start() {
	go func() {
		for {
			select {
			case task := <-w.taskChan:
				task.Process(w.ID)
			case <-w.quit:
				fmt.Printf("Worker %d stopping.\n", w.ID)
				return
			}
		}
	}()
}

// Stop отправляет сигнал для остановки воркера.
func (w *Worker) Stop() {
	w.quit <- true
}
