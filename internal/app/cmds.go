package app

import (
	"fmt"
	"io/ioutil"
	"os"

	"os/exec"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Step struct {
	Name string `mapstructure:"name"`
	Run  string `mapstructure:"run"`
}

type Steps struct {
	Steps []Step `mapstructure:"steps"`
}

// ExtractCommandAndArgsFromString - reads a string and parses it to a terminal command
func shellCommand(command string) *exec.Cmd {
	cmd := exec.Command("bash", "-c", command) // TODO: improve
	return cmd
}

func executeJob(yml *viper.Viper, debug bool) {
	// gets all jobs
	allJobs := yml.Get("jobs").(map[string]interface{})
	jobsCount := len(allJobs)
	log.Info().Msgf("Found %d jobs to execute. Starting...", jobsCount)

	// go to its steps
	for k := range allJobs {
		jobSteps := Steps{Steps: []Step{}}

		if err := viper.UnmarshalKey(fmt.Sprintf("jobs.%s", k), &jobSteps); err != nil {
			panic(err)
		}

		// exec each cmd
		for _, step := range jobSteps.Steps {
			log.Info().Msgf("⏳: %s", step.Name)

			cmd := shellCommand(step.Run)
			stdoutBuffer, _ := cmd.StdoutPipe()
			stderrBuffer, _ := cmd.StderrPipe()
			var stdout []byte
			var stderr []byte

			if err := cmd.Start(); err != nil {
				log.Fatal().AnErr("An error happened while starting the cmd", err)
			}

			if debug {
				// log.Debug().Msgf("Reading from cmd's stdout...: %s", cmd)
				stdout, _ = ioutil.ReadAll(stdoutBuffer)
			}

			// log.Debug().Msgf("Reading from cmd's stderr...: %s", cmd)
			stderr, _ = ioutil.ReadAll(stderrBuffer)

			// wait for cmd completion
			if err := cmd.Wait(); err != nil {
				if err.Error() == "exec: not started" {
					log.Error().Msgf("❌: cmd %s was not found!", cmd.Path)
				} else {
					log.Error().Msgf("❌: cmd %s has raised an error!", cmd.Path)
					log.Error().Msgf("📜: %s", stderr)
				}

				os.Exit(1)
			}

			if debug && len(stdout) > 0 {
				log.Debug().Msgf("📜 stdout:\n%s", stdout)
			}

			log.Info().Msgf("⌛: %s ✓\n", step.Name)
		}
	}
}
