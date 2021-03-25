// Package tests contains code that should be imported only by test-code only in order
// not to be compiled into the final project's binary. This package contains mostly
// helper functions (such as mocks, etc.) to be used by *_test.go files.
package testutils

import "github.com/IgooorGP/xqtR/internal/dtos"

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