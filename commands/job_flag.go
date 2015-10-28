package commands

import (
	"strings"

	"github.com/concourse/go-concourse/concourse"
)

type JobFlag struct {
	PipelineName string
	JobName      string
}

func (job *JobFlag) UnmarshalFlag(value string) error {
	vs := strings.SplitN(value, "/", 2)
	if vs[0] == "" {
		return concourse.NameRequiredError("pipeline")
	}
	if vs[1] == "" {
		return concourse.NameRequiredError("job")
	}

	job.PipelineName = vs[0]
	job.JobName = vs[1]

	return nil
}
