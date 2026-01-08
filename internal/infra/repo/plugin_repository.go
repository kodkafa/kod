package repo

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"kodkafa/internal/domain/entities"
	"kodkafa/internal/domain/ports"

	"gopkg.in/yaml.v3"
)

// PluginManifest represents the plugin.yml structure.
type PluginManifest struct {
	Name        string `yaml:"name"`
	Interpreter string `yaml:"interpreter"`
	Description string `yaml:"description"`
	Entry       string `yaml:"entry"`
	Usage       string `yaml:"usage"`
}

// PluginRepositoryImpl implements ports.PluginRepository using the filesystem.
type PluginRepositoryImpl struct {
	baseDir    string
	pluginsDir string
}

// NewPluginRepository creates a new PluginRepository implementation.
func NewPluginRepository(baseDir string) ports.PluginRepository {
	return &PluginRepositoryImpl{
		baseDir:    baseDir,
		pluginsDir: filepath.Join(baseDir, "plugins"),
	}
}

// List returns all installed plugins.
func (pr *PluginRepositoryImpl) List() ([]entities.Plugin, error) {
	if _, err := os.Stat(pr.pluginsDir); os.IsNotExist(err) {
		return []entities.Plugin{}, nil
	}

	var plugins []entities.Plugin

	entries, err := os.ReadDir(pr.pluginsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugins directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pluginPath := filepath.Join(pr.pluginsDir, entry.Name())
		plugin, err := pr.readPlugin(pluginPath)
		if err != nil {
			// Skip invalid plugins but continue
			continue
		}
		plugins = append(plugins, *plugin)
	}

	return plugins, nil
}

// Get returns a plugin by name, or an error if not found.
func (pr *PluginRepositoryImpl) Get(name string) (*entities.Plugin, error) {
	pluginPath := filepath.Join(pr.pluginsDir, name)
	return pr.readPlugin(pluginPath)
}

// Add registers a new plugin from a local path or remote URL.
// For now, only local paths are supported.
func (pr *PluginRepositoryImpl) Add(source string) (*entities.Plugin, error) {
	// Check if source is a local path
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return nil, fmt.Errorf("source path does not exist: %w", err)
	}

	if !sourceInfo.IsDir() {
		return nil, fmt.Errorf("source must be a directory")
	}

	// Read plugin manifest
	manifestPath := filepath.Join(source, "plugin.yml")
	manifest, err := pr.readManifest(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read plugin manifest: %w", err)
	}

	// Validate manifest
	if manifest.Name == "" {
		return nil, fmt.Errorf("missing config: name is required")
	}
	if manifest.Interpreter == "" {
		return nil, fmt.Errorf("missing config: interpreter is required")
	}
	if manifest.Entry == "" {
		return nil, fmt.Errorf("missing config: entry is required")
	}

	// Check if plugin already exists in plugins/
	targetDir := filepath.Join(pr.pluginsDir, manifest.Name)
	if _, err := os.Stat(targetDir); err == nil {
		return nil, fmt.Errorf("plugin %s already exists", manifest.Name)
	}

	// Copy plugin directory to plugins/
	if err := pr.copyDirectory(source, targetDir); err != nil {
		return nil, fmt.Errorf("failed to copy plugin: %w", err)
	}

	// Create plugin entity
	plugin := &entities.Plugin{
		Name:        manifest.Name,
		Interpreter: manifest.Interpreter,
		Description: manifest.Description,
		Entry:       manifest.Entry,
		Usage:       manifest.Usage,
		Source:      targetDir,
		AddedAt:     time.Now(),
	}

	return plugin, nil
}

// Remove removes a plugin by name.
func (r *PluginRepositoryImpl) Remove(name string) error {
	path := filepath.Join(r.pluginsDir, name)
	return os.RemoveAll(path)
}

// RemoveDeps removes dependencies for a plugin.
func (r *PluginRepositoryImpl) RemoveDeps(name string) error {
	plugin, err := r.Get(name)
	if err != nil {
		return err
	}

	// Shared Node/Python cleanup is now handled by DependencyInstaller.Uninstall
	// This method remains for local (within-plugin-folder) cleanup if any.

	var depsFolders []string
	switch plugin.Interpreter {
	case "python":
		depsFolders = []string{"venv", ".venv", "__pycache__"}
	case "node", "javascript", "typescript":
		depsFolders = []string{"node_modules"}
	case "r":
		depsFolders = []string{"renv", ".Rproj.user"}
	}

	for _, folder := range depsFolders {
		path := filepath.Join(plugin.Source, folder)
		if _, err := os.Stat(path); err == nil {
			if err := os.RemoveAll(path); err != nil {
				fmt.Printf("Warning: failed to remove dependency folder %s: %v\n", path, err)
			}
		}
	}

	return nil
}

// Exists checks if a plugin with the given name exists.
func (pr *PluginRepositoryImpl) Exists(name string) (bool, error) {
	pluginPath := filepath.Join(pr.pluginsDir, name)
	_, err := os.Stat(pluginPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// readPlugin reads a plugin from the filesystem.
func (pr *PluginRepositoryImpl) readPlugin(path string) (*entities.Plugin, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get plugin info: %w", err)
	}

	manifestPath := filepath.Join(path, "plugin.yml")
	manifest, err := pr.readManifest(manifestPath)
	if err != nil {
		return nil, err
	}

	return &entities.Plugin{
		Name:        manifest.Name,
		Interpreter: manifest.Interpreter,
		Description: manifest.Description,
		Entry:       manifest.Entry,
		Usage:       manifest.Usage,
		Source:      path,
		AddedAt:     info.ModTime(),
	}, nil
}

// readManifest reads and parses a plugin.yml file.
func (pr *PluginRepositoryImpl) readManifest(path string) (*PluginManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}

	var manifest PluginManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	return &manifest, nil
}

// copyDirectory copies a directory recursively.
func (pr *PluginRepositoryImpl) copyDirectory(src, dst string) error {
	// Create destination directory
	if err := os.MkdirAll(dst, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// Read source file
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Write destination file
		return os.WriteFile(dstPath, data, info.Mode())
	})
}
