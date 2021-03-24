// Package executor is the package responsible for executing the jobs' steps commands
// in a sync fashion way or in an async way by spawning goroutines in a workpool.
package executor

import (
	"github.com/IgooorGP/xqtR/internal/dtos"
	"github.com/rs/zerolog/log"
)

type JobExecutor func(job dtos.Job, debug bool)

type JobDispatcher struct {
	ExecuteJobSync  func(job dtos.Job, debug bool)
	ExecuteJobAsync func(job dtos.Job, debug bool)
}

func (dispatcher JobDispatcher) DispatchForExecution(job dtos.Job, debug bool) {
	if job.NumWorkers > 0 {
		dispatcher.ExecuteJobAsync(job, debug)
	} else {
		dispatcher.ExecuteJobSync(job, debug)
	}
}

func (dispatcher JobDispatcher) DispatchJobsForExecution(jobs []dtos.Job, debug bool) {
	for _, job := range jobs {
		log.Info().Msgf("📝 job: %s", job.Title)

		dispatcher.DispatchForExecution(job, debug)
	}
}

func NewJobDispatcher(syncJobExecutor JobExecutor, asyncJobExecutor JobExecutor) *JobDispatcher {
	return &JobDispatcher{
		ExecuteJobSync:  syncJobExecutor,
		ExecuteJobAsync: asyncJobExecutor,
	}
}

func ExecuteJobs(yaml dtos.JobsFile, debug bool) {
	dispatcher := NewJobDispatcher(ExecuteJobSync, ExecuteJobAsync)

	dispatcher.DispatchJobsForExecution(yaml.Jobs, debug)
}
