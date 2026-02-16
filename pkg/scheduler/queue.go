package scheduler

import (
	"fmt"
	"sync"
	"time"

	"github.com/velocity-trinity/core/pkg/logger"
)

// MemoryQueue is a simple in-memory implementation of Queue
type MemoryQueue struct {
	jobs map[string]*Job
	ch   chan *Job
	mu   sync.RWMutex
}

func NewMemoryQueue(bufferSize int) *MemoryQueue {
	return &MemoryQueue{
		jobs: make(map[string]*Job),
		ch:   make(chan *Job, bufferSize),
	}
}

func (q *MemoryQueue) Enqueue(job *Job) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	job.CreatedAt = time.Now()
	job.Status = StatusPending
	
	// If it depends on another PR, log speculation
	if job.BasePR > 0 {
		job.Status = StatusSpeculative
		logger.Log.Info(fmt.Sprintf("Speculatively enqueuing Job %s assuming PR %d passes", job.ID, job.BasePR))
	} else {
		logger.Log.Info("Enqueuing Job " + job.ID)
	}

	q.jobs[job.ID] = job
	q.ch <- job
	return nil
}

func (q *MemoryQueue) Dequeue() (*Job, error) {
	select {
	case job := <-q.ch:
		return job, nil
	case <-time.After(5 * time.Second):
		return nil, nil // Timeout
	}
}

func (q *MemoryQueue) Get(id string) (*Job, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	job, ok := q.jobs[id]
	if !ok {
		return nil, fmt.Errorf("job not found")
	}
	return job, nil
}
