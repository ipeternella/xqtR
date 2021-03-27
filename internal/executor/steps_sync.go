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

func ExecuteJobSync(job dtos.Job, debug bool) *dtos.JobResult {
	jobResult := dtos.NewEmptyJobResult(job)
	jobExecutionRules := dtos.NewJobExecutionRules(debug, job.ContinueOnError)
	stepId := 0

	log.Info().Msgf("📝 job: %s", job.Title)

	for _, jobStep := range job.Steps {
		stepId++

		jobStepResult := executeJobStep(stepId, jobStep, jobExecutionRules) // mutates stepResult
		jobResult.StepsResults[jobStepResult.Id-1] = jobStepResult          // stepId is 1-index based and not 0-index based

		// breaks if there's an error and we should not continue upon errors
		if jobStepResult.HasError && !jobExecutionRules.ContinueOnError {
			return jobResult
		}
	}

	return jobResult
}

func executeJobStep(jobStepId int, jobStep dtos.JobStep, executionRules dtos.JobExecutionRules) *dtos.JobStepResult {
	var cmdResult dtos.CmdResult
	var debug = executionRules.Debug
	var continueOnError = executionRules.ContinueOnError
	var stepResult = dtos.NewEmptyJobStepResult(jobStepId, jobStep)

	cmd, cmdStdoutPipe, cmdStderrPipe := shellCommand(jobStep.Run)

	s := spinner.New(spinner.CharSets[26], 500*time.Millisecond, spinner.WithWriter(os.Stderr))
	s.Prefix = fmt.Sprintf("⏳ step: %s ", jobStep.Name)
	s.Start()

	// spawns a new OS process with the cmd
	if err := cmd.Start(); err != nil {
		log.Fatal().Msgf("An error happened while starting the cmd: %s", err.Error())
	}

	// reads cmd stdstreams until they are over (blocks the calling goroutine)
	stdoutData, stderrData := readCmdStdStreams(cmdStdoutPipe, cmdStderrPipe, debug)

	s.Stop()

	// waits for cmd completion to close stdstreams
	if err := cmd.Wait(); err != nil {
		cmdResult = dtos.NewCmdResult(stdoutData, stderrData, err)
		markStepAsExecuted(stepResult, cmdResult)

		ui.PrintCmdFailure(jobStep.Name, stdoutData, stderrData, continueOnError)
	} else {
		cmdResult = dtos.NewCmdResult(stdoutData, stderrData, nil)
		markStepAsExecuted(stepResult, cmdResult)

		ui.PrintCmdFeedback(jobStep.Name, stdoutData, stderrData, debug)
	}

	return stepResult
}

func markStepAsExecuted(stepResult *dtos.JobStepResult, cmdResult dtos.CmdResult) {
	stepResult.Executed = true
	stepResult.CmdResult = &cmdResult

	if cmdResult.Err != nil {
		stepResult.HasError = true
	}
}
