package executor

import (
	"io"
	"io/ioutil"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func NewBytesReadCloserFromString(str string) *io.ReadCloser {
	var readerCloserPtr *io.ReadCloser

	bytesReaderCloser := ioutil.NopCloser(strings.NewReader(str))
	readerCloserPtr = &bytesReaderCloser

	return readerCloserPtr
}

func TestShouldReadAllPipeData(t *testing.T) {
	// arrange
	var allPipeData = "Some command stdout bytes array"
	cmdPipe := NewBytesReadCloserFromString(allPipeData)

	// act
	readPipeData := readPipe(cmdPipe)

	// assert
	assert.Equal(t, []byte(allPipeData), readPipeData)
}

func TestShouldReadAllPipeDataFromCmdStdoutAndStderr(t *testing.T) {
	// arrange
	var readStdout = true

	stdoutPipe := NewBytesReadCloserFromString("Some command stdout bytes")
	stderrPipe := NewBytesReadCloserFromString("")

	// act
	stdoutData, stderrData := readCmdStdStreams(stdoutPipe, stderrPipe, readStdout)

	// assert
	var expectedStdoutData = "Some command stdout bytes"
	var expectedStderrData = ""

	assert.Equal(t, []byte(expectedStdoutData), stdoutData)
	assert.Equal(t, []byte(expectedStderrData), stderrData)
}

func TestShouldReadAllPipeDataFromStderrOnly(t *testing.T) {
	// arrange
	var readStdout = false // no stdout reading

	stdoutPipe := NewBytesReadCloserFromString("Some command stdout bytes")
	stderrPipe := NewBytesReadCloserFromString("Some command stderr bytes")

	// act
	stdoutData, stderrData := readCmdStdStreams(stdoutPipe, stderrPipe, readStdout)

	// assert
	var expectedStdoutData = "" // no stdout reading
	var expectedStderrData = "Some command stderr bytes"

	assert.Equal(t, []byte(expectedStdoutData), stdoutData)
	assert.Equal(t, []byte(expectedStderrData), stderrData)
}

func TestShouldCreateShellCommand(t *testing.T) {
	// arrange
	echoCmd := `echo "hello world!"`

	// act
	shellCmd, _, _ := shellCommand(echoCmd)

	// assert
	expectedArgs := []string{"bash", "-c", "echo \"hello world!\""}

	assert.Regexp(t, regexp.MustCompile(`(/bin/bash|/usr/bin/bash)`), shellCmd.Path)
	assert.Equal(t, expectedArgs, shellCmd.Args) // bash is the shell to execute the cmd
}
