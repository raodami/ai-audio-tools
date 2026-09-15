package queue

import (
	"container/heap"
	"sync"
	"time"
)

// Job represents a processing job in the queue
type Job struct {
	ID         string
	Priority   int
	Duration   time.Duration
	Status     JobStatus
	Result     string
	Progress   int
	Submitted  time.Time
	Started    time.Time
	Completed  time.Time
}

// JobStatus represents the status of a job
type JobStatus int

const (
	StatusPending JobStatus = iota
	StatusProcessing
	StatusCompleted
	StatusFailed
)

// jobHeap implements heap.Interface
type jobHeap []*Job

func (h jobHeap) Len() int { return len(h) }
func (h jobHeap) Less(i, j int) bool {
	return h[i].Priority > h[j].Priority // Higher priority first
}
func (h jobHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *jobHeap) Push(x interface{}) {
	*h = append(*h, x.(*Job))
}

func (h *jobHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*h = old[:n-1]
	return item
}

// JobQueue implements a priority queue for jobs
type JobQueue struct {
	mu    sync.Mutex
	items jobHeap
}

// NewJobQueue creates a new job queue
func NewJobQueue() *JobQueue {
	h := &jobHeap{}
	heap.Init(h)
	return &JobQueue{items: *h}
}

// Enqueue adds a job to the queue
func (q *JobQueue) Enqueue(job *Job) {
	q.mu.Lock()
	defer q.mu.Unlock()

	job.Status = StatusPending
	job.Submitted = time.Now()
	heap.Push(&q.items, job)
}

// Dequeue removes and returns the highest priority job
func (q *JobQueue) Dequeue() *Job {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.items) == 0 {
		return nil
	}

	job := heap.Pop(&q.items).(*Job)
	job.Status = StatusProcessing
	job.Started = time.Now()
	return job
}

// GetStatus returns the current status of a job
func (q *JobQueue) GetStatus(jobID string) (*Job, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	for _, job := range q.items {
		if job.ID == jobID {
			return job, true
		}
	}
	return nil, false
}

// GetQueueLength returns the number of jobs in the queue
func (q *JobQueue) GetQueueLength() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

// GetCompletedCount returns the number of completed jobs
func (q *JobQueue) GetCompletedCount() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	count := 0
	for _, job := range q.items {
		if job.Status == StatusCompleted {
			count++
		}
	}
	return count
}

// GetFailedCount returns the number of failed jobs
func (q *JobQueue) GetFailedCount() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	count := 0
	for _, job := range q.items {
		if job.Status == StatusFailed {
			count++
		}
	}
	return count
}

// GetPendingCount returns the number of pending jobs
func (q *JobQueue) GetPendingCount() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	count := 0
	for _, job := range q.items {
		if job.Status == StatusPending || job.Status == StatusProcessing {
			count++
		}
	}
	return count
}

// GetAvgDuration returns the average processing duration
func (q *JobQueue) GetAvgDuration() float64 {
	q.mu.Lock()
	defer q.mu.Unlock()

	var totalDuration time.Duration
	count := 0

	for _, job := range q.items {
		if job.Status == StatusCompleted && !job.Started.IsZero() && !job.Completed.IsZero() {
			totalDuration += job.Completed.Sub(job.Started)
			count++
		}
	}

	if count == 0 {
		return 0
	}

	return float64(totalDuration) / float64(count)
}
