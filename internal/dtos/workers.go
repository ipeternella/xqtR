package dtos

type CmdResult struct {
	StdoutData []byte
	StderrData []byte
	Err        error
}

type WorkerData struct {
	Id      int
	JobStep JobStep
	Debug   bool
}

type WorkerResult struct {
	WorkerId int
	Name     string
	Result   *CmdResult
}
