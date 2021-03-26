// Package xqtr provides functions to create new xqtR app instances (core of this project).
package xqtr

import (
	"os"
	"path/filepath"
	"time"

	"github.com/IgooorGP/xqtR/internal/config"
	"github.com/IgooorGP/xqtR/internal/executor"
	"github.com/IgooorGP/xqtR/internal/startup"
	"github.com/rs/zerolog/log"
)

// Xqtr represents an xqtr app instance
type XqtR struct {
	config *config.XqtRConfig
}

// NewXqtR creates new XqtR app instances with a given config
func NewXqtR(cfg *config.XqtRConfig) *XqtR {
	return &XqtR{config: cfg}
}

// Startup is the function responsible for booting apps
func (xqtr *XqtR) Startup() {
	startup.Boot(xqtr.config)
}

// Run is the main function for running jobs
func (xqtr *XqtR) Run() {
	xqtr.Startup()

	start := time.Now()
	log.Info().Msgf("🛠️  xqtR is starting...")

	// gets absolute file path location
	absJobsFilePath, err := filepath.Abs(xqtr.config.JobFilePath)

	// bad file path
	if err != nil {
		log.Error().Msgf("Malformed job file path: %s", xqtr.config.JobFilePath)
		os.Exit(1)
	}

	// parse yaml
	viperParser := prepareViper(absJobsFilePath)
	yaml := parseYaml(viperParser, absJobsFilePath)

	// yaml schema validation
	if err := validateYamlSchema(yaml); err != nil {
		log.Fatal().Msgf("yaml structure error: %s", err.Error()) // status code 1
	}

	// extract jobs and create functions to invoke them
	executor.ExecuteJobs(*yaml, true) // no need for pointers from here on (no modifications on the yml)

	log.Info().Msgf("🛠️  xqtR is all done! Total duration: %s", time.Since(start))
}
