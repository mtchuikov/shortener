package closer

import (
	"context"
	"slices"
	"sync"
)

const DefaultMaxConcurrent = 5

type Task struct {
	Sync bool
	Fn   func(context.Context)
}

type Closer struct {
	mu            sync.Mutex
	tasks         []Task
	numTasks      int
	closeOnce     sync.Once
	maxConcurrent int
}

func New(opts ...Option) *Closer {
	c := &Closer{
		mu:            sync.Mutex{},
		tasks:         make([]Task, 0, 3),
		closeOnce:     sync.Once{},
		maxConcurrent: DefaultMaxConcurrent,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Closer) Reset() {
	c.mu.Lock()
	c.tasks = make([]Task, 0, 3)
	c.numTasks = 0
	c.closeOnce = sync.Once{}
	c.mu.Unlock()
}

func (c *Closer) Add(task Task) {
	c.mu.Lock()
	c.numTasks++
	c.tasks = append(c.tasks, task)
	c.mu.Unlock()
}

func (c *Closer) AddWithPriority(task Task, priority int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.numTasks++

	if priority <= 0 {
		priority = 0
	}

	if priority >= c.numTasks {
		c.tasks = append(c.tasks, task)
		return
	}

	c.tasks = slices.Insert(c.tasks, priority, task)
}

func (c *Closer) Close(ctx context.Context) error {
	var err error
	closeFn := func() {
		sem := make(chan struct{}, c.maxConcurrent)
		var wg sync.WaitGroup

		for _, task := range c.tasks {
			select {
			case <-ctx.Done():
				err = ctx.Err()
				return
			case sem <- struct{}{}:
			}

			wg.Add(1)

			if task.Sync {
				task.Fn(ctx)
				continue
			}

			go func() {
				defer func() {
					wg.Done()
					<-sem
				}()

				task.Fn(ctx)
			}()
		}

		close(sem)

		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-ctx.Done():
			err = ctx.Err()
			return
		case <-done:
		}
	}

	c.closeOnce.Do(closeFn)

	return err
}
