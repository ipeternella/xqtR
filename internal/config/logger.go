package config

import "github.com/rs/zerolog"

// ParseLevel parses log level strings into ZeroLog's constants
func ParseLevel(level string) zerolog.Level {
	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		return zerolog.Disabled // wrong setting
	}
}
