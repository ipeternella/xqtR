package app

import (
	"os"
	"path/filepath"

	"github.com/IgooorGP/xqtR/internal/config"
	"github.com/IgooorGP/xqtR/internal/startup"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
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
	absFilePath, _ := filepath.Abs(xqtr.config.JobFilePath)
	fileName := filepath.Base(absFilePath)
	fileFolder := filepath.Dir(absFilePath)

	log.Debug().Msgf("Yaml data: Filename: %s, FileFolder: %s", fileName, fileFolder)

	viper.SetConfigType("yaml")
	viper.SetConfigName(fileName)
	viper.AddConfigPath(fileFolder)

	// error handling when reading config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Error().Msgf("File was not found: %s", absFilePath)
		} else {
			// Config file was found but another error was produced
			log.Error().Msgf("An error occured while reading yaml config: %s", err.Error())
		}

		os.Exit(1)
	}
	yml := viper.GetViper()

	// extract jobs and create functions to invoke them
	executeJobs(yml, true)
}
