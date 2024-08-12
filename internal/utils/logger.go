package utils

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func SetupLogger() {
	writers := []io.Writer{zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}}
	if _, ok := os.LookupEnv("PRODUCTION"); ok {
		logFile, _ := os.OpenFile(
			"logs.txt",
			os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666,
		)
		writers = append(writers, logFile)
	}

	log.Logger = zerolog.New(zerolog.MultiLevelWriter(writers...)).With().Timestamp().Logger()
}
