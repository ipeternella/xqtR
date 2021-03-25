package executor

import (
	"github.com/IgooorGP/xqtR/internal/dtos"
	"github.com/IgooorGP/xqtR/internal/ui"
)

func ExecuteJobAsync(job dtos.Job, debug bool) {
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

	// publish tasks (the job steps commands) to workers
	for _, jobStep := range job.Steps {
		taskId++
		workerData := newWorkerData(taskId, jobStep, debug)
		taskQueue <- workerData
	}

	// no more tasks (breaks loops from workers)
	close(taskQueue)

	// collect results
	for i := 0; i < numTasks; i++ {
		rslt := <-workerResultsChan

		if rslt.Result.Err != nil {
			ui.PrintCmdFailure(rslt.Name, rslt.Result.StdoutData, rslt.Result.StderrData, continueOnError)
		} else {
			ui.PrintCmdFeedback(rslt.Name, rslt.Result.StdoutData, rslt.Result.StderrData, debug)
		}
	}
}
