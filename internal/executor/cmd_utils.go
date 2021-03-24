package executor

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
)

// ExtractCommandAndArgsFromString - reads a string and parses it to a terminal command
func shellCommand(command string) (*exec.Cmd, *io.ReadCloser, *io.ReadCloser) {
	cmd := exec.Command("bash", "-c", command) // TODO: improve

	// adds an os.Pipe to the cmd stdstreams so that xqtR's main process can read the
	// the stdout and stderr of the new cmd spawned process
	stdoutPipe, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()

	return cmd, &stdoutPipe, &stderrPipe
}

func readPipe(cmdPipe *io.ReadCloser) []byte {
	bufferData, err := ioutil.ReadAll(*cmdPipe)

	if err != nil {
		log.Error().Msgf("Error when reading cmd buffer: %s", err.Error())
		os.Exit(1)
	}

	return bufferData
}

func readCmdStdStreams(cmdStdoutPipe *io.ReadCloser, cmdStderrPipe *io.ReadCloser, readStdout bool) ([]byte, []byte) {
	var stdoutData = []byte("")
	var stderrData = []byte("")

	if readStdout {
		stdoutData = readPipe(cmdStdoutPipe)
	}
	stderrData = readPipe(cmdStderrPipe)

	return stdoutData, stderrData
}
