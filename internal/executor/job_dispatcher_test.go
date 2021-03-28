package executor

import (
	"os"
	"testing"

	"github.com/IgooorGP/xqtR/internal/config"
	"github.com/IgooorGP/xqtR/internal/dtos"
	"github.com/IgooorGP/xqtR/internal/startup"
	"github.com/IgooorGP/xqtR/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	cfg := config.NewXqtRConfigWithDefaults()
	startup.Boot(&cfg) // log level as info

	code := m.Run()
	os.Exit(code)
}

func TestDispatchForSyncJobOnly(t *testing.T) {
	// arrange - mocks
	syncCalled := false
	asyncCalled := false

	syncJobExecutorMock := func(job dtos.Job, debug bool) dtos.JobResult {
		var mockJobSteps []dtos.JobStep
		syncCalled = true

		mockResult := dtos.Job{Title: "mock", Steps: mockJobSteps}

		return dtos.NewEmptyJobResult(mockResult)
	}

	asyncJobExecutorMock := func(job dtos.Job, debug bool) dtos.JobResult {
		var mockJobSteps []dtos.JobStep
		asyncCalled = true

		mockResult := dtos.Job{Title: "mock", Steps: mockJobSteps}

		return dtos.NewEmptyJobResult(mockResult)
	}

	yaml := testutils.NewMockJobsFileWithoutNumWorkers() // sync jobs only
	debug := true
	dispatcher := NewJobDispatcher(syncJobExecutorMock, asyncJobExecutorMock)

	// act
	dispatcher.DispatchJobsForExecution(yaml.Jobs, debug)

	// arrange
	assert.True(t, syncCalled) // only sync jobs
	assert.False(t, asyncCalled)
}

func TestDispatchForAsyncJobOnly(t *testing.T) {
	// arrange - mocks
	syncCalled := false
	asyncCalled := false

	syncJobExecutorMock := func(job dtos.Job, debug bool) dtos.JobResult {
		var mockJobSteps []dtos.JobStep
		syncCalled = true

		mockResult := dtos.Job{Title: "mock", Steps: mockJobSteps}

		return dtos.NewEmptyJobResult(mockResult)
	}

	asyncJobExecutorMock := func(job dtos.Job, debug bool) dtos.JobResult {
		var mockJobSteps []dtos.JobStep
		asyncCalled = true

		mockResult := dtos.Job{Title: "mock", Steps: mockJobSteps}

		return dtos.NewEmptyJobResult(mockResult)
	}

	yaml := testutils.NewMockJobsFileWithNumWorkers() // async job only
	debug := true
	dispatcher := NewJobDispatcher(syncJobExecutorMock, asyncJobExecutorMock)

	// act
	dispatcher.DispatchJobsForExecution(yaml.Jobs, debug)

	// arrange
	assert.False(t, syncCalled)
	assert.True(t, asyncCalled) // only async jobs
}

func TestDispatchForSyncAndAsyncJobs(t *testing.T) {
	// arrange - mocks
	syncCalled := false
	asyncCalled := false

	syncJobExecutorMock := func(job dtos.Job, debug bool) dtos.JobResult {
		var mockJobSteps []dtos.JobStep
		syncCalled = true

		mockResult := dtos.Job{Title: "mock", Steps: mockJobSteps}

		return dtos.NewEmptyJobResult(mockResult)
	}

	asyncJobExecutorMock := func(job dtos.Job, debug bool) dtos.JobResult {
		var mockJobSteps []dtos.JobStep
		asyncCalled = true

		mockResult := dtos.Job{Title: "mock", Steps: mockJobSteps}

		return dtos.NewEmptyJobResult(mockResult)
	}

	yaml := testutils.NewMockJobsFileWithSyncAndAsyncJobs() // async job only
	debug := true
	dispatcher := NewJobDispatcher(syncJobExecutorMock, asyncJobExecutorMock)

	// act
	dispatcher.DispatchJobsForExecution(yaml.Jobs, debug)

	// arrange
	assert.True(t, syncCalled)
	assert.True(t, asyncCalled)
}

func TestDispatchJobsForExecutionShouldExecuteJobs_ContinueOnErrorFalse_NoErrors_Sync(t *testing.T) {
	// arrange
	yaml := testutils.NewJobsFileWithTwoSyncTasks()
	dispatcher := NewJobDispatcher(ExecuteJobSync, ExecuteJobAsync)
	debug := true

	// act
	yamlResults := dispatcher.DispatchJobsForExecution(yaml.Jobs, debug)

	// assert -> all jobs are executed without errors
	for _, jobResult := range yamlResults.JobResults {
		assert.True(t, jobResult.Executed)
		assert.False(t, jobResult.HasError)
	}
}

func TestDispatchJobsForExecutionShouldExecuteJobs_ContinueOnErrorFalse_WithError_Sync(t *testing.T) {
	// arrange -> first job has an error, second has not!
	continueOnError := false
	yaml := testutils.NewJobFileTwoJobs_FirstJobWithError("job with error", 0, continueOnError) // sync jobs
	dispatcher := NewJobDispatcher(ExecuteJobSync, ExecuteJobAsync)
	debug := true

	// act
	yamlResults := dispatcher.DispatchJobsForExecution(yaml.Jobs, debug)

	// assert
	jobResult1 := yamlResults.JobResults[0] // first job contains a cmd with an error
	jobResult2 := yamlResults.JobResults[1] // won't execute: continue_on_error is false!

	assert.True(t, jobResult1.Executed)
	assert.True(t, jobResult1.HasError)

	assert.False(t, jobResult2.Executed)
	assert.False(t, jobResult2.HasError)
}

func TestDispatchJobsForExecutionShouldExecuteJobs_ContinueOnErrorTrue_WithError_Sync(t *testing.T) {
	// arrange -> first job has an error, second has not!
	continueOnError := true
	yaml := testutils.NewJobFileTwoJobs_FirstJobWithError("job with error", 0, continueOnError) // sync jobs
	dispatcher := NewJobDispatcher(ExecuteJobSync, ExecuteJobAsync)
	debug := true

	// act
	yamlResults := dispatcher.DispatchJobsForExecution(yaml.Jobs, debug)

	// assert
	jobResult1 := yamlResults.JobResults[0] // first job contains a cmd with an error
	jobResult2 := yamlResults.JobResults[1] // will execute: continue_on_error is true!

	assert.True(t, jobResult1.Executed)
	assert.True(t, jobResult1.HasError)

	assert.True(t, jobResult2.Executed)  // executed
	assert.False(t, jobResult2.HasError) // executed without errors! (just first job has errors)
}

func TestDispatchJobsForExecutionShouldExecuteJobs_ContinueOnErrorFalse_WithError_Async(t *testing.T) {
	// arrange -> first job has an error, second has not!
	continueOnError := false
	yaml := testutils.NewJobFileTwoJobs_FirstJobWithError("job with error", 3, continueOnError) // async job: 3 goroutines
	dispatcher := NewJobDispatcher(ExecuteJobSync, ExecuteJobAsync)
	debug := true

	// act
	yamlResults := dispatcher.DispatchJobsForExecution(yaml.Jobs, debug)

	// assert
	jobResult1 := yamlResults.JobResults[0] // first job contains a cmd with an error
	jobResult2 := yamlResults.JobResults[1] // won't execute: continue_on_error is false!

	assert.True(t, jobResult1.Executed)
	assert.True(t, jobResult1.HasError)

	assert.False(t, jobResult2.Executed)
	assert.False(t, jobResult2.HasError)
}

func TestDispatchJobsForExecutionShouldExecuteJobs_ContinueOnErrorTrue_WithError_Async(t *testing.T) {
	// arrange -> first job has an error, second has not!
	continueOnError := true
	yaml := testutils.NewJobFileTwoJobs_FirstJobWithError("job with error", 3, continueOnError) // async job: 3 goroutines
	dispatcher := NewJobDispatcher(ExecuteJobSync, ExecuteJobAsync)
	debug := true

	// act
	yamlResults := dispatcher.DispatchJobsForExecution(yaml.Jobs, debug)

	// assert
	jobResult1 := yamlResults.JobResults[0] // first job contains a cmd with an error
	jobResult2 := yamlResults.JobResults[1] // will execute: continue_on_error is true!

	assert.True(t, jobResult1.Executed)
	assert.True(t, jobResult1.HasError)

	assert.True(t, jobResult2.Executed)  // executed
	assert.False(t, jobResult2.HasError) // executed without errors
}
