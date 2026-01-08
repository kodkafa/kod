package ports

import (
	"kodkafa/internal/domain/entities"
)

// RunMode specifies the execution mode for a plugin.
type RunMode string

const (
	RunModeStreaming   RunMode = "streaming"
	RunModeInteractive RunMode = "interactive"
)

// RunResult represents the result of a plugin execution.
type RunResult struct {
	ExitCode int
	Duration int64 // nanoseconds
	Status   string
	Output   string
}

// OutputChunk represents a chunk of output from a running process.
type OutputChunk struct {
	Data   []byte
	IsErr  bool
	Plugin string
}

// Runner defines the interface for executing plugins.
type Runner interface {
	// Run executes a plugin with the given arguments.
	// It streams output via the provided channel and returns the result.
	Run(plugin *entities.Plugin, args string, mode RunMode, outputChan chan<- OutputChunk) (*RunResult, error)
}
