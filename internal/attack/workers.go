package attack

import (
	"context"
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

	errCh := make(chan error, 1)
	jobs := make(chan attackJob)

	updated := make([]cameradar.Stream, len(targets))
	copy(updated, targets)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	for range workerCount {
		wg.Go(func() {
			runWorker(ctx, jobs, cancel, fn, updated, errCh)
		})
	}

	queueJobs(ctx, jobs, targets)
	close(jobs)

	wg.Wait()

	select {
	case err := <-errCh:
		return updated, err
	default:
	}

	return updated, nil
}

type attackJob struct {
	index  int
	stream cameradar.Stream
}

func queueJobs(ctx context.Context, jobs chan<- attackJob, targets []cameradar.Stream) {
	for i, stream := range targets {
		select {
		case <-ctx.Done():
			return
		case jobs <- attackJob{index: i, stream: stream}:
		}
	}
}

func runWorker(ctx context.Context, jobs <-chan attackJob, cancelFn func(), fn attackFn, updated []cameradar.Stream, errCh chan error) {
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}

			stream, err := fn(ctx, job.stream)
			if err != nil {
				select {
				case errCh <- err:
				default:
				}

				cancelFn()
				return
			}

			updated[job.index] = stream
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
