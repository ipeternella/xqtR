// Package tests contains code that should be imported only by test-code only in order
// not to be compiled into the final project's binary. This package contains mostly
// helper functions (such as mocks, etc.) to be used by *_test.go files.
package testutils

import (
	"fmt"

	"github.com/IgooorGP/xqtR/internal/dtos"
)

func NewMockJobStep(name string, run string) dtos.JobStep {
	return dtos.JobStep{
		Name: name,
		Run:  run,
	}
}

func NewMockJob(title string, steps []dtos.JobStep, numWorkers int) dtos.Job {
	return dtos.Job{
		Title:      title,
		Steps:      steps,
		NumWorkers: numWorkers,
	}
}

func NewMockJobsFile(jobs []dtos.Job) dtos.JobsFile {
	return dtos.JobsFile{
		Jobs: jobs,
	}
}

func NewJobsFileWithTwoSyncTasks() dtos.JobsFile {
	job1 := []dtos.JobStep{
		NewMockJobStep(`echo "hello world"`, `echo "hello world"`),
		NewMockJobStep(`echo "hi there"`, `echo "hi there"`),
	}

	job2 := []dtos.JobStep{
		NewMockJobStep("sleep for 1s", "sleep 1s"),
	}

	jobs := []dtos.Job{
		NewMockJob("job name 1", job1, 0), // sync job
		NewMockJob("job name 2", job2, 0), // sync job
	}

	// no continuing upon errors
	jobs[0].ContinueOnError = false
	jobs[1].ContinueOnError = false

	return NewMockJobsFile(jobs)
}

func NewSyncJobFileWithEchoStepError(jobTitle string, continueOnError bool) dtos.JobsFile {
	job := []dtos.JobStep{
		NewMockJobStep(`echo "hi"`, `echo "hi"`),                         // ok cmd!
		NewMockJobStep(`echo "hoi with error"`, `wcho "hoi with error"`), // typo on 'echo' -> 'wcho'
		NewMockJobStep(`echo "hey"`, `echo "hey"`),                       // ok cmd!
	}

	jobs := []dtos.Job{
		NewMockJob(jobTitle, job, 0), // sync job
	}

	// no continuing upon errors
	jobs[0].ContinueOnError = continueOnError

	return NewMockJobsFile(jobs)
}

func NewAsyncJobFileWithEchoStepError(jobTitle string, numWorkers int, continueOnError bool) dtos.JobsFile {
	job := []dtos.JobStep{
		NewMockJobStep(`echo "hi"`, `echo "hi"`),
		NewMockJobStep(`echo "hoi with error"`, `wcho "hoi with error"`), // typo on 'echo' -> 'wcho'
		NewMockJobStep(`echo "hey"`, `echo "hey"`),
	}

	jobs := []dtos.Job{
		NewMockJob(jobTitle, job, numWorkers), // sync job
	}

	// no continuing upon errors
	jobs[0].ContinueOnError = continueOnError

	return NewMockJobsFile(jobs)
}

func NewSingleJobFileBuilder(jobTitle string, stepNamePrefix string, stepCmd string, numSteps int, numWorkers int, continueOnError bool) dtos.JobsFile {
	jobSteps := []dtos.JobStep{}

	// build steps with same cmd, whose display name gets the counter i "suffix"
	for i := 0; i < numSteps; i++ {
		jobSteps = append(jobSteps, NewMockJobStep(fmt.Sprintf("%s - %d", stepNamePrefix, i), stepCmd))
	}

	jobs := []dtos.Job{
		NewMockJob(jobTitle, jobSteps, numWorkers),
	}

	jobs[0].ContinueOnError = continueOnError

	return NewMockJobsFile(jobs)
}

func NewMockJobsFileWithSyncAndAsyncJobs() dtos.JobsFile {

	job1 := []dtos.JobStep{
		NewMockJobStep("job1.step1", "job1.step1.run11"),
		NewMockJobStep("job1.step2", "job1.step2.run12"),
	}

	job2 := []dtos.JobStep{
		NewMockJobStep("job2.step1", "job2.step1.run21"),
		NewMockJobStep("job2.step2", "job2.step2.run22"),
	}

	jobs := []dtos.Job{
		NewMockJob("job1", job1, 0), // sync job
		NewMockJob("job2", job2, 3), // async job with 3 workers
	}

	return NewMockJobsFile(jobs)
}

func NewMockJobsFileWithoutNumWorkers() dtos.JobsFile {

	job2 := []dtos.JobStep{
		NewMockJobStep("job1.step1", "job1.step1.run11"),
		NewMockJobStep("job1.step2", "job1.step2.run12"),
		NewMockJobStep("job1.step3", "job1.step3.run13"),
	}

	jobs := []dtos.Job{
		NewMockJob("job1", job2, 0), // sync job
	}

	return NewMockJobsFile(jobs)
}

func NewMockJobsFileWithNumWorkers() dtos.JobsFile {

	job1 := []dtos.JobStep{
		NewMockJobStep("job1.step1", "job1.step1.run11"),
		NewMockJobStep("job1.step2", "job1.step2.run12"),
		NewMockJobStep("job1.step3", "job1.step3.run13"),
	}

	jobs := []dtos.Job{
		NewMockJob("job1", job1, 3), // async: 3 goroutines
	}

	return NewMockJobsFile(jobs)
}
