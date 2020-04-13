package main

import (
	"github.com/fsnotify/fsevents"
	"github.com/reconquest/karma-go"
	"sync"
	"time"
)

var flagsDescription = map[fsevents.EventFlags]string{
	fsevents.MustScanSubDirs:   "MustScanSubdirs",
	fsevents.UserDropped:       "UserDropped",
	fsevents.KernelDropped:     "KernelDropped",
	fsevents.EventIDsWrapped:   "EventIDsWrapped",
	fsevents.HistoryDone:       "HistoryDone",
	fsevents.RootChanged:       "RootChanged",
	fsevents.Mount:             "Mount",
	fsevents.Unmount:           "Unmount",
	fsevents.ItemRemoved:       "Removed",
	fsevents.ItemCreated:       "Created",
	fsevents.ItemInodeMetaMod:  "InodeMetaMod",
	fsevents.ItemRenamed:       "Renamed",
	fsevents.ItemModified:      "Modified",
	fsevents.ItemFinderInfoMod: "FinderInfoMod",
	fsevents.ItemChangeOwner:   "ChangeOwner",
	fsevents.ItemXattrMod:      "XAttrMod",
	fsevents.ItemIsFile:        "IsFile",
	fsevents.ItemIsDir:         "IsDir",
	fsevents.ItemIsSymlink:     "IsSymLink",
}

type watcher struct {
	done chan struct{}
	wg   sync.WaitGroup
}

func newWatcher() *watcher {
	return &watcher{
		done: make(chan struct{}, 0),
	}
}

func (s *watcher) start(
	cfg configSync,
	syncChan chan<- string,
) {
	source := cfg.Source.String()
	ctx := karma.Describe("path", source)
	logger.Infof(ctx, "watching directory")

	go func() {
		s.wg.Add(1)
		defer s.wg.Done()

		excludes := newFileMatchers(source, cfg.Exclude)

		deviceID, err := fsevents.DeviceForPath(source)
		if err != nil {
			logger.Fatalf(ctx.Reason(err), "failed to retrieve device for path")
		}
		stream := &fsevents.EventStream{
			Paths:   []string{source},
			Latency: 500 * time.Millisecond,
			Device:  deviceID,
			Flags:   fsevents.FileEvents | fsevents.WatchRoot,
		}
		stream.Start()

		for {
			select {
			case <-s.done:
				logger.Debugf(ctx, "watching directory is stopped")
				stream.Stop()
				return

			case events := <-stream.Events:
				for _, event := range events {
					path := "/" + event.Path
					ctx := karma.Describe("event.flags", explainEventFlags(event))

					if isEventTemporary(event) {
						logger.Debugf(ctx, "skipping temporary file: %s", path)
						continue
					}

					isEventExcluded := excludes.match("/" + event.Path)
					if isEventExcluded {
						logger.Debugf(ctx, "skipping excluded file: %s", path)
						continue
					}

					logger.Debugf(ctx, "detected changed file: %s", path)
					syncChan <- path[len(source):]
				}
			}
		}
	}()
}

func (s *watcher) stop() {
	logger.Info("stopping watch streams")
	close(s.done)
	s.wg.Wait()
}

func isEventTemporary(event fsevents.Event) bool {
	bits := fsevents.ItemCreated + fsevents.ItemRemoved
	return event.Flags&bits == bits
}

func explainEventFlags(event fsevents.Event) string {
	flags := ""
	for bit, description := range flagsDescription {
		if event.Flags&bit == bit {
			flags += description + " "
		}
	}
	return flags
}
