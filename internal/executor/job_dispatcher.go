// Package executor is the package responsible for executing the jobs' steps commands
// in a sync fashion way or in an async way by spawning goroutines in a workpool.
package executor

import (
	"github.com/IgooorGP/xqtR/internal/dtos"
)

type JobExecutor func(job dtos.Job, debug bool)

type JobDispatcher struct {
	ExecuteJobSync  JobExecutor
	ExecuteJobAsync JobExecutor
}

// DispatchForExecution uses the defined `num_workers` from the yaml file to run a job
// synchronously or asynchronously with more goroutines.
func (dispatcher JobDispatcher) DispatchForExecution(job dtos.Job, debug bool) {
	if job.NumWorkers > 0 {
		dispatcher.ExecuteJobAsync(job, debug)
	} else {
		dispatcher.ExecuteJobSync(job, debug)
	}
}

// DispatchJobsForExecution is a wrapper which calls DispatchForExecution for each distinct
// job given by the job yaml file.
func (dispatcher JobDispatcher) DispatchJobsForExecution(jobs []dtos.Job, debug bool) {
	for _, job := range jobs {
		dispatcher.DispatchForExecution(job, debug)
	}
}

// NewJobDispatcher creates a new job dispatcher with executors to run jobs synchronously
// or asynchronously with more goroutines.
func NewJobDispatcher(syncJobExecutor JobExecutor, asyncJobExecutor JobExecutor) *JobDispatcher {
	return &JobDispatcher{
		ExecuteJobSync:  syncJobExecutor,
		ExecuteJobAsync: asyncJobExecutor,
	}
}
