package scheduler

// Run starts the scheduler and workers
func Run(numWorkers int) {
	queue := NewMemoryQueue(100)

	// Start Workers
	for i := 0; i < numWorkers; i++ {
		go RunWorker(i, queue)
	}

	// Simulate Job Submissions
	// In a real system, this would be a webhook listener
	go func() {
		// Job 1 (PR 1, Base main)
		queue.Enqueue(&Job{ID: "job-1", PRNumber: 1, BasePR: 0})
		
		// Job 2 (PR 2, Base 1)
		// This is the "Quantum" part: Running PR 2 assuming PR 1 passes
		queue.Enqueue(&Job{ID: "job-2", PRNumber: 2, BasePR: 1})
		
		// Job 3 (PR 3, Base 2)
		queue.Enqueue(&Job{ID: "job-3", PRNumber: 3, BasePR: 2})
	}()

	// Keep main process alive
	select {}
}
