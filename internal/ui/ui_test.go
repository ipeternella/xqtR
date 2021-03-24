package ui

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestShouldPrintCmdFeedbackOnlyStderr(t *testing.T) {
	// arrange -> set logger output to local buffer
	var logBuffer bytes.Buffer
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: &logBuffer, TimeFormat: time.Kitchen})

	// arrange mock stdout and stderr
	stdoutMock := []byte("stdout data!")
	stderrMock := []byte("stderr data!")
	stepNameMock := "Mock step!"
	debug := false // no stdout

	// act
	PrintCmdFeedback(stepNameMock, stdoutMock, stderrMock, debug)

	// assert
	processStdstreamData := logBuffer.String()
	expectedStderrData := fmt.Sprintf("%s%s%s", processWarningHeader, stderrMock, processWarningFooter)
	expectedStdoutData := fmt.Sprintf("%s%s%s", processStdoutHeader, stdoutMock, processStdoutFooter)

	assert.Contains(t, processStdstreamData, expectedStderrData)
	assert.NotContains(t, processStdstreamData, expectedStdoutData) // debug is false, so not stdout data, just stderr!
}

func TestShouldPrintCmdFeedbackStderrAndStdout(t *testing.T) {
	// arrange -> set logger output to local buffer
	var logBuffer bytes.Buffer
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: &logBuffer, TimeFormat: time.Kitchen})

	// arrange mock stdout and stderr
	stdoutMock := []byte("stdout data!")
	stderrMock := []byte("stderr data!")
	stepNameMock := "Mock step!"
	debug := true // stdout will be print!

	// act
	PrintCmdFeedback(stepNameMock, stdoutMock, stderrMock, debug)

	// assert
	processStdstreamData := logBuffer.String()
	expectedStderrData := fmt.Sprintf("%s%s%s", processWarningHeader, stderrMock, processWarningFooter)
	expectedStdoutData := fmt.Sprintf("%s%s%s", processStdoutHeader, stdoutMock, processStdoutFooter)

	assert.Contains(t, processStdstreamData, expectedStderrData)
	assert.Contains(t, processStdstreamData, expectedStdoutData) // debug is true, so stdout should have be sent
}
