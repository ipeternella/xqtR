package executor

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/IgooorGP/xqtR/internal/dtos"
	"github.com/IgooorGP/xqtR/internal/ui"
	"github.com/briandowns/spinner"
	"github.com/rs/zerolog/log"
)

func ExecuteJobAsync(job dtos.Job, debug bool) *dtos.JobResult {
	log.Info().Msgf("📝 job: %s", job.Title)

	jobResult := dtos.NewEmptyJobResult(job)
	jobExecutionRules := dtos.NewJobExecutionRules(debug, job.ContinueOnError)

	continueOnError := job.ContinueOnError
	numTasks := len(job.Steps)

	workerJobStepResultsChan := make(chan *dtos.JobStepResult, numTasks) // buffered channel
	taskQueue := make(chan *dtos.WorkerData)                             // unbuffered channel
	stepId := 0

	// spawn NumWorkers goroutines that are initially blocked (no tasks).
	// workers receive a read-only <-chan taskQueue to consume the steps of a given job and
	// a write-only chan<- workerResultsChan to publish the step command results
	for workerId := 1; workerId <= job.NumWorkers; workerId++ {
		go executeJobStepByWorker(workerJobStepResultsChan, taskQueue)
	}

	// spinner for progress
	s := spinner.New(spinner.CharSets[26], 500*time.Millisecond, spinner.WithWriter(os.Stderr))
	stepsInProgressText := "⏳ Step numbers in progress: [%s]"
	stepsInProgress := []string{}
	s.Start()

	// publish tasks (the job steps commands) to workers
	for _, jobStep := range job.Steps {
		stepId++

		// spinner: adds new step as being 'in progress'
		stepsInProgress = append(stepsInProgress, fmt.Sprintf("step %d", stepId))
		updatedStepsInProgress := fmt.Sprintf(stepsInProgressText, strings.Join(stepsInProgress, ", "))

		updateSpinnerPrefix(s, updatedStepsInProgress, true)

		// sends step to worker
		taskQueue <- newWorkerData(stepId, jobStep, jobExecutionRules)
	}

	// no more tasks (breaks loops from workers)
	close(taskQueue)

	// collect results
	for i := 0; i < numTasks; i++ {
		stepResult := <-workerJobStepResultsChan

		if stepResult.CmdResult.Err != nil {
			updateSpinnerWithCompleteStep(stepsInProgress, stepsInProgressText, stepResult.Id, "💀", s)
			ui.PrintCmdFailure(
				stepResult.JobStep.Name,
				stepResult.CmdResult.StdoutData,
				stepResult.CmdResult.StderrData,
				continueOnError,
			)
		} else {
			updateSpinnerWithCompleteStep(stepsInProgress, stepsInProgressText, stepResult.Id, "👍", s)
			ui.PrintCmdFeedback(
				stepResult.JobStep.Name,
				stepResult.CmdResult.StdoutData,
				stepResult.CmdResult.StderrData,
				debug,
			)
		}

		jobResult.StepsResults[stepResult.Id-1] = stepResult
	}

	log.Info().Msgf("📝 job steps results: [%s]", strings.Join(stepsInProgress, ", "))

	return jobResult
}

// updateSpinner updates the spinner steps in progress with either a success or failure mark
func updateSpinnerWithCompleteStep(stepsInProgress []string, stepsInProgressText string, workerId int, mark string, spin *spinner.Spinner) {
	ix := workerId - 1
	stepsInProgress[ix] = mark
	updatedStepsInProgress := fmt.Sprintf(stepsInProgressText, strings.Join(stepsInProgress, ", "))

	// reload spinner
	updateSpinnerPrefix(spin, updatedStepsInProgress, true)
}

// updateSpinnerPrefix changes the spinner prefix with a new one. Optionally, it clears the spinner's stdout by
// stopping it and starting over.
func updateSpinnerPrefix(spin *spinner.Spinner, newPrefix string, reload bool) {

	// optionally reload it (stop -- clears stdout -- and restart it over)
	if reload {
		spin.Stop()
		spin.Prefix = newPrefix
		spin.Start()
	} else {
		spin.Prefix = newPrefix
	}
}
