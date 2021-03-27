package dtos

type CmdResult struct {
	StdoutData []byte
	StderrData []byte
	Err        error
}

type WorkerData struct {
	StepId            int
	JobStep           JobStep
	JobExecutionRules JobExecutionRules
}

// type WorkerResult struct {
// 	StepId int
// 	Name   string // jobStep name
// 	Result *CmdResult
// }

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
