package config

type XqtRConfig struct {
	JobFilePath string
	LogLevel    string
	IsDryRun    bool
}

// NewXqtRConfig creates a new config struct
func NewXqtRConfig(jobFilePath string, isDryRun bool, logLevel string) XqtRConfig {
	return XqtRConfig{JobFilePath: jobFilePath, IsDryRun: isDryRun, LogLevel: logLevel}
}

func NewXqtRConfigWithDefaults() XqtRConfig {
	return XqtRConfig{JobFilePath: "./job.yaml", IsDryRun: false, LogLevel: "info"}
}
