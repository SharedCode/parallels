package jobs

import "time"

type State int
const (
	Initial = iota
	Started
	Completed
	Aborted
	TimedOut
)

// Job structure contains the basic job information. Feel free to inherit from this and add custom fields to your Job structure.
type Job struct{
	ID string
	Name string
	State State
	ActionInputData interface{}
	ActionResultData interface{}
	ActionError error
	Created time.Time
	Updated time.Time
}

type Repository interface{
	Source(jobName string, jobCount int) ([]string,error)
	Sink(jobs []Job) error
	Update(jobs []Job) error
	Get(jobID []string) ([]string,error)
	Remove(jobID []string) error
}
