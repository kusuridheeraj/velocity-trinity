package scheduler

import (
	"fmt"
	"time"

	"github.com/velocity-trinity/core/pkg/logger"
)

// Worker simulates a CI runner
func RunWorker(id int, queue Queue) {
	logger.Log.Info(fmt.Sprintf("Worker %d started", id))
	for {
		job, err := queue.Dequeue()
		if err != nil {
			logger.Log.Error(fmt.Sprintf("Worker %d dequeue error: %v", id, err))
			continue
		}
		if job == nil {
			time.Sleep(1 * time.Second) // Idle
			continue
		}

		// Simulate Running CI
		logger.Log.Info(fmt.Sprintf("Worker %d running Job %s (Speculating PR %d -> PR %d)", id, job.ID, job.BasePR, job.PRNumber))
		job.Status = StatusRunning
		
		// Simulate CI time
		time.Sleep(2 * time.Second)
		
		// Simulate random success/failure
		if job.PRNumber%2 == 0 {
			job.Status = StatusSuccess
			logger.Log.Info(fmt.Sprintf("Job %s passed!", job.ID))
		} else {
			job.Status = StatusFailed
			logger.Log.Info(fmt.Sprintf("Job %s failed!", job.ID))
		}
	}
}
