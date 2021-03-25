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
	brokenYamlRelativePath := brokenYamlFolderPath + "broken_no_step_name.yaml"
	brokenYamlPath_noJobs, _ := filepath.Abs(brokenYamlRelativePath)

	viperParser := prepareViper(brokenYamlPath_noJobs)
	yaml := parseYaml(viperParser, brokenYamlPath_noJobs)

	// act
	err := validateYamlSchema(yaml)

	// assert
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), `missing "name" key for step on (job name: "job 1", job number: 1, step number: 3)`)
}
