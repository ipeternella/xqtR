package main

import (
	"time"

	"github.com/IgooorGP/xqtR/pkg/cmd"
	"github.com/rs/zerolog/log"
)

func main() {
	start := time.Now()
	xqtR := cmd.NewXqtRCmd()

	xqtR.Execute()
	log.Info().Msgf("Duration: %s", time.Since(start))
}
