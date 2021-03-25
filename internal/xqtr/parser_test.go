package xqtr

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	brokenYamlFolderPath = "../testutils/testdata/" // move out of 'xqtr' pkg folder
)

func TestShouldBreakWhenYamlHasMissingStepName(t *testing.T) {
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

func TestShouldBreakWhenYamlHasMissingStepRun(t *testing.T) {
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

func TestShouldBreakWhenYamlHasJobWithoutSteps(t *testing.T) {
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

func TestShouldBreakWhenYamlHasStepNameTypo(t *testing.T) {
	// arrange
	brokenYamlRelativePath := brokenYamlFolderPath + "broken_step_run_typo.yaml" // job with just a "title" but no "steps" list
	brokenYamlPath_noJobs, _ := filepath.Abs(brokenYamlRelativePath)

	viperParser := prepareViper(brokenYamlPath_noJobs)
	yaml := parseYaml(viperParser, brokenYamlPath_noJobs)

	// act
	err := validateYamlSchema(yaml)

	// assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), `missing "run" key for step on (job name: "job 2", job number: 2, step name: "task 21", step number: 1)`)
}
