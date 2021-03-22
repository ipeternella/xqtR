package app

import (
	"github.com/rs/zerolog/log"
)

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

func newWorkerData(id int, jobStep JobStep, debug bool) *WorkerData {
	return &WorkerData{
		Id:      id,
		JobStep: jobStep,
		Debug:   debug,
	}
}

func newCmdResult(stdoutData []byte, stderrData []byte, err error) *CmdResult {
	return &CmdResult{
		StdoutData: stdoutData,
		StderrData: stderrData,
		Err:        err,
	}
}

func newWorkerResult(workerId int, name string, cmdResult *CmdResult) *WorkerResult {
	return &WorkerResult{
		WorkerId: workerId,
		Name:     name,
		Result:   cmdResult,
	}
}

// split jobSteps
func executeJobStepByWorker(workerResults chan<- *WorkerResult, taskQueue <-chan *WorkerData) {

	// keep consuming from queue as long its opened (chan blocks if there are no tasks)
	for workerData := range taskQueue {
		log.Info().Msgf("⏳ step: %s", workerData.JobStep.Name)

		var rslt *CmdResult
		cmd, cmdStdoutPipe, cmdStderrPipe := shellCommand(workerData.JobStep.Run)

		// spawns a new OS process with the cmd
		if err := cmd.Start(); err != nil {
			log.Fatal().Msgf("An error happened while starting the cmd", err.Error())
		}

		stdoutData, stderrData := readCmdStdStreams(cmdStdoutPipe, cmdStderrPipe, workerData.Debug)

		if err := cmd.Wait(); err != nil {
			rslt = newCmdResult(stdoutData, stderrData, err) // something bad happened
		} else {
			rslt = newCmdResult(stdoutData, stderrData, nil) // all good (hopefully!)
		}

		// publish back to main goroutine the cmd result
		workerResult := newWorkerResult(workerData.Id, workerData.JobStep.Name, rslt)
		workerResults <- workerResult
	}
}
