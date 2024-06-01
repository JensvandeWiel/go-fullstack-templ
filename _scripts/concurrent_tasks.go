package main

import (
	"log"
	"os"
	"os/exec"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	// Get the tasks from the command-line arguments
	tasks := os.Args[1:]

	// Increment the WaitGroup counter for each task
	wg.Add(len(tasks))

	// Run each task in a separate goroutine
	for _, task := range tasks {
		go func(task string) {
			defer wg.Done() // Decrement the counter when the goroutine completes
			err := runProcess("task", task)
			if err != nil {
				log.Fatalf("Failed to run %s task: %v", task, err)
			}
		}(task)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}

func runProcess(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
