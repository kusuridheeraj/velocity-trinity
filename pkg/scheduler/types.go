package scheduler

import "time"

type JobStatus string

const (
	StatusPending     JobStatus = "PENDING"
	StatusRunning     JobStatus = "RUNNING"
	StatusSuccess     JobStatus = "SUCCESS"
	StatusFailed      JobStatus = "FAILED"
	StatusSpeculative JobStatus = "SPECULATIVE"
)

type Job struct {
	ID        string
	PRNumber  int
	BasePR    int // 0 if based on main
	Status    JobStatus
	CreatedAt time.Time
	// In a real system, we'd store the Git SHA, Repo URL, etc.
}

type Queue interface {
	Enqueue(job *Job) error
	Dequeue() (*Job, error)
	Get(id string) (*Job, error)
}
