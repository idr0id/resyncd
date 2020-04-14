package main

import (
	"github.com/reconquest/karma-go"
	"github.com/zloylos/grsync"
	"sync"
	"time"
)

type rsync struct {
	done chan struct{}
	wg   sync.WaitGroup
}

func newRsync() *rsync {
	return &rsync{
		done: make(chan struct{}),
	}
}

func (r *rsync) start(cfg configSync, syncChan <-chan string) {
	r.wg.Add(1)
	defer r.wg.Done()

	source := cfg.Source.String()
	target := cfg.Target.String()

	ctx := karma.Describe("source", source).Describe("target", target)
	logger.Infof(ctx, "initial synchronization started")

	task := grsync.NewTask(source, target, grsync.RsyncOptions{
		Rsh:     cfg.Rsync.Rsh,
		ACLs:    cfg.Rsync.ACLs,
		Perms:   cfg.Rsync.Perms,
		Exclude: cfg.Exclude,
		Stats:   true,
		Delete:  true,
	})
	if err := task.Run(); err != nil {
		logger.Fatalf(logTask(ctx, task).Reason(err), "synchronization failed")
	}
	logger.Infof(logTask(ctx, task), "initial synchronization completed")

	go func() {
		r.wg.Add(1)
		defer r.wg.Done()

		buffer := make([]string, 0)
		mu := sync.Mutex{}
		t := time.NewTicker(200 * time.Millisecond)

		for {
			select {
			case file := <-syncChan:
				mu.Lock()
				buffer = append(buffer, file)
				mu.Unlock()

			case <-t.C:
				if len(buffer) == 0 {
					continue
				}

				mu.Lock()
				changed := buffer
				buffer = []string{}
				mu.Unlock()

				go func() {
					r.wg.Add(1)
					defer r.wg.Done()

					ctx := karma.Describe("source", source).
						Describe("target", target).
						Describe("files", changed)

					logger.Debugf(ctx, "synchronization started")
					task := grsync.NewTask(
						source,
						target,
						grsync.RsyncOptions{
							Rsh:          cfg.Rsync.Rsh,
							ACLs:         cfg.Rsync.ACLs,
							Perms:        cfg.Rsync.Perms,
							Include:      expandPaths(changed),
							Exclude:      []string{"*"},
							Progress:     true,
							Stats:        false,
							Verbose:      true,
							Recursive:    true,
							Delete:       true,
							IgnoreErrors: true,
							Force:        true,
						})
					if err := task.Run(); err != nil {
						logger.Errorf(
							logTask(ctx, task).Reason(err),
							"synchronization failed")
					}
					logger.Infof(logTask(ctx, task), "synchronization complete")
				}()

			case <-r.done:
				return
			}
		}
	}()
}

func logTask(ctx *karma.Context, task *grsync.Task) *karma.Context {
	return ctx.
		Describe("stdout", task.Log().Stdout).
		Describe("stderr", task.Log().Stderr)
}

func (r *rsync) stop() {
	logger.Debugf(nil, "stopping rsync")
	close(r.done)
	r.wg.Wait()
}
