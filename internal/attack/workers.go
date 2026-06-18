package attack

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"

	"github.com/Ullaakut/cameradar/v6"
)

type attackFn func(context.Context, cameradar.Stream) (cameradar.Stream, error)

func runParallel(ctx context.Context, targets []cameradar.Stream, fn attackFn) ([]cameradar.Stream, error) {
	if len(targets) == 0 {
		return targets, nil
	}

	workerCount := parallelWorkerCount(len(targets))
	if workerCount == 0 {
		return targets, nil
	}

	errCh := make(chan error, workerCount)
	jobs := make(chan attackJob)

	updated := make([]cameradar.Stream, len(targets))
	copy(updated, targets)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	for range workerCount {
		wg.Go(func() {
			errCh <- runWorker(ctx, jobs, fn, updated)
		})
	}

	queued := queueJobs(ctx, jobs, targets)
	close(jobs)

	wg.Wait()
	close(errCh)

	// Aggregate every worker's errors instead of surfacing only the first one.
	// A failure on one camera (e.g. a host that masscan flagged but is now
	// unreachable) must not drop results for every other camera in the batch,
	// and all failures are returned together so callers can inspect each.
	var errs error
	for err := range errCh {
		errs = errors.Join(errs, err)
	}

	// When the context is cancelled before all jobs are queued, the unqueued
	// targets never run and their slots in updated hold stale pre-attack data.
	// Surface the cancellation so callers know the results are incomplete.
	if queued < len(targets) {
		errs = errors.Join(errs, fmt.Errorf("attack cancelled after %d/%d targets: %w", queued, len(targets), ctx.Err()))
	}

	return updated, errs
}

type attackJob struct {
	index  int
	stream cameradar.Stream
}

func queueJobs(ctx context.Context, jobs chan<- attackJob, targets []cameradar.Stream) int {
	for i, stream := range targets {
		select {
		case <-ctx.Done():
			return i
		case jobs <- attackJob{index: i, stream: stream}:
		}
	}
	return len(targets)
}

func runWorker(ctx context.Context, jobs <-chan attackJob, fn attackFn, updated []cameradar.Stream) error {
	var errs error
	for {
		select {
		case <-ctx.Done():
			return errs
		case job, ok := <-jobs:
			if !ok {
				return errs
			}

			stream, err := fn(ctx, job.stream)
			updated[job.index] = stream
			if err != nil {
				// Aggregate the error but keep processing the remaining
				// targets. A failure on one camera must not drop results for
				// every other camera in the batch. Cancellation is reserved
				// for a genuine ctx.Done().
				errs = errors.Join(errs, fmt.Errorf("attacking %s: %w", stream, err))
			}
		}
	}
}

func parallelWorkerCount(targetCount int) int {
	if targetCount <= 0 {
		return 0
	}

	workers := max(runtime.GOMAXPROCS(0), 1)
	if targetCount < workers {
		return targetCount
	}

	return workers
}
