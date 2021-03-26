package executor

import (
	"fmt"
	"os"
	"time"

	"github.com/IgooorGP/xqtR/internal/dtos"
	"github.com/IgooorGP/xqtR/internal/ui"
	"github.com/briandowns/spinner"
	"github.com/rs/zerolog/log"
)

func ExecuteJobSync(job dtos.Job, debug bool) {
	log.Info().Msgf("📝 job: %s", job.Title)

	for _, jobStep := range job.Steps {
		executeJobStep(jobStep, debug, job.ContinueOnError)
	}
}

func executeJobStep(jobStep dtos.JobStep, debug bool, continueOnError bool) {
	cmd, cmdStdoutPipe, cmdStderrPipe := shellCommand(jobStep.Run)

	s := spinner.New(spinner.CharSets[26], 500*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.Prefix = fmt.Sprintf("⏳ step: %s ", jobStep.Name)
	s.Start()

	// spawns a new OS process with the cmd
	if err := cmd.Start(); err != nil {
		log.Fatal().Msgf("An error happened while starting the cmd: %s", err.Error())
	}

	stdoutData, stderrData := readCmdStdStreams(cmdStdoutPipe, cmdStderrPipe, debug)

	s.Stop()

	// waits for cmd completion (also closes stdstreams)
	if err := cmd.Wait(); err != nil {
		ui.PrintCmdFailure(jobStep.Name, stdoutData, stderrData, continueOnError)
	} else {
		ui.PrintCmdFeedback(jobStep.Name, stdoutData, stderrData, debug)
	}
}
