// Package startup is used to configure basic needs of the application such as
// logging levels, formatting, etc. XqtR app instances rely on the Boot() method.
package startup

import (
	"io"
	"os"
	"time"

	"github.com/IgooorGP/xqtR/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Boot is a func that boots the tool by running the appropriate configurations.
func Boot(cfg *config.XqtRConfig) bool {
	stdout := os.Stdout        // grab stdout file descriptor from the OS
	timeFormat := time.Kitchen // timeformat

	setupLogger(cfg, stdout, timeFormat)
	log.Debug().Msg("Booting xqtR is complete.")

	return true
}

func setupLogger(cfg *config.XqtRConfig, outputStream io.Writer, timeFormat string) {
	// zerolog's global logger setup
	logFormat := zerolog.ConsoleWriter{Out: outputStream, TimeFormat: timeFormat}
	log.Logger = zerolog.New(logFormat).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(config.ParseLevel(cfg.LogLevel))
}
