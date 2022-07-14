package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/docopt/docopt-go"
)

var version = "develop"

const usage = `resyncd – synchronize changed files to a remote servers.

usage:
  resyncd -h | --help
  resyncd [options] <config>

options:
  -v --verbose  Verbose logging.
`

func main() {
	args, err := docopt.ParseArgs(usage, nil, version)
	if err != nil {
		panic(err)
	}

	logger := setupLogger(args["--verbose"].(bool))

	var conf config
	if _, err := toml.DecodeFile(args["<config>"].(string), &conf); err != nil {
		logger.Error("unable to read configuration file", slog.Any("err", err))
		os.Exit(1)
	}

	if len(conf.Syncs) == 0 {
		logger.Error("no configuration found")
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	var wg sync.WaitGroup
	for _, syncConf := range conf.Syncs {
		filesCh := make(chan string)

		wg.Add(1)
		go func() {
			defer wg.Done()

			err := watchDirectory(ctx, logger, syncConf, filesCh)
			if err != nil && err != context.Canceled {
				logger.Error("error watching directory", slog.Any("err", err))
				cancel()
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()

			bufferedCh := bufferize(filesCh, bufferizeDuration)
			rsyncDirectory(logger, syncConf, bufferedCh)
		}()
	}

	<-ctx.Done()
	logger.Info("stopping synchronization")

	wg.Wait()
}
