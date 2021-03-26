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

func ExecuteJobAsync(job dtos.Job, debug bool) {
	log.Info().Msgf("📝 job: %s", job.Title)

	continueOnError := job.ContinueOnError
	numTasks := len(job.Steps)

	workerResultsChan := make(chan *dtos.WorkerResult, numTasks) // buffered channel
	taskQueue := make(chan *dtos.WorkerData)                     // unbuffered channel
	taskId := 0

	// spawn NumWorkers goroutines that are initially blocked (no tasks).
	// workers receive a read-only <-chan taskQueue to consume the steps of a given job and
	// a write-only chan<- workerResultsChan to publish the step command results
	for workerId := 1; workerId <= job.NumWorkers; workerId++ {
		go executeJobStepByWorker(workerResultsChan, taskQueue)
	}

	// spinner for progress
	s := spinner.New(spinner.CharSets[26], 500*time.Millisecond, spinner.WithWriter(os.Stderr))
	stepsInProgressText := "⏳ Step numbers in progress: [%s]"
	stepsInProgress := []string{}
	s.Start()

	// publish tasks (the job steps commands) to workers
	for _, jobStep := range job.Steps {
		taskId++

		// spinner: adds new step as being 'in progress'
		stepsInProgress = append(stepsInProgress, fmt.Sprintf("step %d", taskId))
		updatedStepsInProgress := fmt.Sprintf(stepsInProgressText, strings.Join(stepsInProgress, ", "))

		updateSpinnerPrefix(s, updatedStepsInProgress, true)

		// sends step to worker
		workerData := newWorkerData(taskId, jobStep, debug)
		taskQueue <- workerData
	}

	// no more tasks (breaks loops from workers)
	close(taskQueue)

	// collect results
	for i := 0; i < numTasks; i++ {
		rslt := <-workerResultsChan

		if rslt.Result.Err != nil {
			updateSpinnerWithCompleteStep(stepsInProgress, stepsInProgressText, rslt.WorkerId, "💀", s) // mark step as complete
			ui.PrintCmdFailure(rslt.Name, rslt.Result.StdoutData, rslt.Result.StderrData, continueOnError)
		} else {
			updateSpinnerWithCompleteStep(stepsInProgress, stepsInProgressText, rslt.WorkerId, "👍", s) // mark step as complete
			ui.PrintCmdFeedback(rslt.Name, rslt.Result.StdoutData, rslt.Result.StderrData, debug)
		}
	}

	log.Info().Msgf("📝 job steps results: [%s]", strings.Join(stepsInProgress, ", "))
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
