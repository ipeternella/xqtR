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

const (
	processWarningHeader = "\n>--- Warnings: ---<\n\n"
	processWarningFooter = "\n>-----------------<"
	processStdoutHeader  = "\n>--- Stdout: ---<\n\n"
	processStdoutFooter  = "\n>---------------<"
	processStderrHeader  = "\n>--- Stderr: ---<\n\n"
	processStderrFooter  = "\n>---------------<"
)

func executeJobs(yml *viper.Viper, debug bool) {
	// gets all jobs
	allJobs := yml.Get("jobs").(map[string]interface{})
	jobsCount := len(allJobs)
	log.Info().Msgf("Found %d jobs to execute. Starting...", jobsCount)

	// go to its steps
	for jobName := range allJobs {
		jobSteps := Steps{Steps: []Step{}}

		if err := viper.UnmarshalKey(fmt.Sprintf("jobs.%s", jobName), &jobSteps); err != nil {
			panic(err)
		}

		log.Info().Msgf("📝 job: %s ...", jobName)

		// exec each cmd
		for _, step := range jobSteps.Steps {
			log.Info().Msgf("⏳ step: %s", step.Name)

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
				log.Error().Msgf("%s%s%s", processStderrHeader, stderr, processStderrFooter)
				log.Error().Msgf("⌛ step: %s ✖️", step.Name)
				os.Exit(1)
			}

			// stderr is also used for warnings when the process does not exit with a non-zero status code
			if len(stderr) > 0 {
				log.Warn().Msgf("%s%s%s", processWarningHeader, stderr, processWarningFooter)
			}

			// stdout is print only if debug is on
			if debug && len(stdout) > 0 {
				log.Debug().Msgf("%s%s%s", processStdoutHeader, stdout, processStdoutFooter)
			}

			log.Info().Msgf("⌛ step: %s ✓", step.Name)
		}
	}
}
