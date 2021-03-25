package executor

import "github.com/IgooorGP/xqtR/internal/dtos"

// ExecuteJobs runs all the jobs' steps by using a dispatcher which controls the logic
// to when to execute a given job in a sync way or in an async fashion by spawning
// many goroutines. See `job_dispatcher.go` fore more.
func ExecuteJobs(yaml dtos.JobsFile, debug bool) {
	dispatcher := NewJobDispatcher(ExecuteJobSync, ExecuteJobAsync)

	dispatcher.DispatchJobsForExecution(yaml.Jobs, debug)
}
