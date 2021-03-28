package startup

import (
	"bytes"
	"testing"
	"time"

	"github.com/IgooorGP/xqtR/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestShouldBootWithoutErrors(t *testing.T) {
	// arrange
	cfg := config.NewXqtRConfigWithDefaults()

	// act
	booted := Boot(&cfg)

	// assert
	assert.True(t, booted) // successfully booted
}

func TestShouldPrintToStdoutWhenLogLevelIsInfo(t *testing.T) {
	// arrange
	cfg := config.NewXqtRConfigWithDefaults()
	var logBuffer bytes.Buffer

	// act - sets log level to info
	setupLogger(&cfg, &logBuffer, time.Kitchen)

	// arrange - redirect logs to local buffer
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: &logBuffer, TimeFormat: time.Kitchen})

	// act - use the logger (info is ok!)
	log.Info().Msg("hi!")

	// assert
	var usedStdout = logBuffer.String()

	assert.Contains(t, usedStdout, "hi!")
}

func TestShouldNotPrintToStdoutWhenLogLevelIsDebug(t *testing.T) {
	// arrange
	cfg := config.NewXqtRConfigWithDefaults()
	var logBuffer bytes.Buffer

	// act - sets log level to info
	setupLogger(&cfg, &logBuffer, time.Kitchen)

	// arrange - redirect logs to local buffer
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: &logBuffer, TimeFormat: time.Kitchen})

	// act - use the logger (debug level is lower than info --> nothing printed to stdout)
	log.Debug().Msg("hi!")

	// assert
	var usedStdout = logBuffer.String()

	assert.Equal(t, "", usedStdout)
}
