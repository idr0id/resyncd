package main

import (
	"log/slog"
	"slices"
	"time"

	"github.com/zloylos/grsync"
)

const (
	bufferizeDuration = 200 * time.Millisecond
)

func rsyncDirectory(
	logger *slog.Logger,
	conf configSync,
	filesCh <-chan []string,
) {
	source := endsWithSlash(conf.Source)
	target := endsWithSlash(conf.Target)

	logger = logger.With(slog.String("source", source), slog.String("target", target))
	logger.Info("initial synchronization")

	rsyncCh := make(chan *grsync.Task, 1)
	rsyncCh <- grsync.NewTask(source, target, newRsyncDirectoryOptions(conf))

	go func() {
		for files := range filesCh {
			rsyncCh <- grsync.NewTask(source, target, newRsyncFilesOptions(conf, files))
		}
		close(rsyncCh)
	}()

	for task := range rsyncCh {
		err := task.Run()

		state := task.State()
		logger := logger.With(
			slog.Int("remain", state.Remain),
			slog.Int("total", state.Total),
			slog.Float64("progress", state.Progress),
			slog.String("speed", state.Speed),
		)

		if err != nil {
			logger.Error(
				"synchronization failed",
				slog.Any("err", err),
				slog.String("strerr", task.Log().Stderr),
			)
		} else {
			logger.Info("synchronization complete")
		}
	}
}

func bufferize(filesCh <-chan string, d time.Duration) <-chan []string {
	bufferedCh := make(chan []string, 1)

	go func() {
		defer close(bufferedCh)

		buffer := make([]string, 0)
		flush := func() {
			if len(buffer) > 0 {
				bufferedCh <- slices.Clone(buffer)
				buffer = buffer[:0]
			}
		}

		ticker := time.NewTicker(d)
		defer ticker.Stop()

		for {
			select {
			case file, ok := <-filesCh:
				if !ok {
					flush()
					return
				}
				buffer = append(buffer, file)

			case <-ticker.C:
				flush()
			}
		}
	}()

	return bufferedCh
}

func newRsyncFilesOptions(conf configSync, files []string) grsync.RsyncOptions {
	return grsync.RsyncOptions{
		Rsh:          conf.Rsync.Rsh,
		ACLs:         conf.Rsync.ACLs,
		Perms:        conf.Rsync.Perms,
		Include:      expandPaths(files),
		Exclude:      []string{"*"},
		Contimeout:   conf.Rsync.getConnectTimeoutSeconds(),
		Timeout:      conf.Rsync.getTimeoutSeconds(),
		Progress:     true,
		Stats:        false,
		Verbose:      true,
		Recursive:    true,
		Delete:       true,
		IgnoreErrors: true,
		Force:        true,
	}
}

func newRsyncDirectoryOptions(conf configSync) grsync.RsyncOptions {
	return grsync.RsyncOptions{
		Rsh:     conf.Rsync.Rsh,
		ACLs:    conf.Rsync.ACLs,
		Perms:   conf.Rsync.Perms,
		Exclude: conf.Exclude,
		Timeout: rsyncDefaultTimeoutSeconds,
		Stats:   true,
		Delete:  true,
	}
}
