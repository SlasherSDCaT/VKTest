package main

import (
	"VKTest/workerpool"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	pool := workerpool.NewPool(10)

	for i := 0; i < 3; i++ {
		pool.AddWorker()
	}

	taskFunc := func(data string) string {
		return fmt.Sprintf("processed data: %s", data)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter commands: add_task <data>, add_worker, remove_worker, stop")

	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		commands := strings.SplitN(input, " ", 2)

		switch commands[0] {
		case "add_task":
			if len(commands) < 2 {
				fmt.Println("Usage: add_task <data>")
				continue
			}
			task := workerpool.NewTask(taskFunc, commands[1])
			pool.AddTask(task)
		case "add_worker":
			pool.AddWorker()
		case "remove_worker":
			pool.RemoveWorker()
		case "stop":
			pool.Stop()
			fmt.Println("All workers stopped. Exiting.")
			return
		default:
			fmt.Println("Unknown command. Use: add_task, add_worker, remove_worker, stop")
		}
	}
}
