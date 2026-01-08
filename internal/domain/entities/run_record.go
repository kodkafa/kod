package entities

import "time"

// RunStatus represents the status of a plugin run.
type RunStatus string

const (
	RunStatusRunning RunStatus = "running"
	RunStatusSuccess RunStatus = "success"
	RunStatusFailed  RunStatus = "failed"
	RunStatusAborted RunStatus = "aborted"
)

// RunRecord represents a single execution record for a plugin.
type RunRecord struct {
	Timestamp time.Time
	Args      string
	ExitCode  int
	Duration  time.Duration
	Status    RunStatus
}
