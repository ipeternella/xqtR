// Package executor is the package responsible for executing the jobs' steps commands
// in a sync fashion way or in an async way by spawning goroutines in a workpool.
package executor

import (
	"github.com/IgooorGP/xqtR/internal/dtos"
	"github.com/rs/zerolog/log"
)

type JobExecutor func(job dtos.Job, debug bool) dtos.JobResult

type JobDispatcher struct {
	ExecuteJobSync  JobExecutor
	ExecuteJobAsync JobExecutor
}

// DispatchForExecution uses the defined `num_workers` from the yaml file to run a job
// synchronously or asynchronously with more goroutines.
func (dispatcher JobDispatcher) DispatchForExecution(job dtos.Job, debug bool) dtos.JobResult {
	var jobResult dtos.JobResult

	if job.NumWorkers > 1 {
		jobResult = dispatcher.ExecuteJobAsync(job, debug)
	} else {
		jobResult = dispatcher.ExecuteJobSync(job, debug)
	}

	return jobResult
}

// DispatchJobsForExecution is a wrapper which calls DispatchForExecution for each distinct
// job given by the job yaml file.
func (dispatcher JobDispatcher) DispatchJobsForExecution(jobs []dtos.Job, debug bool) *dtos.JobsYamlResult {
	jobResults := []dtos.JobResult{}

	for _, job := range jobs {
		jobResult := dispatcher.DispatchForExecution(job, debug)
		jobResults = append(jobResults, jobResult)
	}

	for _, jobResult := range jobResults {

		for _, stepResult := range jobResult.StepsResults {
			log.Info().Msgf("%s, executed: %t, error: %t", stepResult.JobStep.Name, stepResult.Executed, stepResult.HasError)
		}
	}

	return dtos.NewJobsYamlResult(jobResults)
}

// NewJobDispatcher creates a new job dispatcher with executors to run jobs synchronously
// or asynchronously with more goroutines.
func NewJobDispatcher(syncJobExecutor JobExecutor, asyncJobExecutor JobExecutor) *JobDispatcher {
	return &JobDispatcher{
		ExecuteJobSync:  syncJobExecutor,
		ExecuteJobAsync: asyncJobExecutor,
	}
}
