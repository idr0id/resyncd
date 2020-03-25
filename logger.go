package main

import (
	"github.com/kovetskiy/lorg"
	"github.com/reconquest/cog"
)

var logger *cog.Logger

func setupLogger() {
	format := lorg.NewFormat("${time} ${level:[%s]:right:short} ${prefix}%s")

	stderr := lorg.NewLog()
	stderr.SetIndentLines(true)
	stderr.SetFormat(format)

	logger = cog.NewLogger(stderr)
	logger.SetLevel(lorg.LevelInfo)
}

func verboseLogging() {
	logger.SetLevel(lorg.LevelDebug)
}