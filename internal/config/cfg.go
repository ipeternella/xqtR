package config

type XqtRConfig struct {
	JobFilePath string
	LogLevel    string
	IsDryRun    bool
}

func NewXqtRConfigWithDefaults() XqtRConfig {
	return XqtRConfig{JobFilePath: "./job.yaml", IsDryRun: false, LogLevel: "info"}
}
