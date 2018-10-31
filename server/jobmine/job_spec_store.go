package jobmine

import (
	"fmt"

	"github.com/pkg/errors"
)

type JobSpecStore struct {
	JobSpecs map[JobType]JobSpec
}

// GetJobSpecForJobtype Finds a job spec given a job type, using a job spec store.
func (s *JobSpecStore) GetJobSpecForJobType(jobType JobType) (*JobSpec, error) {
	val, ok := s.JobSpecs[jobType]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Unable to find JobSpec job type %s", string(jobType)))
	}
	return &val, nil
}

// GetTaskSpecForJobType Finds a task spec given a job type, using a job spec store.
func (s *JobSpecStore) GetTaskSpecForJobType(jobType JobType) (*TaskSpec, error) {
	val, err := s.GetJobSpecForJobType(jobType)
	if err != nil {
		return nil, err
	}

	return &val.TaskSpec, nil
}
