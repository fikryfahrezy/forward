package worker

import (
	"fmt"
	"io"
)

type jobRunner func() any

type job struct {
	runner jobRunner
	result chan any
}

type wrk struct {
	jobs   chan job
	closed chan struct{}
	logger io.Writer
}

type option func(*wrk)

var ErrPoolClosed = fmt.Errorf("worker pool is closed")

func WithLogger(w io.Writer) option {
	return func(wrk *wrk) {
		wrk.logger = w
	}
}

func New(pool int, opts ...option) *wrk {
	worker := &wrk{
		jobs:   make(chan job),
		closed: make(chan struct{}),
		logger: io.Discard,
	}

	for _, opt := range opts {
		opt(worker)
	}

	for i := range pool {
		go worker.work(i)
	}

	return worker
}

func (w *wrk) work(id int) {
	for {
		select {
		case job, ok := <-w.jobs:
			if !ok {
				fmt.Fprintf(w.logger, "Worker %d: jobs channel closed, exiting\n", id)
				return
			}
			fmt.Fprintf(w.logger, "Worker %d: executing job\n", id)
			result := job.runner()

			select {
			case job.result <- result:
				fmt.Fprintf(w.logger, "Worker %d: result sent successfully\n", id)
			case <-w.closed:
				// Pool is shutting down while we're trying to send result
				// Exit immediately instead of blocking forever
				// (This happens if Close() is called but no one reads the result)
				fmt.Fprintf(w.logger, "Worker %d: shutdown signal received while sending result, exiting\n", id)
				return
			}

		case <-w.closed:
			// Pool is shutting down while we're idle (waiting for jobs)
			// Exit immediately instead of waiting forever for work
			fmt.Fprintf(w.logger, "Worker %d: shutdown signal received while idle, exiting\n", id)
			return
		}
	}
}

func (w *wrk) Add(jobRunner jobRunner) (<-chan any, error) {
	result := make(chan any, 1)

	select {
	case w.jobs <- job{runner: jobRunner, result: result}:
		fmt.Fprintf(w.logger, "Job queued successfully\n")
		return result, nil
	case <-w.closed:
		close(result)
		fmt.Fprintf(w.logger, "Failed to add job: pool is closed\n")
		return result, ErrPoolClosed
	}
}

func (w *wrk) Close() {
	fmt.Fprintf(w.logger, "Closing worker pool\n")
	close(w.closed)
	close(w.jobs)
}
