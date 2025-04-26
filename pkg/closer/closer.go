package closer

import (
	"context"
	"sync"
)

const DefaultMaxConcurrentTasks = 5

type TaskFn func(context.Context)

type Task struct {
	Sync bool
	Fn   TaskFn
}

type Closer struct {
	mu            sync.Mutex
	VeryFirst     TaskFn
	tasks         []Task
	VeryLast      TaskFn
	numTasks      int
	maxConcurrent int
}

func New(opts ...Option) *Closer {
	c := &Closer{
		mu:            sync.Mutex{},
		tasks:         make([]Task, 0, 3),
		maxConcurrent: DefaultMaxConcurrentTasks,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Closer) NumTasks() int {
	return c.numTasks
}

func (c *Closer) Add(sync bool, task TaskFn) {
	c.mu.Lock()
	c.numTasks++
	c.tasks = append(c.tasks, Task{Sync: sync, Fn: task})
	c.mu.Unlock()
}

func (c *Closer) Reset() {
	c.mu.Lock()
	c.tasks = make([]Task, 0, 3)
	c.numTasks = 0
	c.mu.Unlock()
}

func (c *Closer) Close(ctx context.Context) error {
	var err error
	closeFn := func() {
		if c.VeryFirst != nil {
			c.VeryFirst(ctx)
		}

		if c.VeryLast != nil {
			defer c.VeryLast(ctx)
		}

		sem := make(chan struct{}, c.maxConcurrent)
		var wg sync.WaitGroup

		for _, task := range c.tasks {
			select {
			case <-ctx.Done():
				err = ctx.Err()
				return
			default:
			}

			if task.Sync {
				// Wait until all async tasks will be done before starting the sync
				// ones. It helps to avoid unexpected errors that would be possible
				// whitout this check. For instance, we planned to close the
				// database asynchronously and close the logger synchronously, the
				// logger might be closed before the database finishes its closure,
				// resulting in the loss of log messages generated during the
				// database shutdown.
				wg.Wait()
				task.Fn(ctx)
				continue
			}

			sem <- struct{}{}
			wg.Add(1)

			doneFn := func() {
				wg.Done()
				<-sem
			}

			go func() {
				task.Fn(ctx)
				doneFn()
			}()
		}

		close(sem)

		waitTillDone := make(chan struct{})
		go func() {
			wg.Wait()
			close(waitTillDone)
		}()

		select {
		case <-ctx.Done():
			err = ctx.Err()
			return
		case <-waitTillDone:
		}
	}

	once := sync.Once{}
	once.Do(closeFn)

	return err
}
