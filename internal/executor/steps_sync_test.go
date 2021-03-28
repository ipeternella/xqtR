package executor

import (
	"testing"

	"github.com/IgooorGP/xqtR/internal/dtos"
	"github.com/IgooorGP/xqtR/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestShouldExecuteJobsSyncWithoutErrorsNoDebug(t *testing.T) {
	// arrange
	yaml := testutils.NewJobsFileWithTwoSyncTasks()
	debug := false
	job1 := yaml.Jobs[0]

	// act
	jobResult := ExecuteJobSync(job1, debug)

	// assert
	expectedStepResults := []dtos.JobStepResult{
		{
			Id:      1,
			JobStep: yaml.Jobs[0].Steps[0],
			CmdResult: &dtos.CmdResult{
				StdoutData: []byte{}, // empty as debug is false: no stdout pipe reading
				StderrData: []byte{},
				Err:        nil,
			},
			Executed: true,
			HasError: false,
		},
		{
			Id:      2,
			JobStep: yaml.Jobs[0].Steps[1],
			CmdResult: &dtos.CmdResult{
				StdoutData: []byte{}, // empty as debug is false: no stdout pipe reading
				StderrData: []byte{},
				Err:        nil,
			},
			Executed: true,
			HasError: false,
		},
	}

	assert.Equal(t, jobResult.Title, "job name 1")
	assert.Equal(t, expectedStepResults, jobResult.StepsResults)
}

func TestShouldExecuteJobsSyncWithoutErrorsDebug(t *testing.T) {
	// arrange
	yaml := testutils.NewJobsFileWithTwoSyncTasks()
	debug := true
	job1 := yaml.Jobs[0]

	// act
	jobResult := ExecuteJobSync(job1, debug)

	// assert
	expectedStepResults := []dtos.JobStepResult{
		{
			Id:      1,
			JobStep: yaml.Jobs[0].Steps[0],
			CmdResult: &dtos.CmdResult{
				StdoutData: []byte("hello world\n"),
				StderrData: []byte{},
				Err:        nil,
			},
			Executed: true,
			HasError: false,
		},
		{
			Id:      2,
			JobStep: yaml.Jobs[0].Steps[1],
			CmdResult: &dtos.CmdResult{
				StdoutData: []byte("hi there\n"),
				StderrData: []byte{},
				Err:        nil,
			},
			Executed: true,
			HasError: false,
		},
	}

	assert.Equal(t, jobResult.Title, "job name 1")
	assert.Equal(t, expectedStepResults, jobResult.StepsResults)
}
