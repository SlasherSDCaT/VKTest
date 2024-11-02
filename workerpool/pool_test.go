package workerpool

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// Мок-функция обработки задач для теста, добавляющая приписку "processed" к данным.
func mockTaskFunc(data string) string {
	return "processed: " + data
}

func TestWorkerPool_AddWorkerAndProcessTasks(t *testing.T) {
	pool := NewPool(5)

	// Добавляем 3 воркера в пул.
	for i := 0; i < 3; i++ {
		pool.AddWorker()
	}

	// Создаем wait group, чтобы дождаться выполнения всех задач.
	var wg sync.WaitGroup

	// Добавляем задачи и увеличиваем счетчик wait group.
	for i := 0; i < 5; i++ {
		wg.Add(1)
		taskData := fmt.Sprintf("Task data %d", i+1) // Используем fmt.Sprintf для преобразования числа в строку.

		// Создаем задачу с использованием mock-функции.
		task := NewTask(func(data string) string {
			defer wg.Done()
			return mockTaskFunc(data)
		}, taskData)

		pool.AddTask(task)
	}

	// Ожидаем выполнения всех задач.
	wg.Wait()

	// Проверяем, что у нас есть 3 воркера после добавления.
	if len(pool.workers) != 3 {
		t.Errorf("Expected 3 workers, but got %d", len(pool.workers))
	}

	// Останавливаем все воркеры и закрываем пул.
	pool.Stop()
	if len(pool.workers) != 0 {
		t.Errorf("Expected 0 workers after stop, but got %d", len(pool.workers))
	}
}

func TestWorkerPool_AddAndRemoveWorkers(t *testing.T) {
	pool := NewPool(5)

	// Добавляем 2 воркера.
	pool.AddWorker()
	pool.AddWorker()

	// Проверяем, что количество воркеров равно 2.
	if len(pool.workers) != 2 {
		t.Errorf("Expected 2 workers, but got %d", len(pool.workers))
	}

	// Удаляем одного воркера.
	pool.RemoveWorker()

	// Проверяем, что количество воркеров теперь 1.
	if len(pool.workers) != 1 {
		t.Errorf("Expected 1 worker, but got %d", len(pool.workers))
	}

	// Останавливаем пул.
	pool.Stop()
	if len(pool.workers) != 0 {
		t.Errorf("Expected 0 workers after stop, but got %d", len(pool.workers))
	}
}

func TestWorkerPool_Stop(t *testing.T) {
	pool := NewPool(5)

	// Добавляем воркера и задачу, чтобы проверить остановку.
	pool.AddWorker()
	task := NewTask(mockTaskFunc, "test data")
	pool.AddTask(task)

	// Останавливаем пул.
	pool.Stop()

	// Проверяем, что все воркеры остановлены.
	if len(pool.workers) != 0 {
		t.Errorf("Expected 0 workers after stop, but got %d", len(pool.workers))
	}

	// Попытка добавить новую задачу после остановки.
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when adding task to stopped pool, but did not panic")
		}
	}()
	pool.AddTask(NewTask(mockTaskFunc, "new task"))
}

func TestWorkerPool_TaskProcessing(t *testing.T) {
	pool := NewPool(5)
	pool.AddWorker()

	// Счетчик для проверки выполнения задачи.
	var taskProcessed bool
	task := NewTask(func(data string) string {
		taskProcessed = true
		return mockTaskFunc(data)
	}, "sample task")

	// Добавляем задачу и немного ждем выполнения.
	pool.AddTask(task)
	time.Sleep(100 * time.Millisecond)

	if !taskProcessed {
		t.Errorf("Expected task to be processed, but it was not")
	}

	// Завершаем работу пула.
	pool.Stop()
}

func TestWorkerPool_HeavyLoad(t *testing.T) {
	taskCount := 1000    // Количество задач
	workerCount := 10    // Количество воркеров
	pool := NewPool(100) // Буферизованный канал для задач

	// Добавляем заданное количество воркеров
	for i := 0; i < workerCount; i++ {
		pool.AddWorker()
	}

	// Счетчик обработанных задач
	var processedTasks int
	var mu sync.Mutex     // Мьютекс для безопасного доступа к счетчику
	var wg sync.WaitGroup // WaitGroup для ожидания завершения всех задач

	// Функция обработки задачи с увеличением счетчика
	taskFunc := func(data string) string {
		mu.Lock()
		processedTasks++
		mu.Unlock()
		return "processed: " + data
	}

	// Добавляем задачи и увеличиваем счетчик wait group
	for i := 0; i < taskCount; i++ {
		wg.Add(1)
		taskData := fmt.Sprintf("Task %d", i+1)

		// Создаем задачу с использованием taskFunc
		task := NewTask(func(data string) string {
			defer wg.Done()
			return taskFunc(data)
		}, taskData)

		pool.AddTask(task)
	}

	// Ожидаем выполнения всех задач
	wg.Wait()

	// Проверяем, что все задачи были обработаны
	if processedTasks != taskCount {
		t.Errorf("Expected %d processed tasks, but got %d", taskCount, processedTasks)
	}

	// Останавливаем пул и проверяем, что все воркеры завершили выполнение
	pool.Stop()
	if len(pool.workers) != 0 {
		t.Errorf("Expected 0 workers after stop, but got %d", len(pool.workers))
	}
}
