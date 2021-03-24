package executor

import (
	"github.com/IgooorGP/xqtR/internal/dtos"
	"github.com/rs/zerolog/log"
)

func newWorkerData(id int, jobStep dtos.JobStep, debug bool) *dtos.WorkerData {
	return &dtos.WorkerData{
		Id:      id,
		JobStep: jobStep,
		Debug:   debug,
	}
}

func newCmdResult(stdoutData []byte, stderrData []byte, err error) *dtos.CmdResult {
	return &dtos.CmdResult{
		StdoutData: stdoutData,
		StderrData: stderrData,
		Err:        err,
	}
}

func newWorkerResult(workerId int, name string, cmdResult *dtos.CmdResult) *dtos.WorkerResult {
	return &dtos.WorkerResult{
		WorkerId: workerId,
		Name:     name,
		Result:   cmdResult,
	}
}

// split jobSteps
func executeJobStepByWorker(workerResults chan<- *dtos.WorkerResult, taskQueue <-chan *dtos.WorkerData) {

	// keep consuming from queue as long its opened (chan blocks if there are no tasks)
	for workerData := range taskQueue {
		log.Info().Msgf("⏳ step: %s", workerData.JobStep.Name)

		var rslt *dtos.CmdResult
		cmd, cmdStdoutPipe, cmdStderrPipe := shellCommand(workerData.JobStep.Run)

		// spawns a new OS process with the cmd
		if err := cmd.Start(); err != nil {
			log.Fatal().Msgf("An error happened while starting the cmd: %s", err.Error())
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
