package main

import "github.com/reconquest/karma-go"

type synchronizer struct {
	rsync   *rsync
	watcher *watcher
}

func newSynchronizer() *synchronizer {
	return &synchronizer{
		newRsync(),
		newWatcher(),
	}
}

func (s *synchronizer) start(cfg configSync) {
	ctx := karma.Describe("source", cfg.Source).
		Describe("target", cfg.Target)
	logger.Debugf(ctx, "starting synchronizer")

	syncChan := make(chan string)
	s.rsync.start(cfg, syncChan)
	s.watcher.start(cfg, syncChan)
}

func (s *synchronizer) stop() {
	s.rsync.stop()
	s.watcher.stop()
}