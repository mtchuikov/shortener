package closer

import (
	"context"
	"slices"
	"sync"
)

const DefaultMaxConcurrent = 5

var c *closer = nil

type Task struct {
	Sync bool
	Fn   func(context.Context)
}

type closer struct {
	mu            sync.Mutex
	tasks         []Task
	numTasks      int
	closeOnce     sync.Once
	maxConcurrent int
}

func Init(opts ...Option) {
	c = &closer{
		mu:            sync.Mutex{},
		tasks:         make([]Task, 0, 3),
		closeOnce:     sync.Once{},
		maxConcurrent: DefaultMaxConcurrent,
	}

	for _, opt := range opts {
		opt(c)
	}

}

func Reset() {
	c.mu.Lock()
	c.tasks = make([]Task, 0, 3)
	c.numTasks = 0
	c.closeOnce = sync.Once{}
	c.mu.Unlock()
}

func Add(task Task) {
	c.mu.Lock()
	c.numTasks++
	c.tasks = append(c.tasks, task)
	c.mu.Unlock()
}

func AddWithPriority(priority int, task Task) {
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

func Close(ctx context.Context) error {
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
			doneFn := func() {
				wg.Done()
				<-sem
			}

			if task.Sync {
				task.Fn(ctx)
				doneFn()
				continue
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

	c.closeOnce.Do(closeFn)

	return err
}
