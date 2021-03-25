package xqtr

import (
	"errors"
	"fmt"
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

	// yaml schema validation
	if err := validateYamlSchema(yaml); err != nil {
		log.Fatal().Msgf("yaml structure error: %s", err.Error()) // status code 1
	}

	return yaml
}

func validateYamlSchema(yaml *dtos.JobsFile) error {
	jobNum := 0

	if len(yaml.Jobs) == 0 {
		return errors.New(`missing "jobs" key on yaml file`)
	}

	for _, job := range yaml.Jobs {
		jobNum++     // increment job number to help the user find the job in case of errors
		stepNum := 0 // for steps errors

		if job.Title == "" {
			return fmt.Errorf(`missing "jobs" key on yaml file (job number: %d)`, jobNum)
		}

		if len(job.Steps) == 0 {
			return fmt.Errorf(`missing "steps" key for (job name: "%s", job number: %d)`, job.Title, jobNum)
		}

		for _, step := range job.Steps {
			stepNum++ // increment step number to help the user find the step in case of errors

			if step.Name == "" {
				return fmt.Errorf(`missing "name" key for step on (job name: "%s", job number: %d, step number: %d)`, job.Title, jobNum, stepNum)
			}

			if step.Run == "" {
				return fmt.Errorf(`missing "run" key for step on (job name: "%s", job number: %d, step name: "%s", step number: %d)`, job.Title, jobNum, step.Name, stepNum)
			}
		}
	}

	return nil
}
