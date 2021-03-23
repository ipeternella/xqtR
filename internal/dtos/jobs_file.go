// Package dtos contains the structs that are merely "data containers, i.e. data structures that
// contain no logic, just bare data.
package dtos

type JobsFile struct {
	Jobs []Job `mapstructure:"jobs"`
}

type Job struct {
	Title      string    `mapstructure:"title"` // job name
	Steps      []JobStep `mapstructure:"steps"`
	NumWorkers int       `mapstructure:"num_workers"`
}

type JobStep struct {
	Name string `mapstructure:"name"`
	Run  string `mapstructure:"run"`
}

func NewJobsFile() *JobsFile {
	return &JobsFile{}
}
