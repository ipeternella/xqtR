// Package startup is used to configure basic needs of the application such as
// logging levels, formatting, etc. XqtR app instances rely on the Boot() method.
package startup

import (
	"os"
	"time"

	"github.com/IgooorGP/xqtR/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Boot is a func that boots the tool by running the appropriate configurations.
func Boot(cfg *config.XqtRConfig) {
	setupLogger(cfg)

	log.Debug().Msg("Booting xqtR is complete.")
}

func setupLogger(cfg *config.XqtRConfig) {
	// zerolog's global logger setup
	logFormat := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Kitchen}
	log.Logger = zerolog.New(logFormat).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(config.ParseLevel(cfg.LogLevel))
}
