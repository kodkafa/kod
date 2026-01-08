package ports

import "kodkafa/internal/domain/entities"

// DependencyInstaller defines the interface for installing plugin dependencies.
type DependencyInstaller interface {
	Install(plugin *entities.Plugin) error
	Uninstall(plugin *entities.Plugin) error
}
