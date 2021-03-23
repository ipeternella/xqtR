// Package config holds the main configuration structure that is used by XqtR app
// instances to find job yaml files, to configure loggers in the startup package, etc.
package config

type XqtRConfig struct {
	JobFilePath string
	LogLevel    string
	IsDryRun    bool
}

func NewXqtRConfigWithDefaults() XqtRConfig {
	return XqtRConfig{JobFilePath: "./job.yaml", IsDryRun: false, LogLevel: "info"}
}
