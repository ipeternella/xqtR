package dtos

type WorkerData struct {
	StepId            int
	JobStep           JobStep
	JobExecutionRules JobExecutionRules
}

type JobExecutionRules struct {
	Debug           bool
	ContinueOnError bool
}

func NewJobExecutionRules(debug bool, continueOnError bool) JobExecutionRules {
	return JobExecutionRules{
		Debug:           debug,
		ContinueOnError: continueOnError,
	}
}
