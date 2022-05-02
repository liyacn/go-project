package worker

import "sync"

// ParallelTask 持续不间断可多开的常驻任务组
type ParallelTask struct {
	stop chan struct{}
	wg   sync.WaitGroup
}

func NewParallelTask(concurrency int, h HandlerFunc) *ParallelTask {
	fn := Wrap(h)
	t := &ParallelTask{stop: make(chan struct{})}
	for range concurrency {
		t.wg.Go(func() {
			for {
				select {
				case <-t.stop:
					return
				default:
					fn()
				}
			}
		})
	}
	return t
}

func (t *ParallelTask) Stop() {
	close(t.stop)
	t.wg.Wait()
}
