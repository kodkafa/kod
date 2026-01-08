package ports

// Config represents the global configuration.
type Config struct {
	TrustedDomains     []string          `json:"trusted_domains"`
	RuntimePaths       map[string]string `json:"runtime_paths"`
	SortBy             string            `json:"sort_by"`
	ItemsPerPage       int               `json:"items_per_page"`
	ShowLastRuns       bool              `json:"show_last_runs"`
	FavLimit           int               `json:"fav_limit"`
	LastRunOrder       string            `json:"last_run_order"`
	LastRunLimit       int               `json:"last_run_limit"`
	HistorySize        int               `json:"history_size"`
	DependencySettings map[string]any    `json:"dependency_settings"`
	SupportedRuntimes  map[string]string `json:"supported_runtimes"`
	Splash             bool              `json:"splash"`
}

// ConfigStore defines the interface for configuration persistence.
type ConfigStore interface {
	// Read reads the global configuration.
	Read() (*Config, error)
	// Write persists the global configuration.
	Write(config *Config) error
}
