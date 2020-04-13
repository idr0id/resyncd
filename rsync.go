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

	source := cfg.Source
	target := cfg.Target

	if source != "" && source[len(source)-1] != '/' {
		source += "/"
	}
	if target != "" && target[len(target)-1] != '/' {
		target += "/"
	}

	ctx := karma.Describe("source", source).Describe("target", target)
	logger.Infof(ctx, "initial synchronization started")

	err := r.sync(source, target, grsync.RsyncOptions{
		Rsh:     cfg.Rsync.Rsh,
		ACLs:    cfg.Rsync.ACLs,
		Perms:   cfg.Rsync.Perms,
		Exclude: cfg.Exclude,
		Stats:   true,
		Delete:  true,
	})
	if err != nil {
		logger.Fatalf(ctx.Reason(err), "synchronization failed")
	}
	logger.Infof(ctx, "initial synchronization completed")

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

					err := r.sync(
						source,
						target,
						grsync.RsyncOptions{
							Rsh:          cfg.Rsync.Rsh,
							ACLs:         cfg.Rsync.ACLs,
							Perms:        cfg.Rsync.Perms,
							Include:      changed,
							Exclude:      cfg.Exclude,
							Progress:     true,
							Stats:        false,
							Verbose:      true,
							Recursive:    true,
							Delete:       true,
							IgnoreErrors: true,
							Force:        true,
						})

					if err != nil {
						logger.Errorf(ctx.Reason(err), "synchronization failed")
					}
					logger.Infof(ctx, "synchronization complete")
				}()

			case <-r.done:
				return
			}
		}
	}()
}

func (r *rsync) sync(source, target string, options grsync.RsyncOptions) error {
	task := grsync.NewTask(source, target, options)
	if err := task.Run(); err != nil {
		ctx := karma.
			Describe("source", source).
			Describe("target", target).
			Describe("stdout", task.Log().Stdout).
			Describe("stderr", task.Log().Stderr)
		return ctx.Reason(err)
	}
	return nil
}

func (r *rsync) stop() {
	logger.Info("stopping rsync")
	close(r.done)
	r.wg.Wait()
}
