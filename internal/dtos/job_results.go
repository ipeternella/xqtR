package dtos

type CmdResult struct {
	StdoutData []byte
	StderrData []byte
	Err        error
}

type JobStepResult struct {
	Id        int     // JobId
	JobStep   JobStep // the job step data
	CmdResult *CmdResult
	Executed  bool
	HasError  bool
}

type JobResult struct {
	Title        string
	StepsResults []JobStepResult
	Executed     bool
	HasError     bool
}

type JobsYamlResult struct {
	JobResults []JobResult
}

func NewJobsYamlResult(jobResults []JobResult) *JobsYamlResult {
	return &JobsYamlResult{JobResults: jobResults}
}

func NewCmdResult(stdoutData []byte, stderrData []byte, err error) CmdResult {
	return CmdResult{
		StdoutData: stdoutData,
		StderrData: stderrData,
		Err:        err,
	}
}

func NewJobStepResult(id int, jobStep JobStep, cmdResult *CmdResult, executed bool, hasError bool) JobStepResult {
	return JobStepResult{
		Id:        id,
		JobStep:   jobStep,
		CmdResult: cmdResult,
		Executed:  executed,
		HasError:  hasError,
	}
}

// Empty structs
func NewEmptyJobStepResult(id int, jobStep JobStep) JobStepResult {
	return JobStepResult{
		Id:        id,
		JobStep:   jobStep,
		CmdResult: nil,
		Executed:  false,
		HasError:  false,
	}
}

func NewEmptyJobResult(job Job) JobResult {
	stepResults := []JobStepResult{}

	// initialize step results with just the step id and its jobStep data (but no cmd results, not executed, etc.)
	for id, jobStep := range job.Steps {
		stepResults = append(stepResults, NewEmptyJobStepResult(id, jobStep))
	}

	return JobResult{
		Title:        job.Title,
		StepsResults: stepResults,
	}
}

func NewEmptyJobsYamlResult(jobs []Job) JobsYamlResult {
	jobResults := []JobResult{}

	for _, job := range jobs {
		jobResults = append(jobResults, NewEmptyJobResult(job))
	}

	return JobsYamlResult{JobResults: jobResults}
}
