package executor

import (
	"github.com/IgooorGP/xqtR/internal/dtos"
	"github.com/rs/zerolog/log"
)

func newWorkerData(stepId int, jobStep dtos.JobStep, jobExecutionRules dtos.JobExecutionRules) *dtos.WorkerData {
	return &dtos.WorkerData{
		StepId:            stepId,
		JobStep:           jobStep,
		JobExecutionRules: jobExecutionRules,
	}
}

func newCmdResult(stdoutData []byte, stderrData []byte, err error) *dtos.CmdResult {
	return &dtos.CmdResult{
		StdoutData: stdoutData,
		StderrData: stderrData,
		Err:        err,
	}
}

// split jobSteps
func executeJobStepByWorker(workerResults chan<- dtos.JobStepResult, taskQueue <-chan *dtos.WorkerData) {
	// keep consuming from queue as long its opened (chan blocks if there are no tasks)
	for workerData := range taskQueue {
		var jobStepResult = dtos.NewEmptyJobStepResult(workerData.StepId, workerData.JobStep)
		var cmdResult *dtos.CmdResult

		cmd, cmdStdoutPipe, cmdStderrPipe := shellCommand(workerData.JobStep.Run)

		// spawns a new OS process with the cmd
		if err := cmd.Start(); err != nil {
			log.Fatal().Msgf("An error happened while starting the cmd: %s", err.Error())
		}

		stdoutData, stderrData := readCmdStdStreams(cmdStdoutPipe, cmdStderrPipe, workerData.JobExecutionRules.Debug)

		if err := cmd.Wait(); err != nil {
			cmdResult = newCmdResult(stdoutData, stderrData, err) // something bad happened
		} else {
			cmdResult = newCmdResult(stdoutData, stderrData, nil) // all good (hopefully!)
		}

		// publish back to main goroutine the cmd result
		// workerResult := newWorkerResult(workerData.StepId, workerData.JobStep.Name, cmdResult)
		// jobStepResult.CmdResult = cmdResult
		markStepAsExecuted(&jobStepResult, *cmdResult)

		workerResults <- jobStepResult
	}
}
