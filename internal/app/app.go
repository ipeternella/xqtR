package app

import (
	"os"
	"path/filepath"

	"github.com/IgooorGP/xqtR/internal/config"
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

	// extract jobs and create functions to invoke them
	executeJobs(yaml, true)
}
