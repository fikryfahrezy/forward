package worker_test

import (
	"testing"
	"time"

	"github.com/fikryfahrezy/forward/csv-processing/caster"
	"github.com/fikryfahrezy/forward/csv-processing/worker"
)

func TestWorkerPool_BasicExecution(t *testing.T) {
	pool := worker.New(3)
	defer pool.Close()

	// Add a simple job
	result, err := pool.Add(func() any {
		return 42
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Wait for result with timeout
	select {
	case got := <-result:
		if got != 42 {
			t.Errorf("expected 42, got %d", got)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for result")
	}
}

func TestWorkerPool_MultipleJobs(t *testing.T) {
	pool := worker.New(3)
	defer pool.Close()

	numJobs := 10
	results := make([]<-chan int, numJobs)

	// Add multiple jobs
	for i := range numJobs {
		jobNum := i
		result, err := pool.Add(func() any {
			return jobNum * 2
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		results[i] = caster.ChanType[int](result)
	}

	// Collect all results
	for i, resultChan := range results {
		select {
		case got := <-resultChan:
			expected := i * 2
			if got != expected {
				t.Errorf("job %d: expected %d, got %d", i, expected, got)
			}
		case <-time.After(time.Second):
			t.Fatalf("timeout waiting for job %d", i)
		}
	}
}
