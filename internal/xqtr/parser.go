package xqtr

import (
	"os"
	"path/filepath"

	"github.com/IgooorGP/xqtR/internal/dtos"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func prepareViper(absJobsFilePath string) *viper.Viper {
	fileName := filepath.Base(absJobsFilePath)
	fileFolder := filepath.Dir(absJobsFilePath)

	log.Debug().Msgf("Preparing Viper with filename: %s, fileFolder: %s", fileName, fileFolder)

	viper.SetConfigType("yaml")
	viper.SetConfigName(fileName)
	viper.AddConfigPath(fileFolder)

	return viper.GetViper() // configured viper instance
}

func parseYaml(viperParser *viper.Viper, absJobsFilePath string) *dtos.JobsFile {
	// error handling when reading config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Error().Msgf("File was not found: %s", absJobsFilePath)
		} else {
			log.Error().Msgf("An error occured while reading yaml config: %s", err.Error())
		}

		os.Exit(1)
	}

	// yaml unmarshalling (viper requires a pointer)
	yaml := dtos.NewJobsFile()
	if err := viper.Unmarshal(yaml); err != nil {
		panic(err)
	}

	return yaml
}
