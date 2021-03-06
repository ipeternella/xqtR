// Package ui provides some pretty string formatting that's used to display commands results
// such as stdout, stderr, warnings, etc.
package ui

import "github.com/rs/zerolog/log"

const (
	processWarningHeader = "\n>--- Warnings: ---<\n\n"
	processWarningFooter = "\n>-----------------<"
	processStdoutHeader  = "\n>--- Stdout: ---<\n\n"
	processStdoutFooter  = "\n>---------------<"
	processStderrHeader  = "\n>--- Stderr: ---<\n\n"
	processStderrFooter  = "\n>---------------<"
)

func PrintCmdFailure(stepName string, stdoutData []byte, stderrData []byte, continueOnError bool) {
	log.Error().Msgf("⌛ step: %s 💀", stepName)
	log.Error().Msgf("%s%s%s", processStderrHeader, stderrData, processStderrFooter)

	// does not exit here anymore...
	// if !continueOnError {
	// 	log.Fatal().Msg("Step has failed and 'continue_on_error' is false. Exiting...")
	// }
}

func PrintCmdFeedback(stepName string, stdoutData []byte, stderrData []byte, debug bool) {
	log.Info().Msgf("⌛ step: %s 👍", stepName)

	// stderr is also used for warnings when the process does not exit with a non-zero status code
	if len(stderrData) > 0 {
		log.Warn().Msgf("%s%s%s", processWarningHeader, stderrData, processWarningFooter)
	}

	// stdout is print only if debug is on
	if debug && len(stdoutData) > 0 {
		log.Debug().Msgf("%s%s%s", processStdoutHeader, stdoutData, processStdoutFooter)
	}
}
