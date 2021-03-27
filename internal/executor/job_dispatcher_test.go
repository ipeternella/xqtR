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

	syncJobExecutorMock := func(job dtos.Job, debug bool) *dtos.JobResult {
		var mockJobSteps []dtos.JobStep
		syncCalled = true

		mockResult := dtos.Job{Title: "mock", Steps: mockJobSteps}

		return dtos.NewEmptyJobResult(mockResult)
	}

	asyncJobExecutorMock := func(job dtos.Job, debug bool) *dtos.JobResult {
		var mockJobSteps []dtos.JobStep
		syncCalled = true

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

	syncJobExecutorMock := func(job dtos.Job, debug bool) *dtos.JobResult {
		var mockJobSteps []dtos.JobStep
		syncCalled = true

		mockResult := dtos.Job{Title: "mock", Steps: mockJobSteps}

		return dtos.NewEmptyJobResult(mockResult)
	}

	asyncJobExecutorMock := func(job dtos.Job, debug bool) *dtos.JobResult {
		var mockJobSteps []dtos.JobStep
		syncCalled = true

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

	syncJobExecutorMock := func(job dtos.Job, debug bool) *dtos.JobResult {
		var mockJobSteps []dtos.JobStep
		syncCalled = true

		mockResult := dtos.Job{Title: "mock", Steps: mockJobSteps}

		return dtos.NewEmptyJobResult(mockResult)
	}

	asyncJobExecutorMock := func(job dtos.Job, debug bool) *dtos.JobResult {
		var mockJobSteps []dtos.JobStep
		syncCalled = true

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
