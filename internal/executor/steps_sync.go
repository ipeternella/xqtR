package executor

import (
	"github.com/IgooorGP/xqtR/internal/dtos"
	"github.com/IgooorGP/xqtR/internal/ui"
	"github.com/rs/zerolog/log"
)

func ExecuteJobSync(job dtos.Job, debug bool) {
	for _, jobStep := range job.Steps {
		executeJobStep(jobStep, debug)
	}
}

func executeJobStep(jobStep dtos.JobStep, debug bool) {
	log.Info().Msgf("⏳ step: %s", jobStep.Name)

	cmd, cmdStdoutPipe, cmdStderrPipe := shellCommand(jobStep.Run)

	// spawns a new OS process with the cmd
	if err := cmd.Start(); err != nil {
		log.Fatal().Msgf("An error happened while starting the cmd: %s", err.Error())
	}

	stdoutData, stderrData := readCmdStdStreams(cmdStdoutPipe, cmdStderrPipe, debug)

	// waits for cmd completion (also closes stdstreams)
	if err := cmd.Wait(); err != nil {
		ui.PrintCmdFailure(jobStep.Name, stdoutData, stderrData, debug)
	}

	ui.PrintCmdFeedback(jobStep.Name, stdoutData, stderrData, debug)
}
