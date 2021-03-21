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
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	zerolog.SetGlobalLevel(config.ParseLevel(cfg.LogLevel))
}
