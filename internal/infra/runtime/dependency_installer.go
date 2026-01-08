package runtime

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"kodkafa/internal/domain/entities"
	"kodkafa/internal/domain/ports"
)

type FSInstaller struct {
	baseDir     string
	configStore ports.ConfigStore
}

func NewFSInstaller(baseDir string, configStore ports.ConfigStore) *FSInstaller {
	return &FSInstaller{
		baseDir:     baseDir,
		configStore: configStore,
	}
}

func (i *FSInstaller) Install(plugin *entities.Plugin) error {
	switch plugin.Interpreter {
	case "python":
		return i.installPython(plugin)
	case "node":
		return i.installNode(plugin)
	default:
		return fmt.Errorf("unsupported interpreter for dependency installation: %s", plugin.Interpreter)
	}
}

func (i *FSInstaller) Uninstall(plugin *entities.Plugin) error {
	switch plugin.Interpreter {
	case "python":
		return i.uninstallPython(plugin)
	case "node":
		return i.uninstallNode(plugin)
	default:
		return nil
	}
}

func (i *FSInstaller) installPython(plugin *entities.Plugin) error {
	reqPath := filepath.Join(plugin.Source, "requirements.txt")
	if _, err := os.Stat(reqPath); os.IsNotExist(err) {
		return nil
	}

	// Internal Core Strategy: install into central venv at core/python/venv
	venvPath := filepath.Join(i.baseDir, "core", "python", "venv")
	pipPath := filepath.Join(venvPath, "bin", "pip")

	cmd := exec.Command(pipPath, "install", "-r", "requirements.txt")
	cmd.Dir = plugin.Source
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("core pip install failed: %w (output: %s)", err, string(out))
	}
	return nil
}

func (i *FSInstaller) uninstallPython(plugin *entities.Plugin) error {
	// Pip uninstall is tricky without a specific target, but we'll leave it for now
	// as shared environments typically grow. Unique plugin cleanup is a future task.
	return nil
}

func (i *FSInstaller) installNode(plugin *entities.Plugin) error {
	pkgPath := filepath.Join(plugin.Source, "package.json")
	if _, err := os.Stat(pkgPath); os.IsNotExist(err) {
		return nil
	}

	// 1. Read plugin's package.json to get dependencies
	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return fmt.Errorf("failed to read package.json: %w", err)
	}

	// Minimal struct to capture dependencies
	type PkgJSON struct {
		Dependencies map[string]string `json:"dependencies"`
	}
	var pkg PkgJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return fmt.Errorf("failed to parse package.json: %w", err)
	}

	if len(pkg.Dependencies) == 0 {
		return nil
	}

	// 2. Install dependencies into core/node
	nodeCoreDir := filepath.Join(i.baseDir, "core", "node")
	if err := os.MkdirAll(nodeCoreDir, 0755); err != nil {
		return fmt.Errorf("failed to create core node dir: %w", err)
	}

	// Construct list of packages to install (e.g. "axios@^1.0.0")
	var installArgs []string
	installArgs = append(installArgs, "install")
	for name, version := range pkg.Dependencies {
		installArgs = append(installArgs, fmt.Sprintf("%s@%s", name, version))
	}

	cmd := exec.Command("npm", installArgs...)
	cmd.Dir = nodeCoreDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("npm install in core failed: %w (output: %s)", err, string(out))
	}

	// 3. Cleanup: Ensure no node_modules in plugin dir
	pluginNodeModules := filepath.Join(plugin.Source, "node_modules")
	_ = os.RemoveAll(pluginNodeModules)

	return nil
}

func (i *FSInstaller) uninstallNode(plugin *entities.Plugin) error {
	// For now, we don't prune the central node_modules to avoid breaking other plugins.
	// Just remove the symlink in the plugin directory.
	pluginNodeModules := filepath.Join(plugin.Source, "node_modules")
	return os.RemoveAll(pluginNodeModules)
}
