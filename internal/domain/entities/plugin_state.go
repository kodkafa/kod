package entities

import "time"

// PluginState stores per-plugin execution state and history.
type PluginState struct {
	PluginName     string
	AddedAt        time.Time
	LastExecutedAt time.Time
	RunCount       int
	History        []RunRecord
	MaxHistorySize int
}

// NewPluginState creates a new PluginState with default max history size.
func NewPluginState(pluginName string) *PluginState {
	return &PluginState{
		PluginName:     pluginName,
		AddedAt:        time.Now(),
		History:        make([]RunRecord, 0),
		MaxHistorySize: 50, // default, configurable
	}
}

// AddRunRecord appends a run record and maintains bounded history.
func (ps *PluginState) AddRunRecord(record RunRecord) {
	// Deduplicate: Remove existing record with same args
	filtered := make([]RunRecord, 0, len(ps.History))
	for _, r := range ps.History {
		if r.Args != record.Args {
			filtered = append(filtered, r)
		}
	}
	ps.History = filtered

	ps.History = append(ps.History, record)
	ps.LastExecutedAt = record.Timestamp
	ps.RunCount++

	// Maintain bounded history
	if len(ps.History) > ps.MaxHistorySize {
		ps.History = ps.History[len(ps.History)-ps.MaxHistorySize:]
	}
}

// GetMostRecentArgs returns the args from the most recent run, or empty string.
func (ps *PluginState) GetMostRecentArgs() string {
	if len(ps.History) == 0 {
		return ""
	}
	return ps.History[len(ps.History)-1].Args
}
