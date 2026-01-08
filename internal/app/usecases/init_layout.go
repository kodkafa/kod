package usecases

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"kodkafa/internal/domain/ports"
	"kodkafa/internal/infra/runtime"
)

// InitLayoutUseCase handles initialization of the KODKAFA directory structure.
type InitLayoutUseCase struct {
	configStore  ports.ConfigStore
	baseDir      string
	templatePath string
}

// NewInitLayoutUseCase creates a new InitLayoutUseCase.
func NewInitLayoutUseCase(baseDir string, templatePath string, configStore ports.ConfigStore) *InitLayoutUseCase {
	return &InitLayoutUseCase{
		configStore:  configStore,
		baseDir:      baseDir,
		templatePath: templatePath,
	}
}

// Execute creates the directory structure and initializes configuration.
func (uc *InitLayoutUseCase) Execute() error {
	// Create base directory
	if err := os.MkdirAll(uc.baseDir, 0755); err != nil {
		return fmt.Errorf("failed to create base directory: %w", err)
	}

	// Create subdirectories
	dirs := []string{
		"plugins",
		"state",
		"core",
		"logs",
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(uc.baseDir, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Initialize Core Runtimes
	if err := uc.initCoreRuntimes(); err != nil {
		fmt.Printf("Warning: failed to initialize some core runtimes: %v\n", err)
	}

	// Runtime detection
	config, err := uc.configStore.Read()
	if err != nil {
		// Load from template
		data, err := os.ReadFile(uc.templatePath)
		if err != nil {
			return fmt.Errorf("failed to read default config template: %w", err)
		}
		config = &ports.Config{}
		if err := json.Unmarshal(data, config); err != nil {
			return fmt.Errorf("failed to parse default config template: %w", err)
		}
	}

	for key, cmd := range config.SupportedRuntimes {
		path, found := runtime.CheckInterpreter(cmd)
		if found {
			config.RuntimePaths[key] = path
		} else {
			config.RuntimePaths[key] = "undefined"
			fmt.Printf("Warning: %s interpreter (%s) not found in PATH\n", key, cmd)
		}
	}

	if err := uc.configStore.Write(config); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func (uc *InitLayoutUseCase) initCoreRuntimes() error {
	// 1. Python venv
	pyCoreDir := filepath.Join(uc.baseDir, "core", "python")
	venvPath := filepath.Join(pyCoreDir, "venv")
	if _, err := os.Stat(venvPath); os.IsNotExist(err) {
		fmt.Println("Initializing core Python environment...")
		_ = os.MkdirAll(pyCoreDir, 0755)
		cmd := exec.Command("python3", "-m", "venv", "venv")
		cmd.Dir = pyCoreDir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to create core python venv: %w", err)
		}
	}

	// 2. Node package.json
	nodeCoreDir := filepath.Join(uc.baseDir, "core", "node")
	pkgJsonPath := filepath.Join(nodeCoreDir, "package.json")
	if _, err := os.Stat(pkgJsonPath); os.IsNotExist(err) {
		fmt.Println("Initializing core Node.js environment...")
		_ = os.MkdirAll(nodeCoreDir, 0755)
		pkgData := []byte(`{
  "name": "kodkafa-core",
  "version": "1.0.0",
  "private": true,
  "dependencies": {}
}`)
		if err := os.WriteFile(pkgJsonPath, pkgData, 0644); err != nil {
			return fmt.Errorf("failed to create core node package.json: %w", err)
		}
	}

	return nil
}
