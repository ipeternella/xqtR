package executor

import (
	"testing"

	"github.com/IgooorGP/xqtR/internal/dtos"
	"github.com/IgooorGP/xqtR/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestShouldExecuteAsyncJobWithoutErrorsDebugOff(t *testing.T) {
	// arrange - echo cmds with 3 workers (async job)
	yaml := testutils.NewSingleJobFileBuilder("some async job 1", "echoing async", `echo "hello world"`, 3, 3, false)
	job := yaml.Jobs[0]
	debug := false

	// act
	jobResult := ExecuteJobAsync(job, debug)

	// assert
	expectedStepResults := []dtos.JobStepResult{
		{
			Id:      0,
			JobStep: yaml.Jobs[0].Steps[0],
			CmdResult: &dtos.CmdResult{
				StdoutData: []byte{}, // no stdout -> debug is false
				StderrData: []byte{},
				Err:        jobResult.StepsResults[0].CmdResult.Err,
			},
			Executed: true,
			HasError: false,
		},
		{
			Id:      1,
			JobStep: yaml.Jobs[0].Steps[1],
			CmdResult: &dtos.CmdResult{
				StdoutData: []byte{},
				StderrData: []byte{},
				Err:        jobResult.StepsResults[0].CmdResult.Err,
			},
			Executed: true,
			HasError: false,
		},
		{
			Id:      2,
			JobStep: yaml.Jobs[0].Steps[2],
			CmdResult: &dtos.CmdResult{
				StdoutData: []byte{},
				StderrData: []byte{},
				Err:        jobResult.StepsResults[0].CmdResult.Err,
			},
			Executed: true,
			HasError: false,
		},
	}

	assert.Equal(t, jobResult.Title, "some async job 1")
	assert.Equal(t, expectedStepResults, jobResult.StepsResults)
}

func TestShouldExecuteAsyncJobWithoutErrorsDebugOn(t *testing.T) {
	// arrange - echo cmds with 3 workers (async job)
	yaml := testutils.NewSingleJobFileBuilder("some async job 1", "echoing async", `echo "hello world"`, 3, 3, false)
	job := yaml.Jobs[0]
	debug := true

	// act
	jobResult := ExecuteJobAsync(job, debug)

	// assert
	expectedStepResults := []dtos.JobStepResult{
		{
			Id:      0,
			JobStep: yaml.Jobs[0].Steps[0],
			CmdResult: &dtos.CmdResult{
				StdoutData: []byte("hello world\n"),
				StderrData: []byte{},
				Err:        jobResult.StepsResults[0].CmdResult.Err,
			},
			Executed: true,
			HasError: false,
		},
		{
			Id:      1,
			JobStep: yaml.Jobs[0].Steps[1],
			CmdResult: &dtos.CmdResult{
				StdoutData: []byte("hello world\n"),
				StderrData: []byte{},
				Err:        jobResult.StepsResults[0].CmdResult.Err,
			},
			Executed: true,
			HasError: false,
		},
		{
			Id:      2,
			JobStep: yaml.Jobs[0].Steps[2],
			CmdResult: &dtos.CmdResult{
				StdoutData: []byte("hello world\n"),
				StderrData: []byte{},
				Err:        jobResult.StepsResults[0].CmdResult.Err,
			},
			Executed: true,
			HasError: false,
		},
	}

	assert.Equal(t, jobResult.Title, "some async job 1")
	assert.Equal(t, expectedStepResults, jobResult.StepsResults)
}

func TestShouldExecuteAsyncJobWithError_ContinueOnErrorTrue(t *testing.T) {
	// arrange - echo cmds with 3 workers (async job)
	// continueOnError doesn't matter for async jobs: they are executed concurrently
	continueOnError := true
	yaml := testutils.NewAsyncJobFileWithEchoStepError("some async job with error 1", 3, continueOnError)
	job := yaml.Jobs[0]
	debug := true

	// act
	jobResult := ExecuteJobAsync(job, debug)

	// assert
	expectedStepResults := []dtos.JobStepResult{
		{
			Id:      0,
			JobStep: yaml.Jobs[0].Steps[0],
			CmdResult: &dtos.CmdResult{
				StdoutData: []byte("hi\n"),
				StderrData: []byte{},
				Err:        nil,
			},
			Executed: true,
			HasError: false,
		},
		// step id = 2 has a typo error on 'echo' cmd -> only step with 'HasError': true
		{
			Id:      1,
			JobStep: yaml.Jobs[0].Steps[1],
			CmdResult: &dtos.CmdResult{
				StdoutData: []byte{}, // no stdout: there was an error
				StderrData: []byte("bash: wcho: command not found\n"),
				Err:        jobResult.StepsResults[1].CmdResult.Err, // stepResult 1 error
			},
			Executed: true,
			HasError: true,
		},
		{
			Id:      2,
			JobStep: yaml.Jobs[0].Steps[2],
			CmdResult: &dtos.CmdResult{
				StdoutData: []byte("hey\n"),
				StderrData: []byte{},
				Err:        nil,
			},
			Executed: true,
			HasError: false,
		},
	}

	assert.Equal(t, jobResult.Title, "some async job with error 1")
	assert.Equal(t, expectedStepResults, jobResult.StepsResults)
}

func TestShouldExecuteAsyncJobWithError_ContinueOnErrorFalse(t *testing.T) {
	// arrange - echo cmds with 3 workers (async job)

	// continueOnError doesn't matter for async jobs: they are executed concurrently
	continueOnError := false
	yaml := testutils.NewAsyncJobFileWithEchoStepError("some async job with error 1", 3, continueOnError)
	job := yaml.Jobs[0]
	debug := true

	// act
	jobResult := ExecuteJobAsync(job, debug)

	// assert
	expectedStepResults := []dtos.JobStepResult{
		{
			Id:      0,
			JobStep: yaml.Jobs[0].Steps[0],
			CmdResult: &dtos.CmdResult{
				StdoutData: []byte("hi\n"),
				StderrData: []byte{},
				Err:        nil,
			},
			Executed: true,
			HasError: false,
		},
		// step id = 2 has a typo error on 'echo' cmd -> only step with 'HasError': true
		{
			Id:      1,
			JobStep: yaml.Jobs[0].Steps[1],
			CmdResult: &dtos.CmdResult{
				StdoutData: []byte{}, // no stdout: there was an error
				StderrData: []byte("bash: wcho: command not found\n"),
				Err:        jobResult.StepsResults[1].CmdResult.Err, // stepResult 1 error
			},
			Executed: true,
			HasError: true,
		},
		{
			Id:      2,
			JobStep: yaml.Jobs[0].Steps[2],
			CmdResult: &dtos.CmdResult{
				StdoutData: []byte("hey\n"),
				StderrData: []byte{},
				Err:        nil,
			},
			Executed: true,
			HasError: false,
		},
	}

	assert.Equal(t, jobResult.Title, "some async job with error 1")
	assert.Equal(t, expectedStepResults, jobResult.StepsResults)
}
