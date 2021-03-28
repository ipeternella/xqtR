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
			Id:      0,
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
			Id:      1,
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
			Id:      0,
			JobStep: yaml.Jobs[0].Steps[0],
			CmdResult: &dtos.CmdResult{
				StdoutData: []byte("hello world\n"), // stdout is read due to debug = true
				StderrData: []byte{},
				Err:        nil,
			},
			Executed: true,
			HasError: false,
		},
		{
			Id:      1,
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

func TestShouldExecuteSyncJobWithStepWithContinueOnErrorFalse(t *testing.T) {
	// arrange - cmd with a typo 'wcho' not 'echo' with 0 workers (sync job)
	yaml := testutils.NewSingleJobFileBuilder("some job 1", "echoing", "wcho 'hello world'", 3, 0, false)
	job := yaml.Jobs[0]
	debug := true

	// act
	jobResult := ExecuteJobSync(job, debug)

	// assert
	expectedStepResults := []dtos.JobStepResult{
		{
			Id:      0,
			JobStep: yaml.Jobs[0].Steps[0],
			CmdResult: &dtos.CmdResult{
				StdoutData: []byte{},
				StderrData: []byte("bash: wcho: command not found\n"),
				Err:        jobResult.StepsResults[0].CmdResult.Err,
			},
			Executed: true,
			HasError: true,
		},
		// since the first step has errors, the others aren't executed
		{
			Id:        1,
			JobStep:   yaml.Jobs[0].Steps[1],
			CmdResult: nil,
			Executed:  false,
			HasError:  false,
		},
		{
			Id:        2,
			JobStep:   yaml.Jobs[0].Steps[2],
			CmdResult: nil,
			Executed:  false,
			HasError:  false,
		},
	}

	assert.Equal(t, jobResult.Title, "some job 1")
	assert.Equal(t, expectedStepResults, jobResult.StepsResults)
}
