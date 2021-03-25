package xqtr

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/IgooorGP/xqtR/internal/config"
	"github.com/IgooorGP/xqtR/internal/startup"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	cfg := config.NewXqtRConfigWithDefaults()
	startup.Boot(&cfg) // log level as info

	code := m.Run()
	os.Exit(code)
}

const (
	brokenYamlFolderPath = "../testutils/testdata/" // move out of 'xqtr' pkg folder
)

func TestYamlSchemaValidationShouldBreakWhenYamlHasMissingStepName(t *testing.T) {
	// arrange
	brokenYamlRelativePath := brokenYamlFolderPath + "broken_no_step_name.yaml" // missing name of a step on yaml
	brokenYamlPath_noJobs, _ := filepath.Abs(brokenYamlRelativePath)

	viperParser := prepareViper(brokenYamlPath_noJobs)
	yaml := parseYaml(viperParser, brokenYamlPath_noJobs)

	// act
	err := validateYamlSchema(yaml)

	// assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), `missing "name" key for step on (job name: "job 1", job number: 1, step number: 3)`)
}

func TestYamlSchemaValidationShouldBreakWhenYamlHasMissingStepRun(t *testing.T) {
	// arrange
	brokenYamlRelativePath := brokenYamlFolderPath + "broken_no_step_run.yaml" // missing run of a step on yaml
	brokenYamlPath_noJobs, _ := filepath.Abs(brokenYamlRelativePath)

	viperParser := prepareViper(brokenYamlPath_noJobs)
	yaml := parseYaml(viperParser, brokenYamlPath_noJobs)

	// act
	err := validateYamlSchema(yaml)

	// assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), `missing "run" key for step on (job name: "job 2", job number: 2, step name: "task 22", step number: 2)`)
}

func TestYamlSchemaValidationShouldBreakWhenYamlHasJobWithoutSteps(t *testing.T) {
	// arrange
	brokenYamlRelativePath := brokenYamlFolderPath + "broken_job_without_steps.yaml" // job with just a "title" but no "steps" list
	brokenYamlPath_noJobs, _ := filepath.Abs(brokenYamlRelativePath)

	viperParser := prepareViper(brokenYamlPath_noJobs)
	yaml := parseYaml(viperParser, brokenYamlPath_noJobs)

	// act
	err := validateYamlSchema(yaml)

	// assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), `missing "steps" key for (job name: "job 2", job number: 2)`)
}

func TestYamlSchemaValidationShouldBreakWhenYamlHasStepNameTypo(t *testing.T) {
	// arrange
	brokenYamlRelativePath := brokenYamlFolderPath + "broken_step_name_typo.yaml" // job with just a typo on the "name" key
	brokenYamlPath_noJobs, _ := filepath.Abs(brokenYamlRelativePath)

	viperParser := prepareViper(brokenYamlPath_noJobs)
	yaml := parseYaml(viperParser, brokenYamlPath_noJobs)

	// act
	err := validateYamlSchema(yaml)

	// assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), `missing "name" key for step on (job name: "job 3", job number: 3, step number: 3)`)
}

func TestYamlSchemaValidationShouldBreakWhenYamlHasStepRunTypo(t *testing.T) {
	// arrange
	brokenYamlRelativePath := brokenYamlFolderPath + "broken_step_run_typo.yaml" // job with just a typo on the "run" key
	brokenYamlPath_noJobs, _ := filepath.Abs(brokenYamlRelativePath)

	viperParser := prepareViper(brokenYamlPath_noJobs)
	yaml := parseYaml(viperParser, brokenYamlPath_noJobs)

	// act
	err := validateYamlSchema(yaml)

	// assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), `missing "run" key for step on (job name: "job 2", job number: 2, step name: "task 21", step number: 1)`)
}

func TestYamlSchemaValidationShouldNotBreakWithNormalJobYaml(t *testing.T) {
	// arrange
	brokenYamlRelativePath := brokenYamlFolderPath + "job_not_broken.yaml" // yaml is fine: no schema errors!
	brokenYamlPath_noJobs, _ := filepath.Abs(brokenYamlRelativePath)

	viperParser := prepareViper(brokenYamlPath_noJobs)
	yaml := parseYaml(viperParser, brokenYamlPath_noJobs)

	// act
	err := validateYamlSchema(yaml)

	// assert
	assert.Nil(t, err)
}
