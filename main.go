package main

import (
	"github.com/BurntSushi/toml"
	"github.com/docopt/docopt-go"
	"github.com/reconquest/karma-go"
	"os"
	"os/signal"
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
	setupLogger()

	args, err := docopt.ParseArgs(usage, nil, version)
	if err != nil {
		panic(err)
	}

	if args["--verbose"].(bool) {
		verboseLogging()
	}

	var conf config
	if _, err := toml.DecodeFile(args["<config>"].(string), &conf); err != nil {
		logger.Fatalf(err, "unable to read configuration file")
	}

	if len(conf.Syncs) == 0 {
		logger.Fatalf(nil, "no configuration found")
	}

	logger.Infof(nil, "loading %d configurations", len(conf.Syncs))

	synchronizers := make([]*synchronizer, 0)
	for _, cfg := range conf.Syncs {
		synchronizer := newSynchronizer()
		synchronizers = append(synchronizers, synchronizer)
		go synchronizer.start(cfg)
	}

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt, os.Kill)

	for {
		select {
		case sig := <-signalChan:
			logger.Infof(karma.Describe("signal", sig.String()), "stopping synchronizers")
			for _, synchronizer := range synchronizers {
				synchronizer.stop()
			}
			return
		}
	}
}
