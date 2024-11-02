package workerpool

import (
	"fmt"
)

// Задачи для воркера
type Task struct {
	Data string
	f    func(string) string
}

// Конструктор новой задачи
func NewTask(f func(string) string, data string) *Task {
	return &Task{
		f:    f,
		Data: data,
	}
}

// Выполнение задачи
func (t *Task) Process(workerID int) {
	result := t.f(t.Data)
	fmt.Printf("Worker %d processed task: %s\n", workerID, result)
}
