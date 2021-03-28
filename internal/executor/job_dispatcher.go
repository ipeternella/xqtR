// Package executor is the package responsible for executing the jobs' steps commands
// in a sync fashion way or in an async way by spawning goroutines in a workpool.
package executor

import (
	"github.com/IgooorGP/xqtR/internal/dtos"
)

type JobExecutor func(job dtos.Job, debug bool) dtos.JobResult

type JobDispatcher struct {
	ExecuteJobSync  JobExecutor
	ExecuteJobAsync JobExecutor
}

func jobResultContainsErrors(jobResult dtos.JobResult) bool {
	for _, stepResult := range jobResult.StepsResults {
		if stepResult.HasError {
			return true
		}
	}

	return false
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
func (dispatcher JobDispatcher) DispatchJobsForExecution(jobs []dtos.Job, debug bool) dtos.JobsYamlResult {
	yamlResults := dtos.NewEmptyJobsYamlResult(jobs)

	for i, job := range jobs {
		jobResult := dispatcher.DispatchForExecution(job, debug)
		yamlResults.JobResults[i] = jobResult // rewrite with results

		// job has steps with errors and `continue_on_error` is set to false: break
		if jobResultContainsErrors(jobResult) && !job.ContinueOnError {
			break
		}
	}

	return yamlResults
}

// NewJobDispatcher creates a new job dispatcher with executors to run jobs synchronously
// or asynchronously with more goroutines.
func NewJobDispatcher(syncJobExecutor JobExecutor, asyncJobExecutor JobExecutor) *JobDispatcher {
	return &JobDispatcher{
		ExecuteJobSync:  syncJobExecutor,
		ExecuteJobAsync: asyncJobExecutor,
	}
}
