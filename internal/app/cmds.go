package app

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
)

type Step struct {
	Name string `mapstructure:"name"`
	Run  string `mapstructure:"run"`
}

type Steps struct {
	Steps []Step `mapstructure:"steps"`
}

// ExtractCommandAndArgsFromString - reads a string and parses it to a terminal command
func shellCommand(command string) (*exec.Cmd, *io.ReadCloser, *io.ReadCloser) {
	cmd := exec.Command("bash", "-c", command) // TODO: improve

	// adds an os.Pipe to the cmd stdstreams so that xqtR's main process can read the
	// the stdout and stderr of the new cmd spawned process
	stdoutPipe, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()

	return cmd, &stdoutPipe, &stderrPipe
}

func readPipe(cmdPipe *io.ReadCloser) []byte {
	bufferData, err := ioutil.ReadAll(*cmdPipe)

	if err != nil {
		log.Error().Msgf("Error when reading cmd buffer: %s", err.Error())
		os.Exit(1)
	}

	return bufferData
}

func readCmdStdStreams(cmdStdoutPipe *io.ReadCloser, cmdStderrPipe *io.ReadCloser, readStdout bool) ([]byte, []byte) {
	var stdoutData []byte
	var stderrData []byte

	if readStdout {
		stdoutData = readPipe(cmdStdoutPipe)
	}
	stderrData = readPipe(cmdStderrPipe)

	return stdoutData, stderrData
}

func executeJobStep(jobStep JobStep, debug bool) {
	log.Info().Msgf("⏳ step: %s", jobStep.Name)

	cmd, cmdStdoutPipe, cmdStderrPipe := shellCommand(jobStep.Run)

	// spawns a new OS process with the cmd
	if err := cmd.Start(); err != nil {
		log.Fatal().Msgf("An error happened while starting the cmd", err.Error())
	}

	// pipes must be passed by reference, naturally
	stdoutData, stderrData := readCmdStdStreams(cmdStdoutPipe, cmdStderrPipe, debug)

	// waits for cmd completion (also closes stdstreams)
	if err := cmd.Wait(); err != nil {
		printCmdFailure(jobStep.Name, stdoutData, stderrData, debug)
	}

	printCmdFeedback(jobStep.Name, stdoutData, stderrData, debug)
}

func executeJob(job Job, debug bool) {
	log.Info().Msgf("📝 job: %s", job.Title)

	for _, jobStep := range job.Steps {
		executeJobStep(jobStep, debug)
	}
}

func executeJobs(yaml *JobsFile, debug bool) {
	for _, job := range yaml.Jobs {
		executeJob(job, debug)
	}
}
