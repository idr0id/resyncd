package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/fsnotify/fsevents"
)

const fseventsLatency = 500 * time.Millisecond

func watchDirectory(
	ctx context.Context,
	logger *slog.Logger,
	conf configSync,
	filesCh chan<- string,
) error {
	defer close(filesCh)

	source := endsWithSlash(conf.Source)

	logger = logger.With(slog.String("source", source))
	logger.Debug("start watching directory")

	deviceID, err := fsevents.DeviceForPath(source)
	if err != nil {
		return fmt.Errorf("failed to retrieve device for path: %s: %w", source, err)
	}

	stream := &fsevents.EventStream{
		Paths:   []string{source},
		Latency: fseventsLatency,
		Device:  deviceID,
		Flags:   fsevents.FileEvents | fsevents.WatchRoot,
	}
	if err = stream.Start(); err != nil {
		return fmt.Errorf("failed to start event stream for path: %s: %w", source, err)
	}
	defer stream.Stop()

	excludes := newFileMatchers(source, conf.Exclude)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case events := <-stream.Events:
			for _, event := range events {
				path := "/" + event.Path
				logger := logger.With(slog.String("path", path))

				if isTemporaryEvent(event) {
					logger.Debug("skipping temporary file")
					continue
				}

				if excludes.match("/" + event.Path) {
					logger.Debug("skipping excluded file")
					continue
				}

				logger.Debug("detected changed file")
				filesCh <- path[len(source):]
			}
		}
	}
}

func isTemporaryEvent(event fsevents.Event) bool {
	bits := fsevents.ItemCreated + fsevents.ItemRemoved
	return event.Flags&bits == bits
}
