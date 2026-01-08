package entities

import "time"

// UsageEntry represents a single usage entry.
type UsageEntry struct {
	PluginName string
	Timestamp  time.Time
	RunCount   int
}

// UsageStats stores global usage statistics.
type UsageStats struct {
	RecentlyUsed []UsageEntry
	MostUsed     []UsageEntry
	MaxRecent    int
	MaxMostUsed  int
}

// NewUsageStats creates a new UsageStats with default limits.
func NewUsageStats() *UsageStats {
	return &UsageStats{
		RecentlyUsed: make([]UsageEntry, 0),
		MostUsed:     make([]UsageEntry, 0),
		MaxRecent:    20,
		MaxMostUsed:  10,
	}
}

// RecordRun records a plugin run and updates both recent and most-used lists.
func (us *UsageStats) RecordRun(pluginName string, limit int) {
	now := time.Now()

	// Update or add to recently used
	found := false
	for i := range us.RecentlyUsed {
		if us.RecentlyUsed[i].PluginName == pluginName {
			us.RecentlyUsed[i].Timestamp = now
			us.RecentlyUsed[i].RunCount++
			found = true
			break
		}
	}
	if !found {
		us.RecentlyUsed = append(us.RecentlyUsed, UsageEntry{
			PluginName: pluginName,
			Timestamp:  now,
			RunCount:   1,
		})
	}

	// Sort recently used by timestamp (newest first)
	us.sortRecentlyUsed()

	// Update most used
	us.updateMostUsed(pluginName, limit)

	// Maintain bounds
	if len(us.RecentlyUsed) > limit {
		us.RecentlyUsed = us.RecentlyUsed[:limit]
	}
}

func (us *UsageStats) sortRecentlyUsed() {
	// Simple bubble sort by timestamp (newest first)
	for i := 0; i < len(us.RecentlyUsed)-1; i++ {
		for j := 0; j < len(us.RecentlyUsed)-1-i; j++ {
			if us.RecentlyUsed[j].Timestamp.Before(us.RecentlyUsed[j+1].Timestamp) {
				us.RecentlyUsed[j], us.RecentlyUsed[j+1] = us.RecentlyUsed[j+1], us.RecentlyUsed[j]
			}
		}
	}
}

func (us *UsageStats) updateMostUsed(pluginName string, limit int) {
	found := false
	for i := range us.MostUsed {
		if us.MostUsed[i].PluginName == pluginName {
			us.MostUsed[i].RunCount++
			found = true
			break
		}
	}
	if !found {
		us.MostUsed = append(us.MostUsed, UsageEntry{
			PluginName: pluginName,
			Timestamp:  time.Now(),
			RunCount:   1,
		})
	}

	// Sort by run count (highest first)
	for i := 0; i < len(us.MostUsed)-1; i++ {
		for j := 0; j < len(us.MostUsed)-1-i; j++ {
			if us.MostUsed[j].RunCount < us.MostUsed[j+1].RunCount {
				us.MostUsed[j], us.MostUsed[j+1] = us.MostUsed[j+1], us.MostUsed[j]
			}
		}
	}

	// Maintain bounds
	if len(us.MostUsed) > limit {
		us.MostUsed = us.MostUsed[:limit]
	}
}
