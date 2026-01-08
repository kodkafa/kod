package entities

import "time"

// Plugin represents a registered plugin with its metadata.
type Plugin struct {
	Name        string
	Interpreter string
	Description string
	Entry       string
	Usage       string
	Source      string
	AddedAt     time.Time
}
