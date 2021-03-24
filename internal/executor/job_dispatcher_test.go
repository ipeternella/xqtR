package executor

import (
	"testing"

	"github.com/IgooorGP/xqtR/internal/dtos"
	"github.com/IgooorGP/xqtR/internal/tests"
	"github.com/stretchr/testify/assert"
)

func TestDispatchForSyncJobOnly(t *testing.T) {
	// arrange - mocks
	syncCalled := false
	asyncCalled := false

	syncJobExecutorMock := func(job dtos.Job, debug bool) {
		syncCalled = true
	}

	asyncJobExecutorMock := func(job dtos.Job, debug bool) {
		asyncCalled = true
	}

	yaml := tests.NewMockJobsFileWithoutNumWorkers() // sync jobs only
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

	syncJobExecutorMock := func(job dtos.Job, debug bool) {
		syncCalled = true
	}

	asyncJobExecutorMock := func(job dtos.Job, debug bool) {
		asyncCalled = true
	}

	yaml := tests.NewMockJobsFileWithNumWorkers() // async job only
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

	syncJobExecutorMock := func(job dtos.Job, debug bool) {
		syncCalled = true
	}

	asyncJobExecutorMock := func(job dtos.Job, debug bool) {
		asyncCalled = true
	}

	yaml := tests.NewMockJobsFileWithSyncAndAsyncJobs() // async job only
	debug := true
	dispatcher := NewJobDispatcher(syncJobExecutorMock, asyncJobExecutorMock)

	// act
	dispatcher.DispatchJobsForExecution(yaml.Jobs, debug)

	// arrange
	assert.True(t, syncCalled)
	assert.True(t, asyncCalled)
}
