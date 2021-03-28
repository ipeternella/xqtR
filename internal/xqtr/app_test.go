package xqtr

import (
	"testing"

	"github.com/IgooorGP/xqtR/internal/config"
)

const (
	yamlBasePath = "../testutils/testdata/"
)

func TestXqtRAppShoudRunJobYamlWithoutErrors(t *testing.T) {
	// arrange
	cfg := config.NewXqtRConfigWithDefaults()
	cfg.JobFilePath = yamlBasePath + "job_app_test.yaml"
	app := NewXqtR(&cfg) // core app instance

	// act and assert -> should not raise any errors
	app.Run()
}
