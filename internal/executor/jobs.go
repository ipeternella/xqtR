// Package executor is the package responsible for executing the jobs' steps commands
// in a sync fashion way or in an async way by spawning goroutines in a workpool.
package executor

import (
	"github.com/IgooorGP/xqtR/internal/dtos"
	"github.com/rs/zerolog/log"
)

// ExecuteJobs dispatches the jobs found on the yaml file to be run sync or async
func ExecuteJobs(yaml dtos.JobsFile, debug bool) {
	for _, job := range yaml.Jobs {
		dispatchForExecution(job, debug)
	}
}

// dispatchForExecution decides whether the job's steps should run in workerpool (async) or not (sync)
func dispatchForExecution(job dtos.Job, debug bool) {
	log.Info().Msgf("📝 job: %s", job.Title)

	if job.NumWorkers > 0 {
		ExecuteJobAsync(job, debug)
	} else {
		ExecuteJobSync(job, debug)
	}
}
