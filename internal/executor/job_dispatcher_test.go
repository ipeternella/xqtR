package executor

import (
	"testing"

	"github.com/IgooorGP/xqtR/internal/dtos"
	"github.com/IgooorGP/xqtR/internal/tests"
	"github.com/stretchr/testify/assert"
)

func TestDispatchForSyncJob(t *testing.T) {
	// arrange - mocks
	syncCalled := false
	asyncCalled := false

	syncJobExecutorMock := func(job dtos.Job, debug bool) {
		syncCalled = true
	}

	asyncJobExecutorMock := func(job dtos.Job, debug bool) {
		asyncCalled = true
	}

	yaml := tests.MockJobsFileSyncOnly()
	debug := true
	dispatcher := NewJobDispatcher(syncJobExecutorMock, asyncJobExecutorMock)

	// act
	dispatcher.DispatchJobsForExecution(yaml.Jobs, debug)

	// arrange
	assert.True(t, syncCalled)
	assert.False(t, asyncCalled)
}
